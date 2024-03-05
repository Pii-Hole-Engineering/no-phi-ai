package handlers

import (
	"context"
	"encoding/json"

	//"fmt"
	//"strings"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/az"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/gh"
)

type IssueCommentHandler struct {
	AI   *az.EntityDetectionAI
	GHCM *gh.ClientManager
}

func (h *IssueCommentHandler) Handles() []string {
	return []string{EventTypeIssueComment}
}

func (h *IssueCommentHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse payload for eventType="+EventTypeIssueComment)
	}
	zerolog.Ctx(ctx).Debug().Msgf("%s received webhook eventType=%s", h.name(), eventType)
	// TODO : remove vulnerable use of payload as unfiltered input to logging function
	//zerolog.Ctx(ctx).Debug().Msgf("%s received webhook event:\n%s", h.name(), string(payload))

	if event.GetIssue().IsPullRequest() {
		return h.handlePullRequestIssueComments(ctx, eventType, deliveryID, event)
	}
	return h.handleIssueComments(ctx, eventType, deliveryID, event)
}

// handleIssueComments() method is used for handling comment events for issues that are not pull requests.
func (h *IssueCommentHandler) handleIssueComments(ctx context.Context, eventType, deliveryID string, event github.IssueCommentEvent) (e error) {
	// TODO : handle issue comment events for regular issues
	if (len(eventType) > 0) && (len(deliveryID) > 0) {
		e = errors.New("issue comment event is NOT for a pull request")
		zerolog.Ctx(ctx).Error().Err(e).Msgf("failed to handle eventType=%s : deliveryID=%s", eventType, deliveryID)
	}

	return e
}

// handlePullRequestIssueComments() method is used for handling comment events for pull requests.
// From the GitHub API perspective, all pull requests are issues, but not all issues are pull requests.
func (h *IssueCommentHandler) handlePullRequestIssueComments(ctx context.Context, eventType, deliveryID string, event github.IssueCommentEvent) (e error) {
	zerolog.Ctx(ctx).Debug().Msg("issue comment event is for a pull request")

	// check the "action" field of the event
	eventAction := event.GetAction()
	switch eventAction {
	case "created":
		break
	default:
		zerolog.Ctx(ctx).Debug().Msgf("ignoring event action=%s for eventType=%s : deliveryID=%s", eventAction, eventType, deliveryID)
		return nil
	}

	// create a slice of documents to send to Azure AI Language service
	documents := []az.Document{}
	// add documents to the slice using data from text fields in the webhook event
	if event_comment := event.GetComment().GetBody(); event_comment != "" {
		document := az.NewDocument(event.GetComment().GetURL(), event_comment, "en")
		documents = append(documents, document)
	}
	// TODO : pull data from other text fields as potential sources of PHI/PII

	// default label value assumes there is no PHI/PII in the documents
	var issue_label = gh.LabelCleanPHI
	if len(documents) > 0 {
		// create a new request to detect PII entities in the documents
		req := az.NewPiiEntityRecognitionRequest(documents)

		zerolog.Ctx(ctx).Debug().Msgf("sending PII entity detection request for %d documents", len(documents))
		var found bool
		// send the detection request to the Azure AI Language service
		found, e = h.AI.DetectPiiEntities(ctx, req)
		if e != nil {
			zerolog.Ctx(ctx).Debug().Msg(e.Error())
			return e
		}
		if found {
			// update the issue_label value to indicate that the AI language
			// service detected PHI/PII in the documents
			issue_label = gh.LabelDirtyPHI
		}
		zerolog.Ctx(ctx).Debug().Msgf("AI scanned %d documents : %s", len(documents), issue_label)
	} else {
		e = errors.New("no documents to process in issue comment handler")
		zerolog.Ctx(ctx).Debug().Msg(e.Error())
		return e
	}

	// TODO : replace static label with values determined by ^ detection
	if err := h.GHCM.ApplyLabelForIssueComment(ctx, event, issue_label); err != nil {
		return err
	}

	return nil
}

// IssueCommentHandler.name() method is NOT required by any interface.
func (h *IssueCommentHandler) name() string {
	return "IssueCommentHandler"
}
