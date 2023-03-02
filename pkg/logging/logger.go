package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the logger
func InitLogger() {
	// Set up the logger
	zerolog.TimeFieldFormat = time.RFC3339
	// Set global log level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Use ConsoleWriter with color output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}
	// Set output to stdout
	log.Logger = log.Output(consoleWriter).With().Caller().Logger()
}

// Info logs an info-level message
func Info(message string) {
	log.Info().Msg(message)
}

// Error logs an error-level message
func Error(message string, err error) {
	log.Error().Err(err).Msg(message)
}

// Log a message with the Fatal level
func Fatal(message string, err error) {
	log.Fatal().Err(err).Msg(message)
}
