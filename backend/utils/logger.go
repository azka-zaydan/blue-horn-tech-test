package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the Zerolog logger
func InitLogger() {
	// Set global logger level to Info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure console writer for human-readable output
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}
