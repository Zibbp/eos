package jobs_client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/rs/zerolog/log"
)

type RiverClientInput struct {
	DB_URL string
}

type RiverClient struct {
	Ctx            context.Context
	PgxPool        *pgxpool.Pool
	RiverPgxDriver *riverpgxv5.Driver
	Client         *river.Client[pgx.Tx]
}

// Create a river client and the database pool
func NewRiverClient(input RiverClientInput) (*RiverClient, error) {
	rc := &RiverClient{}
	rc.Ctx = context.Background()

	// create postgres pool connection
	pool, err := pgxpool.New(rc.Ctx, input.DB_URL)
	if err != nil {
		return rc, err
	}
	rc.PgxPool = pool

	// create river pgx driver
	rc.RiverPgxDriver = riverpgxv5.New(rc.PgxPool)

	// create river client
	riverClient, err := river.NewClient(rc.RiverPgxDriver, &river.Config{})
	if err != nil {
		return rc, fmt.Errorf("error creating river client: %v", err)
	}
	rc.Client = riverClient

	return rc, nil
}

func (rc *RiverClient) Stop() error {
	if err := rc.Client.Stop(rc.Ctx); err != nil {
		return err
	}
	return nil
}

// Run river database migrations
func (rc *RiverClient) RunMigrations() error {
	migrator, err := rivermigrate.New(rc.RiverPgxDriver, nil)
	if err != nil {
		return fmt.Errorf("error creating river migrator: %v", err)
	}

	_, err = migrator.Migrate(rc.Ctx, rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{})
	if err != nil {
		return fmt.Errorf("error running river migrations: %v", err)
	}

	log.Info().Msg("successfully applied river migrations")

	return nil
}
