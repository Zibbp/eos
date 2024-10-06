package goose

import (
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// CustomLogger is a custom logger that uses zerolog to log messages.
type CustomLogger struct{}

func (l *CustomLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func (l *CustomLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

// RunGooseMigrations runs goose migrations.
func RunGooseMigrations(dbString string) error {
	goose.SetLogger(&CustomLogger{})
	goose.SetBaseFS(embedMigrations)
	goose.SetDialect("postgres")

	gooseDB, err := goose.OpenDBWithDriver("postgres", dbString)
	if err != nil {
		return err
	}
	defer gooseDB.Close()

	if err := goose.Up(gooseDB, "migrations"); err != nil {
		return err
	}

	if err := goose.Version(gooseDB, "migrations"); err != nil {
		return err
	}

	gooseDB.Close()

	return nil
}
