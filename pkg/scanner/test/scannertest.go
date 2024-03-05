package test

import (
	"context"
	"errors"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
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

	scanner, err := scanner.NewScanner(
		ctx,
		&config.Git,
		memory.NewMemoryResultRecordIO(ctx),
	)
	if err != nil {
		e = err
		return
	}

	chan_scan_errors := make(chan error)
	chan_requests := make(chan rrr.Request)
	chan_responses := make(chan rrr.Response)
	dry_run_detector := dryrun.NewDryRunPhiDetector()

	go scanner.Scan(chan_scan_errors, chan_requests, chan_responses)
	go dry_run_detector.Run(ctx, chan_requests, chan_responses)

	// wait for an error to be returned from the scanner
	e = <-chan_scan_errors
	if e != nil {
		return
	}

	return
}
