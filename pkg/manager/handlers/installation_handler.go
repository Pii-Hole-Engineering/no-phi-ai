package handlers

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/gh"
)

type InstallationHandler struct {
	GHCM *gh.ClientManager
}

func (h *InstallationHandler) Handles() []string {
	return []string{EventTypeInstallation}
}

func (h *InstallationHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.InstallationEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse payload for event type="+EventTypeInstallation)
	}
	zerolog.Ctx(ctx).Debug().Msgf("%s received webhook event type=%s", h.name(), eventType)
	// TODO : remove vulnerable use of payload as unfiltered input to logging function
	zerolog.Ctx(ctx).Debug().Msgf("%s received webhook event:\n%s", h.name(), string(payload))

	// TODO : do something with the GH app installation event

	return nil
}

// InstallationHandler.name() method is NOT required by any interface.
func (h *InstallationHandler) name() string {
	return "InstallationHandler"
}
