package cfg

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// setupLogger() method uses the Config settings to setup logging
// for the app.
func (c *Config) setupLogger() *zerolog.Logger {
	// set the global log level
	switch strings.ToLower(c.App.Log.Level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var level_writer io.Writer = zerolog.LevelWriterAdapter{Writer: os.Stdout}
	if c.App.Log.ConsolePretty {
		level_writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	// create the zerolog.Logger instance that will be re-used throughout the app
	logger := zerolog.New(level_writer).With().Caller().Timestamp().Logger()

	if c.App.Log.File != "" {
		log_file, err := os.OpenFile(c.App.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to open log file=%s", c.App.Log.File)
		} else if c.App.Log.ConsoleEnable {
			// log to both console and file
			logger = logger.Output(zerolog.MultiLevelWriter(level_writer, log_file))
		} else {
			// just log to the file
			logger = logger.Output(log_file)
		}
		logger.Info().Msgf("logger setup to log to file=%s", c.App.Log.File)
	}

	zerolog.DefaultContextLogger = &logger
	logger.Trace().Msgf("%s app logger setup complete : log_level=%s", c.App.Name, zerolog.GlobalLevel())

	return &logger
}
