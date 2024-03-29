package scanner

import (
	"context"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	gitmemory "github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/memory"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/tracker"
)

var (
	test_context           = context.Background()
	test_failed_msg        = "failed test : %s"
	test_log_level         = "trace"
	test_repo_name         = "learn"
	test_repo_org          = "BlackIQ"
	test_repo_url          = "git@github.com:" + test_repo_org + "/" + test_repo_name + ".git"
	test_work_dir          = "/tmp/no-phi-ai/test/pkg/scanner"
	test_valid_config_func = func() *cfg.Config {
		c := cfg.NewDefaultConfig()
		c.App.Log.Level = test_log_level
		c.AzureAI.AuthKey = "test-auth-key"
		c.AzureAI.DryRun = true
		c.AzureAI.Service = "test-service"
		c.Git.Auth.Token = "test-token"
		c.Git.WorkDir = test_work_dir
		return c
	}
	test_valid_git_config_func = func() *cfg.GitConfig {
		c := test_valid_config_func()
		return &c.Git
	}
)

// TestNewScanner unit test function tests the NewScanner() function.
func TestNewScanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		config_func  func() *cfg.GitConfig
		ctx          context.Context
		err_expected bool
		name         string
	}{
		{
			config_func: func() *cfg.GitConfig {
				return &cfg.GitConfig{}
			},
			ctx:          test_context,
			err_expected: false,
			name:         "Config_Empty",
		},
		{
			config_func: func() *cfg.GitConfig {
				c := cfg.NewDefaultConfig()
				return &c.Git
			},
			ctx:          test_context,
			err_expected: false,
			name:         "Config_Default",
		},
		{
			config_func: func() *cfg.GitConfig {
				c := cfg.NewDefaultConfig()
				c.Git.Auth.Token = "test-token"
				c.Git.WorkDir = test_work_dir
				return &c.Git
			},
			ctx:          test_context,
			err_expected: false,
			name:         "Config_Missing_AzureAIAuthKey",
		},
		{
			config_func:  test_valid_git_config_func,
			ctx:          test_context,
			err_expected: false,
			name:         "Config_Valid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.config_func()
			s, err := NewScanner(
				test.ctx,
				config,
				memory.NewMemoryResultRecordIO(test_context),
			)

			if test.err_expected {
				assert.Errorf(t, err, test_failed_msg, test.name)
				assert.Nilf(t, s, test_failed_msg, test.name)
				return
			}
			assert.NoErrorf(t, err, test_failed_msg, test.name)
			if !assert.NotNilf(t, s, test_failed_msg, test.name) {
				assert.FailNowf(t, "failed to create scanner : %s", err.Error())
			}
			assert.NotEqualf(t, "", s.ID, test_failed_msg, test.name)
			if !assert.NotNil(t, s.chan_commits, "Scanner.chan_commits should not be nil") {
				assert.FailNow(t, "Scanner.chan_commits should not be nil")
			}
			if !assert.NotNil(t, s.chan_errors, "Scanner.chan_errors should not be nil") {
				assert.FailNow(t, "Scanner.chan_errors should not be nil")
			}
			if !assert.NotNil(t, s.chan_requests, "Scanner.chan_requests should not be nil") {
				assert.FailNow(t, "Scanner.chan_requests should not be nil")
			}
		})
	}
}

