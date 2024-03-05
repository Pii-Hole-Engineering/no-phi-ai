package cfg

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Use constant strings to declare the environment variable names used by the app.
const NOPHI_APP_LOG_LEVEL string = "NOPHI_APP_LOG_LEVEL"
const NOPHI_APP_MODE string = "NOPHI_APP_MODE"
const NOPHI_APP_NAME string = "NOPHI_APP_NAME"
const NOPHI_AZURE_AI_AUTH_KEY string = "NOPHI_AZURE_AI_AUTH_KEY"
const NOPHI_AZURE_AI_DRY_RUN string = "NOPHI_AZURE_AI_DRY_RUN"
const NOPHI_AZURE_AI_SERVICE string = "NOPHI_AZURE_AI_SERVICE"
const NOPHI_AZURE_AI_SHOW_STATS string = "NOPHI_AZURE_AI_SHOW_STATS"
const NOPHI_COMMAND_RUN = "NOPHI_COMMAND_RUN"
const NOPHI_CONFIG_PATH string = "NOPHI_CONFIG_PATH"
const NOPHI_GH_INTEGRATION_ID string = "NOPHI_GH_INTEGRATION_ID"
const NOPHI_GH_PRIVATE_KEY string = "NOPHI_GH_PRIVATE_KEY"
const NOPHI_GH_V3APIURL string = "NOPHI_GH_V3APIURL"
const NOPHI_GH_V4APIURL string = "NOPHI_GH_V4APIURL"
const NOPHI_GH_WEBHOOK_SECRET = "NOPHI_GH_WEBHOOK_SECRET"
const NOPHI_GIT_SCAN_MAX_REQUEST_CHUNK_SIZE = "NOPHI_GIT_SCAN_MAX_REQUEST_CHUNK_SIZE"
const NOPHI_GIT_WORKDIR = "NOPHI_GIT_WORKDIR"
const NOPHI_MAX_REQUESTS_OUTSTANDING = "NOPHI_MAX_REQUESTS_OUTSTANDING"
const NOPHI_SERVER_ADDRESS string = "NOPHI_SERVER_ADDRESS"
const NOPHI_SERVER_PORT string = "NOPHI_SERVER_PORT"

// GetAppEnvVars() method returns a list of environment variables used by the app.
func GetAppEnvVars() []string {
	return []string{
		NOPHI_APP_LOG_LEVEL,
		NOPHI_APP_MODE,
		NOPHI_APP_NAME,
		NOPHI_AZURE_AI_AUTH_KEY,
		NOPHI_AZURE_AI_DRY_RUN,
		NOPHI_AZURE_AI_SERVICE,
		NOPHI_AZURE_AI_SHOW_STATS,
		NOPHI_COMMAND_RUN,
		NOPHI_CONFIG_PATH,
		NOPHI_GH_INTEGRATION_ID,
		NOPHI_GH_PRIVATE_KEY,
		NOPHI_GH_V3APIURL,
		NOPHI_GH_V4APIURL,
		NOPHI_GH_WEBHOOK_SECRET,
		NOPHI_GIT_SCAN_MAX_REQUEST_CHUNK_SIZE,
		NOPHI_GIT_WORKDIR,
		NOPHI_MAX_REQUESTS_OUTSTANDING,
		NOPHI_SERVER_ADDRESS,
		NOPHI_SERVER_PORT,
	}
}

// envOverride() method overrides *Config values with values from environment variables.
func (c *Config) envOverride() error {
	if mode := os.Getenv(NOPHI_APP_MODE); mode != "" {
		c.App.Mode = mode
	}
	if name := os.Getenv(NOPHI_APP_NAME); name != "" {
		c.App.Name = name
	}
	if logLevel := os.Getenv(NOPHI_APP_LOG_LEVEL); logLevel != "" {
		c.App.Log.Level = logLevel
	}
	if azAuthKey := os.Getenv(NOPHI_AZURE_AI_AUTH_KEY); azAuthKey != "" {
		c.AzureAI.AuthKey = azAuthKey
	}
	if azDryRun := os.Getenv(NOPHI_AZURE_AI_DRY_RUN); azDryRun != "" {
		azDryRunBool, err := strconv.ParseBool(azDryRun)
		if err != nil {
			return errors.Wrap(err, "failed parsing NOPHI_AZURE_AI_DRY_RUN env var")
		}
		c.AzureAI.DryRun = azDryRunBool
	}
	if azService := os.Getenv(NOPHI_AZURE_AI_SERVICE); azService != "" {
		c.AzureAI.Service = azService
	}
	if azShowStats := os.Getenv(NOPHI_AZURE_AI_SHOW_STATS); azShowStats != "" {
		azShowStatsBool, err := strconv.ParseBool(azShowStats)
		if err != nil {
			return errors.Wrap(err, "failed parsing NOPHI_AZURE_AI_SHOW_STATS env var")
		}
		c.AzureAI.ShowStats = azShowStatsBool
	}
	if commandRun := os.Getenv(NOPHI_COMMAND_RUN); commandRun != "" {
		c.Command.Run = commandRun
	}
	if maxRequestsOutstanding := os.Getenv(NOPHI_MAX_REQUESTS_OUTSTANDING); maxRequestsOutstanding != "" {
		maxRequestsOutstandingInt, err := strconv.Atoi(maxRequestsOutstanding)
		if err == nil {
			c.Git.Scan.Limits.MaxRequestsOutstanding = maxRequestsOutstandingInt
		}
	}
	if chunkSize := os.Getenv(NOPHI_GIT_SCAN_MAX_REQUEST_CHUNK_SIZE); chunkSize != "" {
		chunkSizeInt, err := strconv.Atoi(chunkSize)
		if err != nil {
			return errors.Wrap(err, "failed parsing NOPHI_GIT_SCAN_MAX_REQUEST_CHUNK_SIZE env var")
		}
		c.Git.Scan.Limits.MaxRequestChunkSize = chunkSizeInt
	}
	if gitWorkDir := os.Getenv(NOPHI_GIT_WORKDIR); gitWorkDir != "" {
		c.Git.WorkDir = gitWorkDir
	}
	if integrationID := os.Getenv(NOPHI_GH_INTEGRATION_ID); integrationID != "" {
		integrationIDInt, err := strconv.ParseInt(integrationID, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed parsing NOPHI_GH_INTEGRATION_ID env var")
		}
		c.GitHub.App.IntegrationID = integrationIDInt
	}
	if privateKey := os.Getenv(NOPHI_GH_PRIVATE_KEY); privateKey != "" {
		c.GitHub.App.PrivateKey = privateKey
	}
	if webhookSecret := os.Getenv(NOPHI_GH_WEBHOOK_SECRET); webhookSecret != "" {
		c.GitHub.App.WebhookSecret = webhookSecret
	}
	if V3APIURL := os.Getenv(NOPHI_GH_V3APIURL); V3APIURL != "" {
		c.GitHub.V3APIURL = V3APIURL
	}
	if V4APIURL := os.Getenv(NOPHI_GH_V4APIURL); V4APIURL != "" {
		c.GitHub.V4APIURL = V4APIURL
	}
	if serverAddress := os.Getenv(NOPHI_SERVER_ADDRESS); serverAddress != "" {
		c.Server.Address = serverAddress
	}
	if serverPort := os.Getenv(NOPHI_SERVER_PORT); serverPort != "" {
		serverPortInt, err := strconv.Atoi(serverPort)
		if err != nil {
			return errors.Wrap(err, "failed parsing SERVER_PORT env var")
		}
		c.Server.Port = serverPortInt
	}

	return nil
}
