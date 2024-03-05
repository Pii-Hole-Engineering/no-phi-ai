package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	// assert that expected (default) values match what is set in the config
	assert.Equal(t, DefaultAppMode, config.App.Mode)
	assert.Equal(t, DefaultAppName, config.App.Name)
	assert.Exactly(t, false, config.App.Log.ConsoleEnable)
	assert.Exactly(t, false, config.App.Log.ConsolePretty)
	assert.Equal(t, "", config.App.Log.File)
	assert.Equal(t, DefaultAppLogLevel, config.App.Log.Level)
	assert.Equal(t, DefaultAppUserAgent, config.App.UserAgent)
	assert.Equal(t, DefaultAzureAIShowStats, config.AzureAI.ShowStats)
	assert.Equal(t, DefaultCommandRun, config.Command.Run)
	assert.Equal(t, DefaultScanFileExtensions, config.Git.Scan.Extensions)
	assert.Equal(t, DefaultMaxRequestChunkSize, config.Git.Scan.Limits.MaxRequestChunkSize)
	assert.Equal(t, DefaultMaxRequestsOutstanding, config.Git.Scan.Limits.MaxRequestsOutstanding)
	assert.Equal(t, DefaultCommandWorkDir, config.Git.WorkDir)
	assert.Equal(t, DefaultGitHubV3APIURL, config.GitHub.V3APIURL)
	assert.Equal(t, DefaultServerAddress, config.Server.Address)
	assert.Equal(t, DefaultServerPort, config.Server.Port)
	assert.Equal(t, DefaultRateLimit, config.Server.RateLimit)
	assert.Exactly(t, false, config.AzureAI.DryRun)

	// assert that required values are not set by default
	assert.Equal(t, "", config.AzureAI.AuthKey)
	assert.Equal(t, "", config.AzureAI.Service)
	assert.Equal(t, "", config.Git.Auth.SSHKeyPath)
	assert.Equal(t, "", config.Git.Auth.Token)
	assert.Exactly(t, int64(0), config.GitHub.App.IntegrationID)
	assert.Equal(t, "", config.GitHub.App.PrivateKey)
	assert.Equal(t, "", config.GitHub.App.WebhookSecret)
}