// TestScanner_Scan() unit test function tests the Scan() method of a new Scanner.
func TestScanner_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		config_func       func() *cfg.GitConfig
		ctx               context.Context
		err_chan          chan error
		err_expected      error
		name              string
		repo_err_expected error
		repo_func         func(ctx context.Context, repo_url string, c *cfg.GitConfig) (*git.Repository, error)
		repo_url          string
		req_chan          chan<- rrr.Request
		resp_chan         <-chan rrr.Response
	}{
		{
			config_func: func() *cfg.GitConfig {
				config := test_valid_config_func()
				config.Git.Scan.Repositories = []string{test_repo_url}
				return &config.Git
			},
			ctx:               test_context,
			err_chan:          make(chan error),
			err_expected:      nil,
			name:              "Scanner_Run_Repository_Init",
			repo_err_expected: nil,
			repo_func: func(ctx context.Context, repo_url string, c *cfg.GitConfig) (*git.Repository, error) {
				// initialize the bare *git.Repository
				return git.Init(gitmemory.NewStorage(), nil)
			},
			repo_url:  test_repo_url,
			req_chan:  make(chan<- rrr.Request),
			resp_chan: make(<-chan rrr.Response),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test_config := test.config_func()
			s, s_err := NewScanner(
				test.ctx,
				test_config,
				memory.NewMemoryResultRecordIO(test_context),
			)
			if !assert.NoErrorf(t, s_err, test_failed_msg, test.name) {
				assert.FailNowf(t, "failed to create scanner : %s", s_err.Error())
			}

			test_repository, test_repository_err := test.repo_func(test.ctx, test.repo_url, test_config)
			if test.repo_err_expected != nil {
				assert.ErrorContains(t, test_repository_err, test.repo_err_expected.Error())
				return
			} else {
				assert.NoError(t, test_repository_err)
			}

			go s.Scan(ScanInput{
				ChanErrorsSend:      test.err_chan,
				ChanRequestSend:     test.req_chan,
				ChanResponseReceive: test.resp_chan,
				RepoID:              test.repo_url,
				Repository:          test_repository,
			})

			if test.err_expected != nil {
				err := <-test.err_chan
				assert.ErrorContainsf(t, err, test.err_expected.Error(), test_failed_msg, test.name)
			}

			// TODO
		})
	}
}

// TestScanner_processRequests() unit test function tests the
// processRequests method of the Scanner object type.
func TestScanner_processRequests(t *testing.T) {
	t.Parallel()
	// create a new Scanner instance
	s, s_err := NewScanner(
		test_context,
		test_valid_git_config_func(),
		memory.NewMemoryResultRecordIO(test_context),
	)
	if !assert.NoErrorf(t, s_err, test_failed_msg, "ProcessRequests") {
		assert.FailNowf(t, "failed to create scanner : %s", s_err.Error())
	}

	// create input and output channels
	chan_quit_in := make(chan struct{})
	chan_requests_in := make(chan rrr.Request)
	chan_requests_out := make(chan<- rrr.Request)
	chan_errors_out := make(chan error)

	// start the requests processor
	go s.processRequests(chan_quit_in, chan_requests_in, chan_requests_out, chan_errors_out)

	chan_requests_in <- rrr.Request{}
	err2 := <-chan_errors_out
	assert.Equal(t, ErrProcessRequestNoID, err2)

	// close the input channels to stop goroutines
	close(chan_requests_in)
	close(chan_quit_in)

	// wait for the requests processor to finish
	time.Sleep(time.Millisecond) // Sleep for a short duration to allow goroutine to exit
}

// TestScanner_processResponse unit test function tests the processResponse method of the Scanner object type.
func TestScanner_processResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedErr  error
		name         string
		responseFunc func() rrr.Response
	}{
		{
			expectedErr: nil,
			name:        "Scanner_processResponse_Pass_1",
			responseFunc: func() rrr.Response {
				request, request_err := rrr.NewRequest(rrr.NewRequestInput{
					CommitID: "commit_id",
					Length:   len("test_text_example"),
					ObjectID: "object_id",
					Offset:   0,
					RepoID:   test_repo_url,
					Text:     "test_text_example",
				})
				if !assert.NoError(t, request_err) {
					assert.FailNow(t, "failed to create test request and response")
				}
				response := rrr.NewResponse(&request)
				return response
			},
		},
		{
			expectedErr: ErrProcessResponseNoID,
			name:        "Scanner_processResponse_Fail_1",
			responseFunc: func() rrr.Response {
				request, request_err := rrr.NewRequest(rrr.NewRequestInput{
					CommitID: "commit_id",
					Length:   len("test_text_example"),
					ObjectID: "object_id",
					Offset:   0,
					RepoID:   test_repo_url,
					Text:     "test_text_example",
				})
				if !assert.NoError(t, request_err) {
					assert.FailNow(t, "failed to create test request and response")
				}
				response := rrr.NewResponse(&request)
				// delete the response ID
				response.ID = ""
				return response
			},
		},
	}

	s, s_err := NewScanner(
		test_context,
		test_valid_git_config_func(),
		memory.NewMemoryResultRecordIO(test_context),
	)
	if !assert.NoError(t, s_err) {
		assert.FailNow(t, "failed to create scanner")
	}

	// initialize the bare *git.Repository
	repository, init_err := git.Init(gitmemory.NewStorage(), nil)
	assert.NoError(t, init_err)

	assert.NotNil(t, repository, "scanner repository pointer should not be nil")
	s.repository = repository
	// is_scan_complete must be set in order to ensure that the
	// processResponse method does not block indefinitely
	s.is_scan_complete = true

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			chan_errors_out := make(chan error)

			response := test.responseFunc()
			go s.processResponse(response, chan_errors_out)

			if test.expectedErr != nil {
				err := <-chan_errors_out
				assert.EqualError(t, err, test.expectedErr.Error())
				return
			}

			// sleep for a short duration to allow the response to be processed
			time.Sleep(time.Millisecond * 100)

			req_key_data, req_key_exists := s.TrackerRequests.Get(response.ID)
			assert.Truef(t, req_key_exists, "failed to find response ID in requests tracker : ID=%s", response.ID)
			assert.Equal(t, req_key_data.Code, tracker.KeyCodeComplete)

			commit_key_data, commit_key_exists := s.TrackerCommits.Get(response.Commit.ID)
			assert.Truef(t, commit_key_exists, "failed to find commit ID in commits tracker : ID=%s", response.Commit.ID)
			assert.Contains(t, commit_key_data.Children, response.Object.ID)
			assert.Equal(t, commit_key_data.Code, tracker.KeyCodeComplete)
			assert.Equal(t, commit_key_data.State, tracker.KeyStateComplete)

			file_key_data, file_key_exists := s.TrackerFiles.Get(response.Object.ID)
			assert.Truef(t, file_key_exists, "failed to find file ID in files tracker : ID=%s", response.Object.ID)
			assert.Contains(t, file_key_data.Children, response.ID)
			assert.Equal(t, file_key_data.Code, tracker.KeyCodeComplete)
			assert.Equal(t, file_key_data.State, tracker.KeyStateComplete)
		})
	}
}

