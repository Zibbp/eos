package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/channel"
	"github.com/zibbp/eos/internal/chapter"
	"github.com/zibbp/eos/internal/comment"
	"github.com/zibbp/eos/internal/config"
	db "github.com/zibbp/eos/internal/db/sqlc"
	jobs_worker "github.com/zibbp/eos/internal/jobs/worker"
	"github.com/zibbp/eos/internal/logger"
	"github.com/zibbp/eos/internal/parser"
	"github.com/zibbp/eos/internal/video"
)

func main() {
	ctx := context.Background()

	// setup logger
	logger.Initialize()

	// initialize config
	c := config.GetConfig()

	dbString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", c.DB_USER, c.DB_PASS, c.DB_HOST, c.DB_PORT, c.DB_NAME)

	connPool, err := pgxpool.New(ctx, dbString)
	if err != nil {
		log.Panic().Err(err).Msg("Error connecting to postgres")
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	channelService := channel.NewService(store)
	videoService := video.NewService(store)
	commentService := comment.NewService(store)
	chapterService := chapter.NewService(store)

	var parserYoutube parser.Parser
	// initialize parser
	parserYoutube = &parser.YoutubeParser{
		Store:          store,
		ChannelService: channelService,
		VideoService:   videoService,
		CommentService: commentService,
		ChapterService: chapterService,
	}

	// initialize river
	riverClient, err := jobs_worker.NewRiverWorker(jobs_worker.RiverWorkerInput{
		DB_URL:             c.DB_URL,
		DB:                 store,
		ParserYoutube:      parserYoutube,
		VideoImportWorkers: c.MaxVideoImportWorkers,
	}, store)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating river worker")
	}

	// Start your worker in a goroutine
	go func() {
		if err := riverClient.Start(); err != nil {
			log.Panic().Err(err).Msg("Error running river worker")
		}
	}()

	// Set up channel to listen for OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-sigs

	// Gracefully stop the worker
	if err := riverClient.Stop(); err != nil {
		log.Panic().Err(err).Msg("Error stopping river worker")
	}

	log.Info().Msg("worker stopped")

}
