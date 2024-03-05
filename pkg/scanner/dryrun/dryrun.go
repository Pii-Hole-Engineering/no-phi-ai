package dryrun

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

const DryRunCategory = "dry-run_category"
const DryRunConfidenceScore = 0.0
const DryRunLength = 100
const DryRunOffset = 0
const DryRunService = "dry-run_service"
const DryRunSubcategory = "dry-run_subcategory"
const DryRunText = "dry-run_text"

// DryRunPhiDetector struct type is a wrapper for the Run() method.
type DryRunPhiDetector struct{}

// NewDryRunPhiDetector() function returns a new DryRunPhiDetector instance.
func NewDryRunPhiDetector() *DryRunPhiDetector {
	return &DryRunPhiDetector{}
}

// Run() method listens for requests, performs no operations other than
// translating the rrr.Request to a new rrr.Response, and
// sends responses using the provided channels.
//
// Useful for testing the performance of the Scanner in generating requests
// and processing responses for chunks of text contained within the files and
// commits of the scanned repositories.
func (detector *DryRunPhiDetector) Run(
	ctx context.Context,
	chan_requests_in <-chan rrr.Request,
	chan_responses_out chan<- rrr.Response,
) {
	defer close(chan_responses_out)

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("started dry run detector")
	defer logger.Info().Msg("finished dry run detector")

	for {
		select {
		case <-ctx.Done():
			logger.Warn().Msg("stopping dry run detector : context done")
			// exit the function when the context is done
			return
		case request := <-chan_requests_in:
			// create a rrr.Response from the rrr.Request metadata
			response := rrr.NewResponse(&request)
			// create a dummy rrr.Result for the rrr.Response
			result := rrr.Result{
				Category:        DryRunCategory,
				ConfidenceScore: DryRunConfidenceScore,
				Length:          DryRunLength,
				Offset:          DryRunOffset,
				Service:         DryRunService,
				Subcategory:     DryRunSubcategory,
				Text:            DryRunText,
			}
			// add the dummy Result to the Response
			response.Results = append(response.Results, result)
			// send the Response to the output channel
			chan_responses_out <- response
			logger.Debug().Msgf(
				"dry run detector processed request ID = %s : RepositoryID = %s : CommitID = %s : ObjectID = %s",
				request.ID,
				request.Repository.ID,
				request.Commit.ID,
				request.Object.ID,
			)
		}
	}
}
