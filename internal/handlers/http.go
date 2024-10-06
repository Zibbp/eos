package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/config"
	"github.com/ziflex/lecho/v3"
)

type Services struct {
	ChannelService ChannelService
	VideoService   VideoService
	CommentService CommentService
	ChapterService ChapterService
	ScannerService ScannerService
}

type Handler struct {
	Server   *echo.Echo
	Config   config.Config
	Services Services
}

func NewHandler(c config.Config, channelService ChannelService, videoService VideoService, commentService CommentService, chapterService ChapterService, scannerService ScannerService) *Handler {

	e := echo.New()

	e.HideBanner = true

	e.Logger = lecho.From(log.Logger)

	h := &Handler{
		Server: e,
		Config: c,
		Services: Services{
			VideoService:   videoService,
			ScannerService: scannerService,
			ChannelService: channelService,
			ChapterService: chapterService,
			CommentService: commentService,
		},
	}

	h.mapRoutes(c.VIDEOS_DIR)

	return h
}

func (h *Handler) mapRoutes(videosDir string) {

	// enable gzip
	h.Server.Use(middleware.Gzip())

	// serve public directory for assets
	h.Server.Static("/public", "public")
	h.Server.Static(videosDir, videosDir)

	h.Server.GET("/", h.HandleLandingIndex)

	h.Server.GET("/channels", h.HandleChannelsPage)
	h.Server.GET("/channels/:name", h.HandleChannelPage)

	h.Server.GET("/videos/:video_id", h.HandelVideoPage)
	h.Server.GET("/videos/:video_id/comments", h.HandleVideoCommentsPage)
	h.Server.GET("/videos/:video_id/chapters", h.GetChaptersForVideo)
	h.Server.GET("/videos/:video_id/comments/:comment_id/replies", h.HandleVideoCommentReplies)

	v1 := h.Server.Group("/api/v1")

	// Scanner
	scannerGroup := v1.Group("/scanner")
	scannerGroup.POST("/start", h.StartScanner)
}

func (h *Handler) Serve() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := h.Server.Start(":3000"); err != nil && err != http.ErrServerClosed {
			h.Server.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Server.Shutdown(ctx); err != nil {
		h.Server.Logger.Fatal(err)
	}

	return nil
}
