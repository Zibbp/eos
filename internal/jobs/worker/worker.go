package jobs_worker

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/rs/zerolog/log"
	db "github.com/zibbp/eos/internal/db/sqlc"
	"github.com/zibbp/eos/internal/jobs"
	"github.com/zibbp/eos/internal/parser"
)

type contextKey string

const storeKey contextKey = "store"

type RiverWorkerInput struct {
	DB_URL             string
	DB                 db.Store
	ParserYoutube      parser.Parser
	VideoImportWorkers int
}

type RiverWorkerClient struct {
	Ctx            context.Context
	PgxPool        *pgxpool.Pool
	RiverPgxDriver *riverpgxv5.Driver
	Client         *river.Client[pgx.Tx]
}

type CustomErrorHandler struct{}

func (*CustomErrorHandler) HandleError(ctx context.Context, job *rivertype.JobRow, err error) *river.ErrorHandlerResult {
	log.Error().Str("job_id", fmt.Sprintf("%d", job.ID)).Str("attempt", fmt.Sprintf("%d", job.Attempt)).Str("attempted_by", job.AttemptedBy[job.Attempt-1]).Err(err).Msg("job errored")
	return nil
}

func (*CustomErrorHandler) HandlePanic(ctx context.Context, job *rivertype.JobRow, panicVal any) *river.ErrorHandlerResult {
	log.Error().Str("job_id", fmt.Sprintf("%d", job.ID)).Str("attempt", fmt.Sprintf("%d", job.Attempt)).Str("attempted_by", job.AttemptedBy[job.Attempt-1]).Str("panic_val", fmt.Sprintf("%v", panicVal)).Msg("job errored")
	return nil
}

func NewRiverWorker(input RiverWorkerInput, store db.Store) (*RiverWorkerClient, error) {
	rc := &RiverWorkerClient{}

	workers := river.NewWorkers()
	if err := river.AddWorkerSafely(workers, &jobs.VideoImportWorker{}); err != nil {
		return rc, err
	}
	if err := river.AddWorkerSafely(workers, &jobs.VideoImportCommentsWorker{}); err != nil {
		return rc, err
	}
	if err := river.AddWorkerSafely(workers, &jobs.VideoScanWorker{}); err != nil {
		return rc, err
	}

	rc.Ctx = context.Background()

	// create postgres pool connection
	pool, err := pgxpool.New(rc.Ctx, input.DB_URL)
	if err != nil {
		return rc, fmt.Errorf("error connecting to postgres: %v", err)
	}
	rc.PgxPool = pool

	// create river pgx driver
	rc.RiverPgxDriver = riverpgxv5.New(rc.PgxPool)

	// create river client
	riverClient, err := river.NewClient(rc.RiverPgxDriver, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault:    {MaxWorkers: 100}, // non-resource intensive tasks or time sensitive tasks
			jobs.QueueScanVideo:   {MaxWorkers: 1},
			jobs.QueueImportVideo: {MaxWorkers: input.VideoImportWorkers},
		},
		Workers:      workers,
		ErrorHandler: &jobs.CustomErrorHandler{},
	})
	if err != nil {
		return rc, fmt.Errorf("error creating river client: %v", err)
	}
	rc.Client = riverClient

	log.Info().Str("default_workers", "100").Str("video_import_workers", fmt.Sprintf("%d", input.VideoImportWorkers)).Msg("created river worker")

	// put store in the rc.Ctx
	rc.Ctx = context.WithValue(rc.Ctx, jobs.StoreKey, store)

	// put parser in the rc.Ctx
	rc.Ctx = context.WithValue(rc.Ctx, jobs.ParserKey, input.ParserYoutube)

	return rc, nil
}

func (rc *RiverWorkerClient) Start() error {
	log.Info().Str("name", rc.Client.ID()).Msg("starting wortker")
	if err := rc.Client.Start(rc.Ctx); err != nil {
		return err
	}
	return nil
}

func (rc *RiverWorkerClient) Stop() error {
	if err := rc.Client.Stop(rc.Ctx); err != nil {
		return err
	}
	return nil
}
