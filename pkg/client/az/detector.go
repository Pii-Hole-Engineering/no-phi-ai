package az

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

// DocumentRequestWrapper struct is a wrapper for an input rrr.Request and a
// Document created from the ID and text of the rrr.Request. Allows for tracking
// results back to the original request such that a response can be sent back
// with full context/metadata.
type DocumentRequestWrapper struct {
	Document *Document
	Request  *rrr.Request
}

// AzAiLanguagePhiDetector struct type is a wrapper for the Run() method.
type AzAiLanguagePhiDetector struct {
	ai *EntityDetectionAI
}

// NewAzAiLanguagePhiDetector() function returns a new AzAiLanguagePhiDetector instance.
func NewAzAiLanguagePhiDetector(ai *EntityDetectionAI) *AzAiLanguagePhiDetector {
	return &AzAiLanguagePhiDetector{ai: ai}
}

// Run() method listens for requests, performs no operations other than
// translating the rrr.Request to a new rrr.Response, and sends responses
// using the provided channels.
//
// Useful for testing the performance of the Scanner in generating requests
// and processing responses for chunks of text contained within the files and
// commits of the scanned repositories.
func (detector *AzAiLanguagePhiDetector) Run(
	ctx context.Context,
	chan_requests_in <-chan rrr.Request,
	chan_responses_out chan<- rrr.Response,
) {
	defer close(chan_responses_out)

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("started dry run detector")
	defer logger.Info().Msg("finished dry run detector")

	document_requests := make([]DocumentRequestWrapper, 0)

	processDocumentRequests := func() {
		if len(document_requests) == 0 {
			return
		}
		var batch_size int = RequestDocumentLimit
		if len(document_requests) < RequestDocumentLimit {
			batch_size = len(document_requests)
		}
		var document_requests_pending []DocumentRequestWrapper
		document_requests_pending, document_requests = document_requests[:batch_size], document_requests[batch_size:]

		documents := make([]Document, 0)
		for _, document_request := range document_requests_pending {
			documents = append(documents, *document_request.Document)
		}
		// create a new PiiEntityRecognitionRequest to send to AZ API
		pii_request := NewPiiEntityRecognitionRequest(documents)
		// send the request to AZ API and await the results
		pii_results, err := detector.ai.requestAiResponse(ctx, pii_request)
		if err != nil {
			// send an error via the response channel
			// TODO
		}
		// split the pii_results into individual responses
		for _, document_response := range pii_results.Results.Documents {
			// find the original request for the document response
			var original_request *rrr.Request
			for _, document_request := range document_requests_pending {
				if document_request.Request.ID == document_response.ID {
					original_request = document_request.Request
					// convert and send the response to output channel
					chan_responses_out <- convertDocumentResponseToResponse(
						detector.ai.endpoint,
						original_request,
						&document_response,
					)
					break
				}
			}
		}
		// handle any remnants from document_requests_pending
		// by sending an error response for each (orphaned) request
		// TODO
	}

	timer := time.NewTimer(RequestTimerDuration)

	for {
		select {
		case <-ctx.Done():
			logger.Warn().Msg("stopping dry run detector : context done")
			// exit the function when the context is done
			return
		case request := <-chan_requests_in:
			document_requests = append(document_requests, wrapDocumentRequest(&request))
			if len(document_requests) >= RequestDocumentLimit {
				processDocumentRequests()
			}
			// stop and reset the timer
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(RequestTimerDuration)
		case <-timer.C:
			processDocumentRequests()
			continue
		}
	}
}

// convertDocumentResponseToResponse() function converts from a DocumentResponse
// struct to a rrr.Response struct, using the original rrr.Request struct to
// initialize the new rrr.Response struct.
func convertDocumentResponseToResponse(
	endpoint string,
	request *rrr.Request,
	doc_response *DocumentResponse,
) rrr.Response {
	// use the original request to initialize the new response
	response := rrr.NewResponse(request)

	// convert each Entity to an rrr.Result
	for _, entity := range doc_response.Entities {
		result := convertEntitytToResult(endpoint, entity)
		// append the converted rrr.Result to the results slice
		response.Results = append(response.Results, result)
	}

	return response
}

// convertEntitytToResult() function converts from an Entity struct to a
// rrr.Result struct, using the provided endpoint string to set the Service
// field of the rrr.Result struct.
func convertEntitytToResult(endpoint string, entity Entity) rrr.Result {
	return rrr.Result{
		Category:        entity.Category,
		ConfidenceScore: entity.ConfidenceScore,
		Length:          entity.Length,
		Offset:          entity.Offset,
		Service:         endpoint,
		Subcategory:     entity.Subcategory,
		Text:            entity.Text,
	}
}

func wrapDocumentRequest(request *rrr.Request) DocumentRequestWrapper {
	document := NewDocument(request.ID, request.Text, "")
	return DocumentRequestWrapper{
		Document: &document,
		Request:  request,
	}
}
