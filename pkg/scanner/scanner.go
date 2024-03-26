package scanner

import (
	"context"
	"sync"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/tracker"
)

// Scanner struct uses private fields to store scanner state and provides methods
// for running a scan of a git repository for PHI/PII data, recursively scanning
// the contents of files committed to each repository.
type Scanner struct {
	ID string `json:"id"`
	// URL associated with the object, such as the repository URL
	URL string `json:"url"`

	TrackerCommits  *tracker.KeyTracker
	TrackerFiles    *tracker.KeyTracker
	TrackerRequests *tracker.KeyTracker

	chan_commits     chan *object.Commit
	chan_requests    chan rrr.Request
	chan_errors      chan error
	ctx              context.Context
	git_config       *cfg.GitConfig
	is_scan_complete bool
	logger           *zerolog.Logger
	repository       *git.Repository
	result_io        rrr.ResultRecordIO
	scan_mutex       *sync.RWMutex
}

// NewScanner() function initializes a new Scanner object.
func NewScanner(
	ctx context.Context,
	git_config *cfg.GitConfig,
	result_io rrr.ResultRecordIO,
) (*Scanner, error) {
	// ensure the context is not nil
	if ctx == nil {
		ctx = context.Background()
	}

	// create a logger from the context
	logger := zerolog.Ctx(ctx)

	tracker_commits, tc_err := tracker.NewKeyTracker(tracker.ScanObjectTypeCommit, logger)
	if tc_err != nil {
		return nil, errors.Wrap(tc_err, ErrMsgScannerCreate)
	}
	tracker_files, tf_err := tracker.NewKeyTracker(tracker.ScanObjectTypeFile, logger)
	if tf_err != nil {
		return nil, errors.Wrap(tf_err, ErrMsgScannerCreate)
	}
	// create a tracker.KeyTracker for tracking (responses to) requests
	tracker_requests, tr_err := tracker.NewKeyTracker(tracker.ScanObjectTypeRequestResponse, logger)
	if tr_err != nil {
		return nil, errors.Wrap(tr_err, ErrMsgScannerCreate)
	}

	return &Scanner{
		ID:              uuid.NewString(),
		TrackerCommits:  tracker_commits,
		TrackerFiles:    tracker_files,
		TrackerRequests: tracker_requests,
		chan_commits:    make(chan *object.Commit),
		chan_errors:     make(chan error),
		chan_requests:   make(chan rrr.Request),
		ctx:             ctx,
		git_config:      git_config,
		logger:          logger,
		result_io:       result_io,
		scan_mutex:      &sync.RWMutex{},
	}, nil
}

// ScanInput struct defines the required input parameters for the Scanner.Scan()
// method.
type ScanInput struct {
	ChanErrorsSend      chan<- error
	ChanRequestSend     chan<- rrr.Request
	ChanResponseReceive <-chan rrr.Response
	RepoID              string
	Repository          *git.Repository
}

// Scan() method uses channels and goroutines to coordinate the scanning of
// a git repository for PHI/PII data.
func (s *Scanner) Scan(in ScanInput) {
	s.logger.Debug().Msg("started Scanner run")
	defer s.logger.Debug().Msg("finished Scanner run")

	// check if a previous scan created a Checkpoint file from which to resume
	cpoint, cpoint_err := CheckpointGet(s.ctx, s.git_config.WorkDir, in.RepoID, "")
	if cpoint_err != nil {
		s.logger.Error().Err(cpoint_err).Msg("failed to initialize scan tracker with checkpoint data")
	}
	if cpoint != nil {
		// use the Checkpoint data to restore state from a previous scan
		s.TrackerCommits.Restore(cpoint.TrackerCommitsData)
		s.TrackerFiles.Restore(cpoint.TrackerFilesData)
		s.TrackerRequests.Restore(cpoint.TrackerRequestsData)
	}

	// create channels for coordinating between goroutines
	chan_scan_done := make(chan struct{})
	chan_quit := make(chan struct{})

	// track the progress of the scan
	go s.trackScanProgress(chan_scan_done, chan_quit)
	// listen for errors generated by the scan
	go s.processErrors(chan_quit, s.chan_errors, in.ChanErrorsSend)
	// process requests generated by the scan
	go s.processRequests(
		chan_quit,
		s.chan_requests,
		in.ChanRequestSend,
		s.chan_errors,
	)
	// process responses for requests
	go s.processResponses(
		chan_quit,
		in.ChanResponseReceive,
		s.chan_errors,
	)
	// scan the repository
	go s.scanRepository(
		in.RepoID,
		in.Repository,
		s.chan_errors,
		chan_scan_done,
	)

	// listen for quit signal
	// TODO : replace with `go s.processResults()`
	<-chan_quit
}

