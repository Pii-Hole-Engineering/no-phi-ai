package cfg

import (
	"flag"
	"os"

	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

// AppConfig struct contains the configuration items used by the running app.
type AppConfig struct {
	// Log contains the configuration items used by the logger.
	Log AppLogConfig `yaml:"log" json:"log"`
	// Mode can be:
	//   - AppModeCLI to use CommandConfig to run a command once, like a CLI tool;
	//   - AppModeServer to use ServerConfig to run a GitHub App that listens for
	//     webhook events;
	//
	// Mode default is defined in DefaultAppMode const.
	Mode string `yaml:"mode" json:"mode"`
	// Name is the name of the app.
	//
	// Name default is defined as DefaultAppName const.
	Name string `yaml:"name" json:"name"`
	// UserAgent is the user agent string used by the app when making HTTP requests.
	//
	// UserAgent default is defined in DefaultAppUserAgent const
	UserAgent string `yaml:"user_agent" json:"user_agent"`
}

// AppLogConfig struct contains the configuration used to initialize the logger.
type AppLogConfig struct {
	// ConsoleEnable controls whether or not to log to standard output
	ConsoleEnable bool `yaml:"console_enable" json:"console_enable"`
	// ConsolePretty controls whether console output is pretty printed (true)
	// or printed as structured JSON (false)
	ConsolePretty bool `yaml:"console_pretty" json:"console_pretty"`
	// "trace", "debug", "info", "warn", or "error"
	Level string `yaml:"level" json:"level"`
	// log to standard output -> "" or "stdout"
	// log to a file -> "../relative/file/path" or "/absolute/file/path"
	File string `yaml:"file" json:"file"`
}

// AzureAIConfig struct contains the configuration items used to create a client
// for interacting with the APIs of the Azure AI Language service.
type AzureAIConfig struct {
	// AuthKey should be set to the value of a "key" associated with the
	// AI Language service resource in Azure.
	AuthKey string `yaml:"auth_key" json:"auth_key"`
	// ConnfidenceThreshold is the minimum confidence score required for
	// a detection result to be considered valid. This must be a value
	// between 0 and 1.
	ConfidenceThreshold float64 `yaml:"confidence_threshold" json:"confidence_threshold"`
	// DryRun prevents the actual sending of requests to the AI Language
	// service API when set to true. Default is false.
	DryRun bool `yaml:"dry_run" json:"dry_run"`
	// Service should be set to the URL of the AI Language service deployment.
	Service string `yaml:"service" json:"service"`
	// ShowStats controls whether the "showStats=true" query parameter will
	// be added to "analyze" requests sent to the AI Language service.
	ShowStats bool `yaml:"show_stats" json:"show_stats"`
}

// CommandConfig struct contains the configuration used to run a command.
// Only used when AppConfig.Mode == AppModeCLI.
type CommandConfig struct {
	// available commands include:
	//   - "help" to print help text
	//   - "list-org-repos" to list repos in an org (for testing) // TODO
	//   - "scan-org" to scan an org for PHI // TODO
	//   - "scan-repos" to scan a repo for PHI // TODO
	//   - "version" to print the app version
	Run string `yaml:"run" json:"run"`
}

// GitAuthConfig struct contains the configuration used to setup
// authentication for GitHub API clients, including cloning repos via
// the git protocol.
//
// User should supply one of the following:
//   - SSHKeyPath
//   - Token
type GitAuthConfig struct {
	// SSHKeyPath is the path to the SSH private key used to authenticate
	// to GitHub via the git protocol.
	SSHKeyPath string `yaml:"ssh_key_path" json:"ssh_key_path"`
	// set to the value of your Personal Access Token in order to allow
	// the app to authenticate to GitHub via HTTPS with OAuth2.
	//
	// TODO : implement this
	Token string `yaml:"token" json:"token"`
}

// GitConfig struct contains the configuration used clone from and
// push (commits) to repos using the git protocol.
type GitConfig struct {
	// Auth config for git protocol and GitHub API clients
	Auth GitAuthConfig `yaml:"auth" json:"auth"`
	// control the behavior of the CLI by specifying the organization
	// and/or repositories to scan
	Scan GitScanConfig `yaml:"scan" json:"scan"`
	// WorkDir is the base directory for working with git repos.
	//
	// Any git repositories processed by the app will be cloned into
	// a subdirectory of WorkDir.
	WorkDir string `yaml:"work_dir" json:"workDir"`
}

// GitHubConfig struct contains the configuration used to create clients for
// interacting with GitHub APIs (outbound) and webhook events (inbound).
type GitHubConfig struct {
	// App configuration is required for conversion to githubapp.Config struct,
	// which is required when running the app in "server" mode. Running a secure
	// installation of a GitHub app requires prior setup of the values used for
	// IntegrationID, WebhookSecret, and PrivateKey.
	App struct {
		IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
		PrivateKey    string `yaml:"private_key" json:"privateKey"`
		WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
	} `yaml:"app" json:"app"`
	// OAuth (and other) configurations are not currently required / used in
	//  the app, but are required for conversion to a githubapp.Config struct.
	//
	// OAuth is required for conversion to a githubapp.Config struct.
	OAuth struct {
		ClientID     string `yaml:"client_id" json:"clientId"`
		ClientSecret string `yaml:"client_secret" json:"clientSecret"`
	} `yaml:"oauth" json:"oauth"`
	// WebURL is required for conversion to a githubapp.Config struct.
	WebURL string `yaml:"web_url" json:"webUrl"`
	// V3APIURL is required for conversion to a githubapp.Config struct.
	V3APIURL string `yaml:"v3_api_url" json:"v3ApiUrl"`
	// V4APIURL is required for conversion to a githubapp.Config struct.
	V4APIURL string `yaml:"v4_api_url" json:"v4ApiUrl"`
}

// GetGitHubAppConfig() method converts the GitHub portion of the Config to
// a githubapp.Config struct that can be used with the githubapp package.
func (c *GitHubConfig) GetGitHubAppConfig() *githubapp.Config {
	out := &githubapp.Config{}

	out.App.IntegrationID = c.App.IntegrationID
	out.App.WebhookSecret = c.App.WebhookSecret
	out.App.PrivateKey = c.App.PrivateKey
	out.OAuth.ClientID = c.OAuth.ClientID
	out.OAuth.ClientSecret = c.OAuth.ClientSecret
	out.WebURL = c.WebURL
	out.V3APIURL = c.V3APIURL
	out.V4APIURL = c.V4APIURL

	return out
}

// GitScanConfig struct contains the configuration used to setup a PHI scan
// for some organization and/or set of repositories.
type GitScanConfig struct {
	// Extensions is a list of file extensions to include in the scan, where
	// each entry is a string in the format ".<ext>". If this list empty,
	// then the DefaultScanFileExtensions list will be used.
	Extensions []string `yaml:"extensions" json:"extensions"`

	// IgnoreExtensions is a list of file extensions to exclude/ignore from
	// the scan, where each entry is a string in the format ".<ext>".
	IgnoreExtensions []string `yaml:"ignore_extensions" json:"ignore_extensions"`

	// IgnoreRepositories is a list of GitHub repositories to exclude/ignore
	// from the scan, where each entry is a string in the format "<org>/<repo>"
	// or "<user>/<repo>". Values in this list take precedence over values in
	// the Repositories list.
	IgnoreRepositories []string `yaml:"ignore_repositories" json:"ignore_repositories"`

	// Limits config
	Limits GitScanLimitsConfig `yaml:"limits" json:"limits"`

	// Organization is the URL of the GitHub organization to scan, where the
	// app will query the GitHub API for a list of repositories to scan.
	Organization string `yaml:"organization" json:"organization"`

	// Repositories is a list of GitHub repositories to scan, where each entry
	// is a string in the format "<org>/<repo>" or "<user>/<repo>".
	//
	// Repositories can be used in parallel with Organization, where the app
	// will scan all repositories in the Organization, plus any repositories
	// listed in Repositories, minus any duplicates and minus any repositories
	// listed in the IgnoreRepositories list.
	Repositories []string `yaml:"repositories" json:"repositories"`
}

type GitScanLimitsConfig struct {
	MaxRequestChunkSize    int `yaml:"max_request_chunk_size" json:"max_request_chunk_size"`
	MaxRequestsOutstanding int `yaml:"max_requests_outstanding" json:"max_requests_outstanding"`
}

// ServerConfig struct contains the configuration used to start the HTTP server.
// Only used when AppConfig.Mode == "server".
type ServerConfig struct {
	Address   string  `yaml:"address" json:"address"`
	Port      int     `yaml:"port" json:"port"`
	RateLimit float64 `yaml:"rate_limit" json:"rate_limit"`
}

// Config struct is the top-level configuration object for the app.
type Config struct {
	App     AppConfig     `yaml:"app" json:"app"`
	AzureAI AzureAIConfig `yaml:"azure_ai" json:"azure_ai"`
	Command CommandConfig `yaml:"command" json:"command"`
	Git     GitConfig     `yaml:"git" json:"git"`
	GitHub  GitHubConfig  `yaml:"github" json:"github"`
	Server  ServerConfig  `yaml:"server" json:"server"`
}

// NewDefaultConfig() function returns a new Config object with default values
// set for optional fields. The returned config can be used as a basis for some
// customer configuration that actually works, where the default values are not
// enough to provide required information for some app commands / functions.
func NewDefaultConfig() *Config {
	// create an empty config
	config := &Config{}
	// set default values for optional config fields
	config.defaultConfig()

	return config
}

// default() method sets default values for optional config fields.
func (c *Config) defaultConfig() {
	// set defaults for optional c.App config values
	if c.App.Mode == "" {
		c.App.Mode = DefaultAppMode
	}
	if c.App.Name == "" {
		c.App.Name = DefaultAppName
	}
	if c.App.Log.Level == "" {
		c.App.Log.Level = DefaultAppLogLevel
	}
	if c.App.UserAgent == "" {
		c.App.UserAgent = DefaultAppUserAgent
	}
	// golang boolean default is false, but we want to c.AzureAI.ShowStats
	// default to true, so we set it here and force the user to override
	// with env var NOPHI_AZURE_AI_SHOW_STATS=false
	c.AzureAI.ShowStats = DefaultAzureAIShowStats
	if c.Command.Run == "" {
		c.Command.Run = DefaultCommandRun
	}
	if len(c.Git.Scan.Extensions) == 0 {
		c.Git.Scan.Extensions = DefaultScanFileExtensions
	}
	if c.Git.Scan.Limits.MaxRequestChunkSize == 0 {
		c.Git.Scan.Limits.MaxRequestChunkSize = DefaultMaxRequestChunkSize
	}
	if c.Git.Scan.Limits.MaxRequestsOutstanding == 0 {
		c.Git.Scan.Limits.MaxRequestsOutstanding = DefaultMaxRequestsOutstanding
	}
	if c.Git.WorkDir == "" {
		c.Git.WorkDir = DefaultCommandWorkDir
	}
	// set defaults for optional c.GitHub config values
	if c.GitHub.V3APIURL == "" {
		c.GitHub.V3APIURL = DefaultGitHubV3APIURL
	}
	// set defaults for optional c.Server config values
	if c.Server.Address == "" {
		c.Server.Address = DefaultServerAddress
	}
	if c.Server.Port == 0 {
		c.Server.Port = DefaultServerPort
	}
	if c.Server.RateLimit == 0 {
		c.Server.RateLimit = DefaultRateLimit
	}
}

// verifyConfig() method returns an error if a required value has not
// been set in the Config, sets defaults for optional values, and/or
// returns a nil error if all required values are set.
func (c *Config) verifyConfig() (e error) {
	switch c.App.Mode {
	case AppModeCLI:
		e = c.verifyConfigCLI()
		return
	case AppModeServer:
		e = c.verifyConfigServer()
		return
	default:
		e = errors.New("invalid config value: app.mode = " + c.App.Mode)
		return
	}
}

// verifyConfigCLI() method verifies required config values when running the app
// in "cli" mode.
func (c *Config) verifyConfigCLI() (e error) {
	// check the c.AzureAI config values
	if c.AzureAI.Service == "" {
		e = errors.New("missing required config value: azure_ai.service")
		return
	}
	if c.AzureAI.AuthKey == "" {
		e = errors.New("missing required config value: azure_ai.auth_key")
		return
	}
	if c.AzureAI.ConfidenceThreshold == 0 {
		c.AzureAI.ConfidenceThreshold = DefaultConfidenceThreshold
	}

	// check the c.Git.Auth.Token config value
	if c.Git.Auth.SSHKeyPath == "" && c.Git.Auth.Token == "" {
		e = errors.New("missing required config value: either 'github.auth.ssh_key_path' or github.auth.token' must be set")
		return
	}

	return
}

// verifyConfigServer() method verifies required config values when running the app
// in "server" mode.
func (c *Config) verifyConfigServer() (e error) {
	// check the c.AzureAI config values
	if c.AzureAI.Service == "" {
		e = errors.New("missing required config value: azure_ai.service")
		return
	}
	if c.AzureAI.AuthKey == "" {
		e = errors.New("missing required config value: azure_ai.auth_key")
		return
	}
	if c.AzureAI.ConfidenceThreshold == 0 {
		c.AzureAI.ConfidenceThreshold = DefaultConfidenceThreshold
	}

	// check the c.GitHub config values
	if c.GitHub.App.IntegrationID == 0 {
		e = errors.New("missing required config value: github.app.integration_id")
		return
	}
	if c.GitHub.App.PrivateKey == "" {
		e = errors.New("missing required config value: github.app.private_key")
		return
	}
	if c.GitHub.App.WebhookSecret == "" {
		e = errors.New("missing required config value: github.app.webhook_secret")
		return
	}

	return
}

// readConfigFile() function reads config data from the input file path,
// or returns an empty Config if path is empty in order to allow for the
// entire config to be provided via environment variables.
func readConfigFile(path string) (*Config, error) {
	var c Config

	if path == "" {
		return &c, nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading server config file: %s", path)
	}

	if err := yaml.UnmarshalStrict(bytes, &c); err != nil {
		return nil, errors.Wrap(err, "failed parsing server config file")
	}

	return &c, nil
}

// ParseConfig() function parses the config file and environment variables.
func ParseConfig() (*Config, *zerolog.Logger, error) {
	// define flags
	configPath := flag.String("config", "", "local relative path to the config file")

	// parse flags
	flag.Parse()

	// get config path from environment variable if not set via flag
	if *configPath == "" {
		if envConfigPath, found := os.LookupEnv(NOPHI_CONFIG_PATH); found {
			*configPath = envConfigPath
		}
	}

	// get config from file
	c, err := readConfigFile(*configPath)
	if err != nil {
		return c, nil, err
	}

	// fill in default values for optional config fields
	c.defaultConfig()

	// override config values with environment variables
	if err := c.envOverride(); err != nil {
		return c, nil, err
	}

	// verify required config values are set (i.e. not empty)
	if err := c.verifyConfig(); err != nil {
		return c, nil, err
	}

	return c, c.setupLogger(), nil
}
