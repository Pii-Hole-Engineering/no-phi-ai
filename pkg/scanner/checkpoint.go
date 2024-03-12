package scanner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	nogit "github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/no-git"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/tracker"
)

// Checkpoints struct defines the structure of the data used to save and restore
// the state of the scanner from a checkpoint in time.
type Checkpoint struct {
	CreatedAt           int64              `json:"created_at"`
	TrackerCommitsData  tracker.KeyDataMap `json:"commits"`
	TrackerFilesData    tracker.KeyDataMap `json:"files"`
	TrackerRequestsData tracker.KeyDataMap `json:"requests"`
}

// NewCheckpoint() function creates a new Checkpoint struct with the given data
// and sets the CreatedAt field to the current time.
func NewCheckpoint(
	data_commits tracker.KeyDataMap,
	data_files tracker.KeyDataMap,
	data_requests tracker.KeyDataMap,
) *Checkpoint {
	return &Checkpoint{
		CreatedAt:           rrr.TimestampNow(),
		TrackerCommitsData:  data_commits,
		TrackerFilesData:    data_files,
		TrackerRequestsData: data_requests,
	}
}

// CheckpointDelete() function deletes the Checkpoint file from the expected file
// path, based on the given repository and (optional) commit ID. Returns a non-nil
// error if unable to locate and delete the expected file path.
func CheckpointDelete(ctx context.Context, work_dir, repo_url, commit_id string) error {
	logger := zerolog.Ctx(ctx)
	file_path, err := getCheckpointPath(work_dir, repo_url, commit_id)
	if err != nil {
		return errors.Wrap(ErrCheckpointDeleteFailed, err.Error())
	}
	if file_path == "" {
		return ErrCheckpointDeleteFailed
	}
	logger.Debug().Msgf("deleting scan checkpoint file: %s", file_path)
	err = os.Remove(file_path)
	if err != nil {
		return errors.Wrap(ErrCheckpointDeleteFailed, err.Error())
	}
	logger.Info().Msgf("deleted scan checkpoint file: %s", file_path)
	return nil
}

// CheckpointGet() function retrieves the Checkpoint data from the checkpoint file
// for the given repository and commit ID. Returns a non-nil error if unable to read
// valid Checkpoint data from the expected file path.
func CheckpointGet(ctx context.Context, work_dir, repo_url, commit_id string) (cpoint *Checkpoint, e error) {
	logger := zerolog.Ctx(ctx)
	var file_path string
	file_path, e = getCheckpointPath(work_dir, repo_url, commit_id)
	if e != nil {
		return
	}

	var file *os.File
	file, e = openCheckpointFile(work_dir, repo_url, commit_id)
	if e != nil {
		return
	}
	file_info, err := file.Stat()
	if err != nil {
		e = err
		return
	}

	if file_info.Size() == 0 {
		e = errors.Wrap(ErrCheckpointFileReadFailed, "file size is 0")
		return
	}

	data_encoded := make([]byte, file_info.Size())
	_, e = file.Read(data_encoded)
	if e != nil {
		e = errors.Wrap(ErrCheckpointFileReadFailed, e.Error())
		return
	}

	data_json, err := base64.StdEncoding.DecodeString(string(data_encoded))
	if err != nil {
		e = err
		return
	}

	// initialize the pointer to the Checkpoint struct
	cpoint = &Checkpoint{}
	// unmarshal the JSON data into the Checkpoint struct
	e = json.Unmarshal(data_json, cpoint)
	if e != nil {
		e = errors.Wrap(e, ErrCheckpointDataUnmarshalFailed.Error())
		return
	}
	logger.Info().Msgf("retrieved scan checkpoint data from file: %s", file_path)

	return
}