// checkpointScan() method is intended to be run as a separate goroutine
// to periodically checkpoint the progress of the scan.
func (s *Scanner) checkpointScan(
	repo_id string,
	commit_id string,
	chan_quit_in <-chan struct{},
	chan_errors_out chan<- error,
) {
	s.logger.Debug().Msg("started scan progress checkpoint processor")
	defer s.logger.Debug().Msg("finished scan progress checkpoint processor")

	setNewCheckpoint := func() (e error) {
		// store the scan progress in a Checkpoint file
		e = CheckpointSet(
			s.ctx,
			s.git_config.WorkDir,
			repo_id,
			commit_id,
			NewCheckpoint(
				s.TrackerCommits.GetKeysData(),
				s.TrackerFiles.GetKeysData(),
				s.TrackerRequests.GetKeysData(),
			),
		)
		if e != nil {
			return
		}
		return
	}

	// set the first Checkpoint before starting (and waiting for) the ticker
	if err := setNewCheckpoint(); err != nil {
		chan_errors_out <- errors.Wrap(err, ErrMsgCheckpointScanProgress)
	}

	// create a ticker to periodically trigger a refresh of the scan checkpoint
	timer := time.NewTicker(CheckpointRefreshInterval)

	for {
		select {
		case <-timer.C:
			// store the scan progress in a Checkpoint file
			if err := setNewCheckpoint(); err != nil {
				chan_errors_out <- errors.Wrap(err, ErrMsgCheckpointScanProgress)
			}
		case <-chan_quit_in:
			return
		}
	}
}

// processCommits() method is intended to be run as a goroutine to process
// commits from the channel of commits generated by the commit iterator.
func (s *Scanner) processCommits(wg_main *sync.WaitGroup) {
	defer wg_main.Done()

	wg_loop := &sync.WaitGroup{}
	processCommit := func(commit *object.Commit) {
		defer wg_loop.Done()

		s.logger.Debug().Msgf(
			"repository %s : scanning commit %s",
			s.URL,
			commit.Hash.String(),
		)

		// get the tree of objects associated with the commit
		tree, err := commit.Tree()
		if err != nil {
			_, err = s.TrackerCommits.Update(
				commit.Hash.String(),
				tracker.KeyCodeError,
				err.Error(),
				[]string{},
			)
			if err != nil {
				err = errors.Wrapf(err, ErrMsgTrackerUpdateCommit, commit.Hash.String())
				s.chan_errors <- err
				return
			}
		}

		// iterate through the files in the commit tree
		err = tree.Files().ForEach(s.scanFile(commit))
		if err != nil {
			err = errors.Wrapf(err, ErrMsgTrackerUpdateCommit, commit.Hash.String())
			s.TrackerCommits.Update(
				commit.Hash.String(),
				tracker.KeyCodeError,
				err.Error(),
				[]string{},
			)
			s.chan_errors <- err
			return
		}

		// attempt to update the commit code to "complete" status, but ignore any error
		// and accept that the commit may be left in "pending" status if the key has
		// children that are still in an incomplete (bool=false) state
		s.TrackerCommits.Update(
			commit.Hash.String(),
			tracker.KeyCodeComplete,
			"",
			[]string{},
		)
	}
	for commit := range s.chan_commits {
		wg_loop.Add(1)
		go processCommit(commit)
	}
	wg_loop.Wait()
}

// processErrors() method processes errors generated by the scan.
func (s *Scanner) processErrors(
	chan_quit_in <-chan struct{},
	chan_errors_in <-chan error,
	chan_errors_out chan<- error,
) {
	s.logger.Debug().Msg("started error processor")
	defer s.logger.Debug().Msg("finished error processor")

	for {
		select {
		case <-chan_quit_in:
			s.logger.Warn().Msg("scanner error processor received quit signal")
			close(chan_errors_out)
			return
		case e := <-chan_errors_in:
			if e != nil {
				err_wrap_msg := "error running scanner"
				s.logger.Error().Err(e).Msg(err_wrap_msg)
				// handle error to determine if the scanner should continue
				// TODO
				chan_errors_out <- e
				close(chan_errors_out)
				return
			}
			return
		}
	}
}

