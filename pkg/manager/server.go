package manager

import (
	"fmt"
	"net/http"
	"time"

	tollbooth "github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/palantir/go-githubapp/githubapp"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/az"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/gh"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/manager/handlers"
)

// initServer() method initializes the HTTP server and registers handlers.
func (m *Manager) initServer() {
	// setup the rate limiter for the HTTP server
	lmt := tollbooth.NewLimiter(m.config.Server.RateLimit, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetHeader("Authorization", []string{})
	lmt.SetHeaderEntryExpirationTTL(time.Hour)
	lmt.SetMessage(`{"code":429,"response":"You have reached maximum request limit. Please try again in a few seconds."}`)
	lmt.SetMessageContentType("application/json")

	// set the "mode" for the (gin) HTTP server
	gin.SetMode(gin.ReleaseMode)
	// setup a new router for the HTTP server
	router := gin.New()
	router.Use(requestid.New(
		requestid.WithGenerator(func() string {
			return uuid.NewString()
		}),
	))
	// setup structured logging for the HTTP server
	router.Use(handlers.StructuredLogger(m.logger))
	// allow the HTTP server to recover from panics
	router.Use(gin.Recovery())

	// setup an http.Handler as the event dispatcher for GitHub webhook events
	eventDispatcher, err := setupEventDispatcher(m.config)
	if err != nil {
		m.logger.Fatal().Err(err).Msg("failed to setup event handler for new Manager")
	}

	// setup routes for the HTTP server
	v1 := router.Group(cfg.RouteGroupGHv1)
	v1.POST(cfg.RouteWebhook, handlers.LimitHandler(lmt), gin.WrapH(eventDispatcher))

	// set the m.server field to a new http.Server instance
	m.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", m.config.Server.Address, m.config.Server.Port),
		Handler: router,
	}
}

// logRoutes() function logs the HTTP routes registered with the gin Router.
func (m *Manager) logRoutes() {
	for _, route := range m.server.Handler.(*gin.Engine).Routes() {
		m.logger.Info().Msgf("serving endpoint -> %s %s%s", route.Method, m.server.Addr, route.Path)
	}
}

// runServer() method auto-loads configuration from file and env vars,
// then uses this configuration to setup clients for Azure AI Language
// service APIs and GitHub APIs before registering HTTP handlers for
// GitHub webhook events.
func (m *Manager) runServer() error {
	m.logRoutes()
	// run the HTTP server
	return m.server.ListenAndServe()
}

// setupEventDispatcher() function returns an http.Handler that can be used
// as the event dispatcher for GitHub webhook events sent to the HTTP server;
// returns a non-nil error if unable to setup the event dispatcher handler.
func setupEventDispatcher(config *cfg.Config) (http.Handler, error) {
	// create a common *gh.ClientManager, which can be used for interacting
	// with the GitHub API via the go-github libarary
	ghcm, err := gh.NewClientManager(config)
	if err != nil {
		return nil, err
	}

	// create a common *az.EntityDetectionAI, which can be used for detecting
	// "entities" of interest within text documents submitted to the the API
	// for the Azure AI Language service
	ai, ai_err := az.NewEntityDetectionAI(config)
	if ai_err != nil {
		return nil, ai_err
	}

	// define the event handlers
	installationHandler := &handlers.InstallationHandler{
		GHCM: ghcm,
	}
	issueCommentHandler := &handlers.IssueCommentHandler{
		AI:   ai,
		GHCM: ghcm,
	}
	pullRequestHandler := &handlers.PullRequestHandler{
		AI:   ai,
		GHCM: ghcm,
	}
	pushHandler := &handlers.PushHandler{
		AI:   ai,
		GHCM: ghcm,
	}

	// register the event handlers with a new/default event dispatcher
	eventDispatcher := githubapp.NewDefaultEventDispatcher(
		*config.GitHub.GetGitHubAppConfig(),
		installationHandler,
		issueCommentHandler,
		pullRequestHandler,
		pushHandler,
	)

	return eventDispatcher, nil
}
