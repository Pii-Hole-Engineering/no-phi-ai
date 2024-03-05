package az

import (
	"testing"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
)

func TestNewEntityDetectionAI(t *testing.T) {
	// Test case 1: Valid service and key
	c1 := &cfg.Config{}
	c1.AzureAI.AuthKey = "valid-key"
	c1.AzureAI.Service = "https://example.com"
	engine, err := NewEntityDetectionAI(c1)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if engine == nil {
		t.Error("Expected engine to be created, but got nil")
	}

	// Test case 2: Empty key
	c2 := &cfg.Config{}
	c2.AzureAI.AuthKey = ""
	c2.AzureAI.Service = "https://example.com"
	engine, err = NewEntityDetectionAI(c2)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if engine != nil {
		t.Error("Expected engine to be nil, but got non-nil")
	}
	if err.Error() != "EntityDetectionAI requires a valid authentication key" {
		t.Errorf("Expected error message: 'EntityDetectionAI requires a valid authentication key', but got: %v", err.Error())
	}

	// Test case 3: Empty service
	c3 := &cfg.Config{}
	c3.AzureAI.AuthKey = "valid-key"
	c3.AzureAI.Service = ""
	engine, err = NewEntityDetectionAI(c3)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if engine != nil {
		t.Error("Expected engine to be nil, but got non-nil")
	}
	if err.Error() != "EntityDetectionAI requires a valid service address" {
		t.Errorf("Expected error message: 'EntityDetectionAI requires a valid service address', but got: %v", err.Error())
	}
}

// Test_getDefaultDetectionApi() unit test function tests the
// getDefaultDetectionApi() function.
func Test_getDefaultDetectionApi(t *testing.T) {
	// Test case 1: showStats is true
	showStats := true
	expectedResult := DefaultDetectionApi + ShowStatsParam
	result := getDefaultDetectionApi(showStats)
	if result != expectedResult {
		t.Errorf("Expected result: %s, but got: %s", expectedResult, result)
	}

	// Test case 2: showStats is false
	showStats = false
	expectedResult = DefaultDetectionApi
	result = getDefaultDetectionApi(showStats)
	if result != expectedResult {
		t.Errorf("Expected result: %s, but got: %s", expectedResult, result)
	}
}
