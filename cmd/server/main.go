package main

import (
	"github.com/labstack/gommon/log"

	"github.com/zibbp/eos/internal/channel"
	"github.com/zibbp/eos/internal/comment"
	"github.com/zibbp/eos/internal/database"
	"github.com/zibbp/eos/internal/scanner"
	transportHttp "github.com/zibbp/eos/internal/transport/http"
	"github.com/zibbp/eos/internal/video"
)

func Run() error {
	// log start server
	log.Info("Starting Server")

	var err error
	store, err := database.NewDatabase()
	if err != nil {
		log.Error("Failed to connect to database")
		return err
	}
	err = store.MigrateDB()
	if err != nil {
		log.Error("Failed to migrate database")
		return err
	}

	// Create services
	videoService := video.NewService(store)
	channelService := channel.NewService(store)
	commentService := comment.NewService(store)
	scannerService := scanner.NewService(store)

	httpHandler := transportHttp.NewHandler(videoService, channelService, commentService, scannerService)

	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