// processRequest() method processes a single request for internal tracking
// purposes before sending the request for external processing.
func (s *Scanner) processRequest(
	r rrr.Request,
	chan_requests_out chan<- rrr.Request,
	chan_errors_out chan<- error,
) {
	// validate the request
	if r.ID == "" {
		chan_errors_out <- ErrProcessRequestNoID
		return
	}
	// check if the request is already being tracked
	if _, exists := s.TrackerRequests.Get(r.ID); exists {
		s.logger.Debug().Msgf("skipping processing for existing request ID=%s", r.ID)
		return
	}
	// update TrackerRequests to track the ID of the pending request
	_, err := s.TrackerRequests.Update(r.ID, tracker.KeyCodePending, "", []string{})
	if err != nil {
		chan_errors_out <- err
		return
	}
	// send the request for external processing
	chan_requests_out <- r
}

// processRequests() method processes requests for documents generated by
// the scan.
func (s *Scanner) processRequests(
	chan_quit_in <-chan struct{},
	chan_requests_in <-chan rrr.Request,
	chan_requests_out chan<- rrr.Request,
	chan_errors_out chan<- error,
) {
	s.logger.Debug().Msg("started requests processor")
	defer s.logger.Debug().Msg("finished requests processor")

	// listen for requests to process
	for {
		select {
		case <-chan_quit_in:
			return
		case r := <-chan_requests_in:
			// keep the input channel clear by processing the request in the
			// background via a separate goroutine, which sends any errors to
			// chan_errors_out
			s.processRequest(r, chan_requests_out, chan_errors_out)
		}
	}
}

// processResponse() method processes a single response to some request.
func (s *Scanner) processResponse(
	r rrr.Response,
	chan_errors_out chan<- error,
) {
	// validate the response
	if r.ID == "" {
		chan_errors_out <- ErrProcessResponseNoID
		return
	}
	// log the response
	s.logger.Trace().Msgf(
		"processing %d results for request/response ID = %s : Repository.ID : %s : Commit.ID = %s : Object.ID = %s",
		len(r.Results),
		r.ID,
		r.Repository.ID,
		r.Commit.ID,
		r.Object.ID,
	)
	// write the result(s) to the result_io store
	if len(r.Results) > 0 {
		// convert the response to a slice of rrr.ResultRecords, where
		// each rrr.ResultRecord is uniquely identified by its SHA1 hash
		result_records := rrr.ResultRecordsFromResponse(&r)
		if err := s.result_io.Write(result_records); err != nil {
			chan_errors_out <- errors.Wrap(err, ErrMsgResultWriteFailed)
		}
	}
	// update TrackerRequests to mark the associated request (ID) as complete
	s.TrackerRequests.Update(r.ID, tracker.KeyCodeComplete, "", []string{})

	// update the tracker for the associated File object to mark this
	// request/response as complete. if all requests/responses for a
	// File object are complete, the File object should be marked as
	// tracker.KeyCodeComplete.
	var file_update_code int
	var update_err error
	file_update_code, update_err = s.TrackerFiles.Update(
		r.Object.ID,
		tracker.KeyCodeComplete,
		"",
		[]string{r.ID},
	)
	if update_err != nil {
		chan_errors_out <- update_err
	}
	// only update the associated commit if file_update_code is tracker.KeyCodeComplete.
	if file_update_code == tracker.KeyCodeComplete {
		_, update_err = s.TrackerCommits.Update(
			r.Commit.ID,
			tracker.KeyCodeComplete,
			"",
			[]string{r.Object.ID},
		)
		if update_err != nil {
			chan_errors_out <- update_err
		}
	}
}

