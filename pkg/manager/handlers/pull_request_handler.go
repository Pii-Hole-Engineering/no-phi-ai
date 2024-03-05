package handlers

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/az"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/gh"
)

type PullRequestHandler struct {
	AI   *az.EntityDetectionAI
	GHCM *gh.ClientManager
}

func (h *PullRequestHandler) Handles() []string {
	return []string{EventTypePullRequest}
}

func (h *PullRequestHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequestEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse payload for event type="+EventTypePullRequest)
	}
	zerolog.Ctx(ctx).Debug().Msgf("%s received webhook event type=%s", h.name(), eventType)
	// TODO : remove vulnerable use of payload as unfiltered input to logging function
	zerolog.Ctx(ctx).Debug().Msgf("%s received webhook event:\n%s", h.name(), string(payload))

	// TODO

	return nil
}

// PullRequestHandler.name() method is NOT required by any interface.
func (h *PullRequestHandler) name() string {
	return "PullRequestHandler"
}
