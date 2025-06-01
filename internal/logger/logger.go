package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func Initialize() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	debugLog := os.Getenv("DEBUG")
	if debugLog == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// stack traces
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// pretty print in dev mode
	devMode := os.Getenv("DEV")
	if devMode == "true" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
