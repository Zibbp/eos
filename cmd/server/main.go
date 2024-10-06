package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/channel"
	"github.com/zibbp/eos/internal/chapter"
	"github.com/zibbp/eos/internal/comment"
	"github.com/zibbp/eos/internal/config"
	goose "github.com/zibbp/eos/internal/db"
	db "github.com/zibbp/eos/internal/db/sqlc"
	"github.com/zibbp/eos/internal/handlers"
	jobs_client "github.com/zibbp/eos/internal/jobs/client"
	"github.com/zibbp/eos/internal/logger"
	"github.com/zibbp/eos/internal/scanner"
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

	err = goose.RunGooseMigrations(dbString)
	if err != nil {
		log.Panic().Err(err).Msg("Error running goose migrations")
	}

	store := db.NewStore(connPool)

	// initialize river
	riverClient, err := jobs_client.NewRiverClient(jobs_client.RiverClientInput{
		DB_URL: c.DB_URL,
	})
	if err != nil {
		log.Panic().Err(err).Msg("Error creating river client")
	}

	err = riverClient.RunMigrations()
	if err != nil {
		log.Panic().Err(err).Msg("Error running river migrations")
	}

	log.Info().Msg("starting server")

	// set services
	channelService := channel.NewService(store)
	videoService := video.NewService(store)
	commentService := comment.NewService(store)
	chapterService := chapter.NewService(store)
	scannerService := scanner.NewScannerService(riverClient.Client, store, channelService, videoService, c.VIDEOS_DIR)

	handler := handlers.NewHandler(c, channelService, videoService, commentService, chapterService, scannerService)

	handler.Serve()

	// server := server.NewServer()

	// err := server.ListenAndServe()
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot start server: %s", err))
	// }
}
