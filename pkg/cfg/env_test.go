package cfg

import (
	"testing"
)

func TestGetAppEnvVars(t *testing.T) {
	expected := []string{
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

	result := GetAppEnvVars()

	if len(result) != len(expected) {
		t.Errorf("Expected %d environment variables, but got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected environment variable at index %d to be %s, but got %s", i, v, result[i])
		}
	}
}