// processResponses() method processes all responses for requests generated by
// the scan.
func (s *Scanner) processResponses(
	chan_quit_in <-chan struct{},
	chan_responses_in <-chan rrr.Response,
	chan_errors_out chan<- error,
) {
	s.logger.Debug().Msg("started response processor")
	defer s.logger.Debug().Msg("finished response processor")

	// listen for responses to process
	for {
		select {
		case <-chan_quit_in:
			return
		case r := <-chan_responses_in:
			// keep the input channel clear by processing the response in the
			// background via a separate goroutine, which sends any errors to
			// chan_errors_out
			s.processResponse(r, chan_errors_out)
		}
	}
}

// reconcilePending() method reconciles pending commits and files by checking
// if all requests for a file are actually complete and if all files for a commit
// are actually complete, then updates the trackers for files and commits to
// reflect the actual state of the scan.
func (s *Scanner) reconcilePending() {
	s.logger.Debug().Msg("started scanner reconciler")
	defer s.logger.Debug().Msg("finished scanner reconciler")

	// get the pending files for the repository
	pending_files, files_err := s.TrackerFiles.GetKeysDataForCode(tracker.KeyCodePending)
	if files_err != nil {
		s.logger.Error().Err(files_err).Msg("error getting pending files")
		return
	}
	for file_key, file_key_data := range pending_files {
		// keep track of requests that should be marked as complete in the
		// next update to the tracker for this file
		var requests_complete []string
		// iterate over the request children of the file
		for request_id, is_complete := range file_key_data.Children {
			if is_complete {
				continue
			}
			// get the tracker.KeyData for the pending request ID
			request_data, request_exists := s.TrackerRequests.Get(request_id)
			if !request_exists {
				s.logger.Error().Msgf("error getting request ID=%s", request_id)
				continue
			}
			if request_data.Code == tracker.KeyCodeComplete {
				requests_complete = append(requests_complete, request_id)
			}
		}
		if len(requests_complete) > 0 {
			// update the tracker for the file to mark the requests as complete
			_, update_err := s.TrackerFiles.Update(
				file_key,
				tracker.KeyCodeComplete,
				"",
				requests_complete,
			)
			if update_err != nil {
				s.logger.Error().Err(update_err).Msg("error updating file tracker")
			}
		}
	}

	// get the pending commits for the repository
	pending_commits, commits_err := s.TrackerCommits.GetKeysDataForCode(tracker.KeyCodePending)
	if commits_err != nil {
		s.logger.Error().Err(commits_err).Msg("error getting pending commits")
		return
	}
	for commit_key, commit_key_data := range pending_commits {
		// keep track of files that should be marked as complete in the
		// next update to the tracker for this commit
		var files_complete []string
		// iterate over the file children of the commit
		for file_key, is_complete := range commit_key_data.Children {
			if is_complete {
				continue
			}
			// get the tracker.KeyData for the pending file ID/key
			file_data, file_exists := s.TrackerFiles.Get(file_key)
			if !file_exists {
				s.logger.Error().Msgf("error getting file ID=%s", file_key)
				continue
			}
			if file_data.Code == tracker.KeyCodeComplete {
				files_complete = append(files_complete, file_key)
			}
		}
		if len(files_complete) > 0 {
			// update the tracker for the commit to mark the files as complete
			_, update_err := s.TrackerCommits.Update(
				commit_key,
				tracker.KeyCodeComplete,
				"",
				files_complete,
			)
			if update_err != nil {
				s.logger.Error().Err(update_err).Msg("error updating commit tracker")
			}
		}
	}
}

// scanCommit() method scans the tree of the object.Commit for files
// containing any PHI/PII entities.
func (s *Scanner) scanCommit(commit *object.Commit) error {
	update_code, init_err := s.TrackerCommits.Update(
		commit.Hash.String(),
		tracker.KeyCodeInit,
		"",
		[]string{},
	)
	if init_err != nil {
		return errors.Wrapf(init_err, ErrMsgTrackerUpdateCommit, commit.Hash.String())
	}

	// skip commits that have already been scanned
	if update_code > tracker.KeyCodeInit {
		s.logger.Trace().Msgf(
			"repository %s : skipping previously scanned commit %s",
			s.URL,
			commit.Hash.String(),
		)
		return nil
	}

	// send the commit to the channel for processing
	s.chan_commits <- commit

	return nil
}

