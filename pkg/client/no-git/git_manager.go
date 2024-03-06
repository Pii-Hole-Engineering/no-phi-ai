package nogit

import (
	"context"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
)

type GitRepoCloner interface {
	CloneRepo(repo_url string) (*git.Repository, error)
}

// GitManager struct provides a management wrapper for interacting with raw
// git repositories using the go-git library. Provides methods for cloning a
// repository and scanning for PHI/PII by recursively walking the
// repository's file tree.
type GitManager struct {
	GitRepoCloner

	config *cfg.GitConfig
	ctx    context.Context
	logger *zerolog.Logger
}

// NewGitManager returns a new GitManager instance for cloning, scanning, and
// other interactions with git repositories.
func NewGitManager(config *cfg.GitConfig, ctx context.Context) *GitManager {
	return &GitManager{
		config: config,
		ctx:    ctx,
		logger: zerolog.Ctx(ctx),
	}
}

// CleanRepo() method cleans the local filesystem by removing the directory created
// by the CloneRepo() method, which forces a fresh clone of the repository on the
// next call to CloneRepo().
func (gm *GitManager) CleanRepo(repo_url string) (e error) {
	var clone_dir string
	clone_dir, e = gm.getRepoCloneDir(repo_url)
	if e != nil {
		return
	}

	// return a nil error if the clone_dir does not exist
	if _, e = os.Stat(clone_dir); os.IsNotExist(e) {
		gm.logger.Debug().Msgf("skiping clean of non-existent repo clone dir %s", clone_dir)
		e = nil
		return
	}

	gm.logger.Info().Msgf("cleaning git repo clone directory : %s", clone_dir)
	if e = os.RemoveAll(clone_dir); e != nil {
		gm.logger.Error().Err(e).Msgf("error removing repo clone dir %s", clone_dir)
		e = errors.Wrapf(e, "error removing repo clone dir %s", clone_dir)
		return
	}

	return
}

// CloneRepo() method clones the repository specified by the repo_url to a
// subdirectory of the configured gm.config.WorkDir.
func (gm *GitManager) CloneRepo(repo_url string) (*git.Repository, error) {

	var key_err error
	var auth_method transport.AuthMethod
	auth_method, key_err = gm.getAuthMethod()
	if key_err != nil {
		return nil, key_err
	}

	clone_dir, dir_err := gm.getRepoCloneDir(repo_url)
	if dir_err != nil {
		return nil, dir_err
	}

	clone_options := &git.CloneOptions{
		//Progress: os.Stdout,
		URL: repo_url,
	}
	if auth_method != nil {
		clone_options.Auth = auth_method
	}
	gm.logger.Debug().Msgf("cloning git repo from %s to %s", repo_url, clone_dir)
	repo, err := git.PlainCloneContext(gm.ctx, clone_dir, false, clone_options)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			gm.logger.Info().Msgf("git repo already cloned : opening from %s", clone_dir)
			return git.PlainOpen(clone_dir)
		} else {
			gm.logger.Error().Err(err).Msgf("failed to clone git repo from %s", repo_url)
			return nil, err
		}
	}
	gm.logger.Info().Msgf("cloned git repo to %s", clone_dir)

	return repo, nil
}

// GetContext() method returns the context.Context associated with the GitManager.
func (gm *GitManager) GetContext() context.Context {
	return gm.ctx
}

// getAuthMethod() method returns the appropriate transport.AuthMethod for the
// given repo_url based on the configuration provided to the GitManager.
func (gm *GitManager) getAuthMethod() (transport.AuthMethod, error) {
	// use the provided config values to determine which auth method to use
	//
	// TODO : also use the repo_url to determine which auth method to use
	if gm.config.Auth.SSHKeyPath != "" {
		// use SSH key auth if configured
		return gm.getAuthMethodPublicKey()
	} else if gm.config.Auth.Token != "" {
		// TODO : implement token auth
		return nil, nil
	} else {
		return nil, errors.New("failed to get auth method due to missing config")
	}
}

// getAuthMethodPublicKey() method returns a transport.AuthMethod using the
// configured, local SSH key for authentication via git protocol over SSH.
func (gm *GitManager) getAuthMethodPublicKey() (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	sshPath := gm.config.Auth.SSHKeyPath
	sshKey, _ := os.ReadFile(sshPath)
	publicKey, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

// getRepoCloneDir() method is used to get the directory where a git repository
// will be cloned by this GitManager instance.
func (gm *GitManager) getRepoCloneDir(repoURL string) (string, error) {
	repoName, err := ParseRepoNameFromURL(repoURL)
	if err != nil {
		return "", err
	}
	cloneDir := gm.config.WorkDir + "/" + repoName
	return cloneDir, nil
}
