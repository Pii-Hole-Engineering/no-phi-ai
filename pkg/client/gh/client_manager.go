package gh

import (
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
)

// ClientManager implements the methods of the githubapp.ClientCreator interface
// and adds additional methods for implementing the business logic of the app.
type ClientManager struct {
	githubapp.ClientCreator
}

// NewClientManager() function initializes a new ClientManager object
// using the provided config.
func NewClientManager(config *cfg.Config) (*ClientManager, error) {
	// TODO : do something more with the metrics registry
	metricsRegistry := metrics.DefaultRegistry

	// create a common githubapp.ClientCreator, which can be used to get an
	// installation client for interacting with GitHub APIs
	cc, err := githubapp.NewDefaultCachingClientCreator(
		*config.GitHub.GetGitHubAppConfig(),
		githubapp.WithClientUserAgent(config.App.UserAgent),
		githubapp.WithClientTimeout(cfg.DefaultClientTimeout),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(metricsRegistry),
			githubapp.ClientLogging(zerolog.GlobalLevel()),
		),
	)
	if err != nil {
		return nil, err
	}

	return &ClientManager{cc}, nil
}
