package dryrun

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

// TestNewDryRunPhiDetector() unit test function tests the
// NewDryRunPhiDetector() function.
func TestNewDryRunPhiDetector(t *testing.T) {
	t.Parallel()

	d := NewDryRunPhiDetector()
	assert.NotNil(t, d)
}

// TestDryRunPhiDetector_Run() unit test function tests the
// Run() method of the DryRunPhiDetector struct.
func TestDryRunPhiDetector_Run(t *testing.T) {
	t.Parallel()

	d := &DryRunPhiDetector{}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	chan_requests_in := make(chan rrr.Request)
	chan_responses_out := make(chan rrr.Response)

	go d.Run(ctx, chan_requests_in, chan_responses_out)

	request_1 := rrr.Request{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "request-1",
			Commit: rrr.MetadataRequestResponseCommit{
				ID: "commit-1",
			},
			Object: rrr.MetadataRequestResponseObject{
				ID: "object-1",
			},
			Repository: rrr.MetadataRequestResponseRepository{
				ID: "repository-1",
			},
		},
	}
	request_2 := rrr.Request{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "request-2",
			Commit: rrr.MetadataRequestResponseCommit{
				ID: "commit-2",
			},
			Object: rrr.MetadataRequestResponseObject{
				ID: "object-2",
			},
			Repository: rrr.MetadataRequestResponseRepository{
				ID: "repository-2",
			},
		},
	}
	request_3 := rrr.Request{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "request-3",
			Commit: rrr.MetadataRequestResponseCommit{
				ID: "commit-3",
			},
			Object: rrr.MetadataRequestResponseObject{
				ID: "object-3",
			},
			Repository: rrr.MetadataRequestResponseRepository{
				ID: "repository-3",
			},
		},
	}
	// define the expected result common to both responses
	expectedResult := rrr.Result{
		Category:        DryRunCategory,
		ConfidenceScore: DryRunConfidenceScore,
		Length:          DryRunLength,
		Offset:          DryRunOffset,
		Service:         DryRunService,
		Subcategory:     DryRunSubcategory,
		Text:            DryRunText,
	}

	chan_requests_in <- request_1
	response_1 := <-chan_responses_out
	assert.Equal(t, []rrr.Result{expectedResult}, response_1.Results)
	assert.Equal(t, request_1.MetadataRequestResponse, response_1.MetadataRequestResponse)

	chan_requests_in <- request_2
	response_2 := <-chan_responses_out
	assert.Equal(t, []rrr.Result{expectedResult}, response_2.Results)
	assert.Equal(t, request_2.MetadataRequestResponse, response_2.MetadataRequestResponse)

	chan_requests_in <- request_3
	response_3 := <-chan_responses_out
	assert.Equal(t, []rrr.Result{expectedResult}, response_3.Results)
	assert.Equal(t, request_3.MetadataRequestResponse, response_3.MetadataRequestResponse)

	// cancel the context to stop the DryRunPhiDetector() function
	cancel()
	// the response channel should be closed
	_, ok := <-chan_responses_out
	assert.False(t, ok, "attempt to read from closed channel should return false")
}
