package scanner

import (
	"context"
	"testing"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

var (
	test_repo_name        = "learn"
	test_repo_org         = "BlackIQ"
	test_repo_scan_config = &cfg.GitScanConfig{
		IgnoreRepositories: []string{},
		Repositories:       []string{test_repo_url},
	}
	test_repo_url = "git@github.com:" + test_repo_org + "/" + test_repo_name + ".git"
)

// TestNewScanRepository unit test function tests the NewScanRepository() function.
func TestNewScanRepository(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	channel_errors := make(chan<- error)
	channel_requests := make(chan<- rrr.Request)

	t.Run("ValidInput", func(t *testing.T) {
		// initialize the bare *git.Repository
		repository, init_err := git.Init(memory.NewStorage(), nil)
		assert.NoError(t, init_err)

		repo, err := NewScanRepository(NewScanRepositoryInput{
			ChannelErrors:   channel_errors,
			ChannelRequests: channel_requests,
			Config:          test_repo_scan_config,
			Context:         ctx,
			Repository:      repository,
			URL:             test_repo_url,
		})

		assert.NoError(t, err)
		assert.NotNil(t, repo)
		assert.NotEmpty(t, repo.ID)
		assert.Equal(t, test_repo_url, repo.ID)
		assert.Equal(t, test_repo_name, repo.Name)
		assert.Equal(t, test_repo_url, repo.URL)
		assert.Equal(t, channel_requests, repo.channel_requests)
		assert.Equal(t, ctx, repo.ctx)
		assert.NotNil(t, repo.logger)
		assert.NotNil(t, repo.repository)
		assert.NotNil(t, repo.TrackerCommits)
		assert.NotNil(t, repo.TrackerFiles)
	})

	t.Run("NilContext", func(t *testing.T) {
		// initialize the bare *git.Repository
		repository, init_err := git.Init(memory.NewStorage(), nil)
		assert.NoError(t, init_err)

		repo, err := NewScanRepository(NewScanRepositoryInput{
			ChannelErrors:   channel_errors,
			ChannelRequests: channel_requests,
			Config:          test_repo_scan_config,
			Context:         nil,
			Repository:      repository,
			URL:             test_repo_url,
		})

		assert.ErrorContains(t, err, ErrMsgScanRepositoryCreate)
		assert.ErrorContains(t, err, ErrScanRepositoryContextNil.Error())
		assert.Nil(t, repo)
	})

	t.Run("NilChannelErrors", func(t *testing.T) {
		// initialize the bare *git.Repository
		repository, init_err := git.Init(memory.NewStorage(), nil)
		assert.NoError(t, init_err)

		repo, err := NewScanRepository(NewScanRepositoryInput{
			ChannelErrors:   nil,
			ChannelRequests: channel_requests,
			Config:          test_repo_scan_config,
			Context:         ctx,
			Repository:      repository,
			URL:             test_repo_url,
		})

		assert.ErrorContains(t, err, ErrMsgScanRepositoryCreate)
		assert.ErrorContains(t, err, ErrScanRepositoryChannelErrorsNil.Error())
		assert.Nil(t, repo)
	})

	t.Run("NilChannelRequests", func(t *testing.T) {
		// initialize the bare *git.Repository
		repository, init_err := git.Init(memory.NewStorage(), nil)
		assert.NoError(t, init_err)

		repo, err := NewScanRepository(NewScanRepositoryInput{
			ChannelErrors:   channel_errors,
			ChannelRequests: nil,
			Config:          test_repo_scan_config,
			Context:         ctx,
			Repository:      repository,
			URL:             test_repo_url,
		})

		assert.ErrorContains(t, err, ErrMsgScanRepositoryCreate)
		assert.ErrorContains(t, err, ErrScanRepositoryChannelRequestsNil.Error())
		assert.Nil(t, repo)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		// initialize the bare *git.Repository
		repository, init_err := git.Init(memory.NewStorage(), nil)
		assert.NoError(t, init_err)

		repo, err := NewScanRepository(NewScanRepositoryInput{
			ChannelErrors:   channel_errors,
			ChannelRequests: channel_requests,
			Config:          test_repo_scan_config,
			Context:         ctx,
			Repository:      repository,
			URL:             "invalid-url",
		})

		assert.ErrorContains(t, err, ErrMsgScanRepositoryCreate)
		assert.Nil(t, repo)
	})
}

// TestScanRepository_GetRepository() tests the GetRepository() method
// of the ScanRepository struct.
func TestScanRepository_GetRepository(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	channel_errors := make(chan<- error)
	channel_requests := make(chan<- rrr.Request)
	name := "TestRepositoryInit"

	t.Run(name, func(t *testing.T) {
		// initialize the bare *git.Repository
		repository, init_err := git.Init(memory.NewStorage(), nil)
		assert.NoError(t, init_err)

		repo, err := NewScanRepository(NewScanRepositoryInput{
			ChannelErrors:   channel_errors,
			ChannelRequests: channel_requests,
			Config:          test_repo_scan_config,
			Context:         ctx,
			Repository:      repository,
			URL:             test_repo_url,
		})
		if !assert.NoError(t, err) || !assert.NotNil(t, repo) {
			t.FailNow()
		}

		assert.NotEmpty(t, repo.ID)

		result := repo.GetRepository()
		assert.NotNil(t, result)
	})
}