// scanFile() method returns an anonymous function that can be used to iterate through
// the files in the associated commit tree and scan each file for PHI/PII entities.
func (s *Scanner) scanFile(commit *object.Commit) func(*object.File) error {
	return func(file *object.File) error {
		code, err := s.TrackerFiles.Update(
			file.Hash.String(),
			tracker.KeyCodeInit,
			"",
			[]string{},
		)
		if err != nil {
			return errors.Wrapf(err, ErrMsgScanTrackerUpdateFile, file.Hash.String())
		}
		// skip files that have already been scanned
		if code > tracker.KeyCodeInit {
			s.logger.Trace().Msgf(
				"commit %s : skipping previously scanned file %s : code=%d",
				commit.Hash.String(),
				file.Hash.String(),
				code,
			)
			return nil
		}

		// check if the file should be ignored instead of scanned
		should_ignore, ignore_reason := IgnoreFileObject(
			file,
			s.git_config.Scan.Extensions,
			s.git_config.Scan.IgnoreExtensions,
		)
		if should_ignore {
			s.logger.Trace().Msgf(
				"commit %s : skipping scan of file %s : %s",
				commit.Hash.String(),
				file.Hash.String(),
				ignore_reason,
			)
			_, err = s.TrackerFiles.Update(
				file.Hash.String(),
				tracker.KeyCodeIgnore,
				ignore_reason,
				[]string{},
			)
			return err
		}
		if ignore_reason != "" {
			s.logger.Warn().Msgf(
				"commit %s : file %s : ignore reason => %s",
				commit.Hash.String(),
				file.Hash.String(),
				ignore_reason,
			)
		}
		s.logger.Debug().Msgf(
			"commit %s : scanning file %s : %s",
			commit.Hash.String(),
			file.Hash.String(),
			file.Name,
		)
		// generate and send requests for the contents of the file
		requests, r_err := rrr.ChunkFileToRequests(rrr.ChunkFileInput{
			CommitID:     commit.Hash.String(),
			File:         file,
			MaxChunkSize: s.git_config.Scan.Limits.MaxRequestChunkSize,
			RepoID:       s.ID,
		})
		if r_err != nil {
			s.logger.Error().Err(r_err).Msgf("commit %s : failed to generate requests for file %s", commit.Hash.String(), file.Hash.String())
			s.TrackerFiles.Update(
				file.Hash.String(),
				tracker.KeyCodeError,
				r_err.Error(),
				[]string{},
			)
			return r_err
		}
		// any zero-size file should have been ignored by the IgnoreFileObject() function
		if len(requests) == 0 && file.Size > 0 {
			err = errors.New("no requests generated for file ID=" + file.Hash.String())
			s.TrackerFiles.Update(
				file.Hash.String(),
				tracker.KeyCodeError,
				err.Error(),
				[]string{},
			)
			s.logger.Warn().Msgf(
				"commit %s : %s : Name=%s : size=%d",
				commit.Hash.String(),
				err.Error(),
				file.Name,
				file.Size,
			)
			return err
		}
		var child_keys []string
		// send each request to the channel for processing
		for _, req := range requests {
			child_keys = append(child_keys, req.ID)
			s.chan_requests <- req
		}
		// update tracker to mark the scan of this file as "pending"
		_, err = s.TrackerFiles.Update(
			file.Hash.String(),
			tracker.KeyCodePending,
			"",
			child_keys,
		)
		if err != nil {
			return errors.Wrapf(err, ErrMsgScanTrackerUpdateFile, file.Hash.String())
		}
		// update tracker for the associated commit to indicate "pending" status
		_, err = s.TrackerCommits.Update(
			commit.Hash.String(),
			tracker.KeyCodePending,
			"",
			[]string{file.Hash.String()},
		)
		if err != nil {
			return errors.Wrapf(err, ErrMsgTrackerUpdateCommit, commit.Hash.String())
		}

		return nil
	}
}

