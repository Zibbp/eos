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
	if debugLog == "TRUE" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// stack traces
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// pretty print in dev mode
	devMode := os.Getenv("DEV")
	if devMode == "TRUE" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
