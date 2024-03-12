package test

import (
	"context"
	"errors"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	nogit "github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/no-git"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/dryrun"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/memory"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

const ScannerTestDataDir = "./testdata"
const ScannerTestRepoPath = ScannerTestDataDir + "/test-repo-1"
const ScannerTestRepoURL = "git@github.com:Pii-Hole-Engineering/test-repo-1.git"

// ScannerTestEndToEnd() function is used to run an end-to-end test of the
// scanner.Scanner using:
//   - the dryrun.DryRunPhiDetector for simulating responses and responses;
//   - the memory.MemoryResultRecordIO for storing rrr.Result records.
func ScannerTestEndToEnd(ctx context.Context, repo_url string) (e error) {
	if repo_url == "" {
		e = errors.New("repo_url is required")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	config := cfg.NewDefaultConfig()
	config.App.Log.Level = "trace"
	config.AzureAI.AuthKey = "test-auth-key"
	config.AzureAI.DryRun = true
	config.AzureAI.Service = "test-service"
	config.Git.Auth.Token = "test-token"
	config.Git.Scan.Repositories = []string{repo_url}
	config.Git.WorkDir = ScannerTestDataDir

	s, err := scanner.NewScanner(
		ctx,
		&config.Git,
		memory.NewMemoryResultRecordIO(ctx),
	)
	if err != nil {
		e = err
		return
	}

	git_manager := nogit.NewGitManager(&config.Git, ctx)

	repo_url = config.Git.Scan.Repositories[0]
	// clone the repository
	repository, repository_err := git_manager.CloneRepo(repo_url)
	if repository_err != nil {
		e = repository_err
		return
	}

	dry_run_detector := dryrun.NewDryRunPhiDetector()

	chan_scan_errors := make(chan error)
	chan_requests := make(chan rrr.Request)
	chan_responses := make(chan rrr.Response)

	go s.Scan(scanner.ScanInput{
		ChanErrorsSend:      chan_scan_errors,
		ChanRequestSend:     chan_requests,
		ChanResponseReceive: chan_responses,
		RepoID:              repo_url,
		Repository:          repository,
	})
	go dry_run_detector.Run(ctx, chan_requests, chan_responses)

	// wait for an error to be returned from the scanner
	e = <-chan_scan_errors
	if e != nil {
		return
	}

	return
}