// scanRepository() method scans the repositories defined in the git config
// and sends the results to the requests channel. If an error occurs during
// the scan, the error is sent to the error channel.
func (s *Scanner) scanRepository(
	repo_url string,
	repository *git.Repository,
	errors_out chan<- error,
	done chan struct{},
) {
	s.logger.Debug().Msgf("started scan of repository %s", s.URL)
	defer s.logger.Debug().Msgf("finished scan of repository %s", s.URL)
	defer close(done)

	if errors_out == nil {
		s.logger.Panic().Msg(ErrMsgErrorChannelNil)
	}
	if repository == nil {
		errors_out <- ErrScannerRepositoryNil
	}

	s.scan_mutex.Lock()
	s.repository = repository
	s.URL = repo_url
	s.scan_mutex.Unlock()

	// run a goroutine that periodically checkpoints of the state of the scan
	go s.checkpointScan(repo_url, "", done, s.chan_errors)

	var e error
	// get an iterator for the commits in the repository
	var commit_iterator object.CommitIter
	commit_iterator, e = s.repository.CommitObjects()
	if e != nil {
		if commit_iterator != nil {
			commit_iterator.Close()
		}
		return
	}
	defer commit_iterator.Close()

	wg := &sync.WaitGroup{}
	// start a goroutine to process commits generated by the iterator
	wg.Add(1)
	go s.processCommits(wg)

	// iterate through the commits in the repository history
	e = commit_iterator.ForEach(s.scanCommit)
	if e != nil {
		// wrap the error and send it to the errors channel
		errors_out <- errors.Wrapf(e, "failed to iterate through commits in repository %s", s.URL)
		//return // TODO: should we return here?
	}
	// close the channel for commits to signal the processCommits goroutine
	// to finish
	close(s.chan_commits)
	// wait for the processCommits goroutine to finish
	wg.Wait()

	// set the scan complete flag to true
	s.is_scan_complete = true
}

// trackScanProgress() method tracks the progress of the scan by periodically
// checking if all requests have been completed. If the scan is complete, the
// method returns. If the scan is not complete, the method continues to track
// the progress of the scan by printing the status counts for the requests.
func (s *Scanner) trackScanProgress(
	scan_done_in <-chan struct{},
	quit_out chan<- struct{},
) {
	s.logger.Debug().Msg("started scan progress tracker")
	defer s.logger.Debug().Msg("finished scan progress tracker")

	printScanCounts := func() {
		// print the counts from s.TrackerCommits
		s.TrackerCommits.PrintCounts()
		// print the counts from s.TrackerFiles
		s.TrackerFiles.PrintCounts()
		// print the counts from TrackerRequests
		s.TrackerRequests.PrintCounts()
	}

	trackScanCounts := func() (done bool) {
		done = false
		// print the scan counts regardless of whether the scan is complete
		printScanCounts()

		// check if the scan of the repository is complete, including all
		// requests/responses created from the scan of the repository
		if !s.is_scan_complete {
			s.logger.Debug().Msgf("tracking scan : scan in-progress for repository %s", s.URL)
			return
		}
		// reconcile tracker states across requests, files, and commits
		s.reconcilePending()
		// check if any files are still pending
		if !s.TrackerFiles.CheckAllComplete() {
			s.logger.Debug().Msgf("tracking scan : not all FILES complete for repository %s", s.URL)
			return
		}
		// check if any commits are still pending
		if !s.TrackerCommits.CheckAllComplete() {
			s.logger.Debug().Msgf("tracking scan : not all COMMITS complete for repository %s", s.URL)
			return
		}
		// check if any requests are still pending
		if !s.TrackerRequests.CheckAllComplete() {
			s.logger.Debug().Msgf("tracking scan : not all REQUESTS complete for repository %s", s.URL)
			return
		}
		s.logger.Debug().Msgf("tracking scan : cleaning up scan for repository %s", s.URL)

		// remove the checkpoint file when tracking indicates the scan is complete
		if err := CheckpointDelete(s.ctx, s.git_config.WorkDir, s.URL, ""); err != nil {
			s.logger.Error().Err(err).Msg("Scanner failed to delete Checkpoint file")
		}

		// print the scan counts again before actually cleaning up
		printScanCounts()

		close(quit_out)
		done = true
		return
	}

	// create a ticker to periodically trigger a refresh of progress tracking
	timer := time.NewTicker(ScanRefreshInterval)

	// use scan_done var to avoid repeated processing of scan_done_in
	var scan_done bool = false

	for {
		select {
		case <-timer.C:
			// print tracker counts, then wait for the next tick
			if trackScanCounts() {
				return
			}
		case <-scan_done_in:
			if !scan_done {
				s.logger.Debug().Msg("received scan done signal")
				scan_done = true
				if trackScanCounts() {
					return
				}
			}
			continue
		}
	}
}