// CheckpointSet() function saves the Checkpoint data to the checkpoint file for the
// given repository and (optional) commit ID. Returns a non-nil error if unable to
// write the Checkpoint data to the expected file path.
func CheckpointSet(ctx context.Context, work_dir, repo_url, commit_id string, c *Checkpoint) (e error) {
	logger := zerolog.Ctx(ctx)
	var file *os.File
	var file_path string
	file_path, e = getCheckpointPath(work_dir, repo_url, commit_id)
	if e != nil {
		e = errors.Wrap(e, ErrMsgCheckpointSaveFailed)
		return
	}
	// attempt to open the checkpoint file
	file, e = openCheckpointFile(work_dir, repo_url, commit_id)
	if e != nil {
		// create the checkpoint file if it does not exist
		file, e = createCheckpointFile(work_dir, repo_url, commit_id)
		if e != nil {
			e = errors.Wrap(e, ErrMsgCheckpointSaveFailed)
			return
		}
		logger.Debug().Msgf("created scan checkpoint file: %s", file_path)
	}

	// marshal the Checkpoint struct into JSON bytes
	data_json, err := json.Marshal(c)
	if err != nil {
		e = errors.Wrap(err, ErrMsgCheckpointSaveFailed)
	}
	data_encoded := base64.StdEncoding.EncodeToString(data_json)

	// truncate the file to ensure it is empty before writing new data
	err = file.Truncate(0)
	if err != nil {
		e = errors.Wrap(err, ErrMsgCheckpointSaveFailed)
	}

	// write the base64-encoded JSON data to the file
	_, err = file.WriteString(data_encoded)
	if err != nil {
		e = errors.Wrap(err, ErrMsgCheckpointSaveFailed)
	}
	logger.Info().Msgf("saved scan checkpoint to file: %s", file_path)

	return
}

// createCheckpointFile() function is used to create a new checkpoint file for
// the given repository URL and commit ID. Returns a non-nil error if the file
// creation fails.
func createCheckpointFile(work_dir, repo_url, commit_id string) (file *os.File, e error) {
	var path string
	// get the expected path of the checkpoint file
	path, e = getCheckpointPath(work_dir, repo_url, commit_id)
	if e != nil {
		e = errors.Wrap(e, ErrMsgCheckpointSaveFailed)
		return
	}
	// create the parent directories as needed
	if e = os.MkdirAll(filepath.Dir(path), os.ModePerm); e != nil {
		e = errors.Wrap(e, ErrMsgCheckpointSaveFailed)
		return
	}

	// create the file if it does not exist
	file, err := os.Create(path)
	if err != nil {
		e = errors.Wrap(err, ErrMsgCheckpointSaveFailed)
		return
	}
	return
}

// getCheckpointPath() function is used to get the expected filesystem path of
// the checkpoint file for a given repository URL and commit ID, where the
// commit ID is optional. Returns a non-nil error if any required input is
// empty or if the path lookup fails.
func getCheckpointPath(work_dir, repo_url, commit_id string) (path string, e error) {

	if work_dir == "" {
		e = errors.Wrap(ErrCheckpointPathLookupFailed, "work_dir is empty")
		return
	}
	if repo_url == "" {
		e = errors.Wrap(ErrCheckpointPathLookupFailed, "repo_url is empty")
		return
	}

	var org_name string
	org_name, e = nogit.ParseOrgNameFromURL(repo_url)
	if e != nil {
		e = errors.Wrap(ErrCheckpointPathLookupFailed, e.Error())
		return
	}
	var repo_name string
	repo_name, e = nogit.ParseRepoNameFromURL(repo_url)
	if e != nil {
		e = errors.Wrap(ErrCheckpointPathLookupFailed, e.Error())
		return
	}

	// use the org_name and repo_name as the base name of the file
	name_list := []string{org_name, repo_name}
	// append the commit_id to the file name if it is not empty
	if commit_id != "" {
		name_list = append(name_list, commit_id)
	}
	file_name := strings.Join(name_list, "_") + CheckpointFileExtension
	path_list := []string{work_dir, cfg.WorkDirCheckpoints, file_name}
	path = strings.Join(path_list, "/")
	return
}

// openCheckpointFile() function is used to open the checkpoint file from its
// expected filesystem path.
func openCheckpointFile(work_dir, repo_url, commit_id string) (file *os.File, e error) {
	path, e := getCheckpointPath(work_dir, repo_url, commit_id)
	if e != nil {
		return
	}
	file, e = os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if e != nil {
		e = errors.Wrap(ErrCheckpointFileOpenFailed, e.Error())
		return
	}
	return
}