// TestScanner_processResponses() unit test function tests the
// processResponses() method of the Scanner object type.
func TestScanner_processResponses(t *testing.T) {
	t.Parallel()
	// create a new Scanner instance
	s, s_err := NewScanner(
		test_context,
		test_valid_git_config_func(),
		memory.NewMemoryResultRecordIO(test_context),
	)
	if !assert.NoErrorf(t, s_err, test_failed_msg, "ProcessResponses") {
		assert.FailNowf(t, "failed to create scanner : %s", s_err.Error())
	}

	// create input and output channels
	chan_quit := make(chan struct{})
	chan_responses_in := make(chan rrr.Response)
	chan_errors_out := make(chan error)

	// start the response processor
	go s.processResponses(chan_quit, chan_responses_in, chan_errors_out)

	chan_responses_in <- rrr.NewResponse(&rrr.Request{})
	err2 := <-chan_errors_out
	assert.Equal(t, ErrProcessResponseNoID, err2)

	// close the input channel to stop the response processor
	close(chan_responses_in)

	// wait for the response processor to finish
	time.Sleep(time.Millisecond) // Sleep for a short duration to allow goroutine to exit
}

// TestScanner_scanRepository() unit test function tests the scanRepository()
// method of a new Scanner.
func TestScanner_scanRepository(t *testing.T) {
	t.Parallel()

	// initialize the bare *git.Repository
	repository, init_err := git.Init(gitmemory.NewStorage(), nil)
	assert.NoError(t, init_err)

	tests := []struct {
		config_func  func() *cfg.GitConfig
		ctx          context.Context
		err_chan     chan error
		err_expected error
		name         string
	}{
		{
			config_func:  test_valid_git_config_func,
			ctx:          test_context,
			name:         "Scanner_Scan_Panic_Channel_Nil",
			err_chan:     nil,
			err_expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.config_func()
			s, s_err := NewScanner(
				test.ctx,
				config,
				memory.NewMemoryResultRecordIO(test_context),
			)
			if !assert.NoErrorf(t, s_err, test_failed_msg, test.name) {
				assert.FailNowf(t, "failed to create scanner : %s", s_err.Error())
			}

			if test.err_chan == nil {
				assert.Panics(t, func() {
					s.scanRepository(
						"test_repo_url",
						repository,
						nil,
						make(chan struct{}),
					)
				})
				return
			}
			go s.scanRepository(
				"test_repo_url",
				repository,
				test.err_chan,
				make(chan struct{}),
			)
			t.Logf("waiting for error on channel")
			err := <-test.err_chan

			if test.err_expected == nil {
				assert.NoErrorf(t, err, test_failed_msg, test.name)
			} else {
				assert.ErrorContainsf(t, err, test.err_expected.Error(), test_failed_msg, test.name)
			}
		})
	}
}
