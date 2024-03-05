package az

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
)

// EntityDetectionAI struct provides methods for (1) sending requests to
// detect entities of interest in natural language documents and (2) processing
// responses from the Azure AI Language service.
type EntityDetectionAI struct {
	client     *http.Client
	confidence float64
	dryRun     bool
	endpoint   string
	key        string
}

// NewEntityDetectionAI() function requires the Azure service host and
// authentication key as inputs, and returns an initialized EntityDetectionAI.
func NewEntityDetectionAI(c *cfg.Config) (*EntityDetectionAI, error) {
	if len(c.AzureAI.AuthKey) == 0 {
		err := errors.New("EntityDetectionAI requires a valid authentication key")
		return nil, err
	}
	if len(c.AzureAI.Service) == 0 {
		err := errors.New("EntityDetectionAI requires a valid service address")
		return nil, err
	}

	// add query param(s) to the endpoint for the Azure AI Languge service
	//  e.g. to "show stats" on document analysis
	endpoint := fmt.Sprintf(
		"%s/%s",
		c.AzureAI.Service,
		getDefaultDetectionApi(c.AzureAI.ShowStats),
	)

	return &EntityDetectionAI{
		client:     &http.Client{},
		confidence: c.AzureAI.ConfidenceThreshold,
		dryRun:     c.AzureAI.DryRun,
		endpoint:   endpoint,
		key:        c.AzureAI.AuthKey,
	}, nil
}

// DetectPiiEntities() method accepts a PiiEntityRecognitionRequest as
// input and returns an error (or nil) status depending upon the
// PiiEntityRecognitionResponse from the Azure AI language service.
func (ai *EntityDetectionAI) DetectPiiEntities(
	ctx context.Context,
	request_data *PiiEntityRecognitionRequest,
) (detected bool, e error) {
	// explicitly set detected to false
	detected = false

	entity_recognition_results, entity_recognition_err := ai.requestAiResponse(ctx, request_data)
	if entity_recognition_err != nil {
		e = errors.Wrap(entity_recognition_err, "error requesting AI API response")
		return
	}

	for _, doc := range entity_recognition_results.Results.Documents {
		var entities []Entity
		for _, entity := range doc.Entities {
			if entity.ConfidenceScore >= ai.confidence {
				entities = append(entities, entity)
			}
		}
		if len(entities) > 0 {
			// set detected to true if any entities are found over the confidence threshold
			detected = true

			// TODO : remove/replace processing of entities
			result := struct {
				ID       string   `json:"id"`
				Entities []Entity `json:"entities"`
			}{
				ID:       doc.ID,
				Entities: entities,
			}
			var resultBytes []byte
			resultBytes, e = json.MarshalIndent(result, "", "  ")
			if e != nil {
				log.Ctx(ctx).Error().Msgf("error marshalling result within entity detection response: %s", e)
				return
			}
			// TODO : remove debug logging of resultBytes
			log.Ctx(ctx).Debug().Msgf("entity detection result:\n%s", string(resultBytes))
		}
	}

	return
}

// GetServiceEndpoint() method return the full URL of the service API endpoint
func (ai *EntityDetectionAI) GetServiceEndpoint() string {
	return ai.endpoint
}

func (ai *EntityDetectionAI) dryRunRespond(ctx context.Context, entity_request *PiiEntityRecognitionRequest) (*PiiEntityRecognitionResults, error) {
	// create a fake response for dry run mode
	fake_response := &PiiEntityRecognitionResults{
		Results: Results{
			Documents: make([]DocumentResponse, 0),
		},
	}

	for _, doc := range entity_request.AnalysisInput.Documents {
		fake_response.Results.Documents = append(fake_response.Results.Documents, DocumentResponse{
			ID:       doc.ID,
			Entities: make([]Entity, 0),
		})

	}

	log.Ctx(ctx).Trace().Msgf(
		"sending DRY RUN MODE response for %d documents",
		len(fake_response.Results.Documents),
	)

	return fake_response, nil
}

// requestAiResponse() method converts a PiiEntityRecognitionRequest to a
// JSON byte array and sends the request to the Azure AI Language service API,
// then converts the JSON response from the API to PiiEntityRecognitionResults.
// This method returns an error if the request or response fails.
func (ai *EntityDetectionAI) requestAiResponse(
	ctx context.Context,
	entity_request *PiiEntityRecognitionRequest,
) (*PiiEntityRecognitionResults, error) {
	var e error

	// check for dry run mode
	if ai.dryRun {
		return ai.dryRunRespond(ctx, entity_request)
	}

	entity_request_bytes, err := json.Marshal(entity_request)
	if err != nil {
		e = errors.Wrap(err, "failed to marshal PiiEntityRecognitionRequest")
		return nil, e
	}

	http_request, err := http.NewRequestWithContext(ctx, "POST", ai.endpoint, bytes.NewBuffer(entity_request_bytes))
	if err != nil {
		e = errors.Wrap(err, "failed creating HTTP request to Azure AI Language service")
		return nil, e
	}

	// set the required headers for the HTTP request before sending
	ai.setHttpRequestHeaders(http_request)

	log.Ctx(ctx).Trace().Msgf(
		"requesting entity recognition for %d documents",
		len(entity_request.AnalysisInput.Documents),
	)

	// send the HTTP request to the Azure AI Language service API
	http_response, err := ai.client.Do(http_request)
	if err != nil {
		e = errors.Wrap(err, "failed sending HTTP request to Azure AI Language service")
		return nil, e
	}
	defer http_response.Body.Close()

	// read all bytes from the HTTP response body
	http_response_body, err := io.ReadAll(http_response.Body)
	if err != nil {
		e = errors.Wrap(err, "failed reading HTTP response from Azure AI Language service")
		return nil, e
	}

	var entity_recognition_results PiiEntityRecognitionResults
	// unmarshal the bytes from the response body into a PiiEntityRecognitionResults
	if e = json.Unmarshal(http_response_body, &entity_recognition_results); e != nil {
		e = errors.Wrap(e, "failed unmarshalling response from Azure AI Language service")
		return nil, e
	}

	log.Ctx(ctx).Trace().Msgf(
		"received entity recognition results with %d document responses",
		len(entity_recognition_results.Results.Documents),
	)

	return &entity_recognition_results, nil
}

// setHttpRequestHeaders() method sets the required headers for any HTTP
// request to the Azure AI Language service API.
func (ai *EntityDetectionAI) setHttpRequestHeaders(req *http.Request) {
	req.Header.Add("Ocp-Apim-Subscription-Key", ai.key)
	req.Header.Add("Content-Type", "application/json")
}

func getDefaultDetectionApi(showStats bool) string {
	if showStats {
		return DefaultDetectionApi + ShowStatsParam
	}
	return DefaultDetectionApi
}
