package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/config"
	"github.com/ziflex/lecho/v3"
	"riverqueue.com/riverui"
)

type Services struct {
	ChannelService ChannelService
	VideoService   VideoService
	CommentService CommentService
	ChapterService ChapterService
	ScannerService ScannerService
	MetricsService MetricsService
	BlockedPaths   BlockedPaths
	RiverUIServer  *riverui.Server
}

type Handler struct {
	Server   *echo.Echo
	Config   config.Config
	Services Services
}

func NewHandler(c config.Config, channelService ChannelService, videoService VideoService, commentService CommentService, chapterService ChapterService, scannerService ScannerService, metricsService MetricsService, blockedPaths BlockedPaths, riverUIServer *riverui.Server) *Handler {

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
			MetricsService: metricsService,
			BlockedPaths:   blockedPaths,
			RiverUIServer:  riverUIServer,
		},
	}

	h.mapRoutes(c.VIDEOS_DIR)

	return h
}

func (h *Handler) mapRoutes(videosDir string) {

	// Setup Prometheus metrics route
	h.Server.GET("/metrics", func(c echo.Context) error {
		r, err := h.GatherMetrics()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Serve videos directory
	h.Server.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: videosDir,
	}))

	// RiverUI
	h.Server.Any("/riverui/", echo.WrapHandler(h.Services.RiverUIServer))
	h.Server.Any("/riverui/*", echo.WrapHandler(h.Services.RiverUIServer))

	// enable gzip
	h.Server.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "metrics") // skip gzip for metrics endpoint
		},
	}))

	// serve public directory for assets
	h.Server.Static("/public", "public")
	h.Server.Static(videosDir, videosDir)

	h.Server.GET("/", h.HandleLandingIndex)

	h.Server.GET("/channels/:name", h.HandleChannelPage)

	h.Server.GET("/videos/:video_id", h.HandelVideoPage)
	h.Server.GET("/videos/:video_id/comments", h.HandleVideoCommentsPage)
	h.Server.GET("/videos/:video_id/chapters", h.GetChaptersForVideo)
	h.Server.GET("/videos/:video_id/comments/:comment_id/replies", h.HandleVideoCommentReplies)
	h.Server.GET("/videos/search", h.HandleVideoSearchPage)

	h.Server.GET("/admin/blocked-paths", h.HandleBlockedPathsPage)

	v1 := h.Server.Group("/api/v1")

	// Scanner
	scannerGroup := v1.Group("/scanner")
	scannerGroup.POST("/start", h.StartScanner)

	// blocked paths
	blockedPathsGroup := v1.Group("/blocked-paths")
	blockedPathsGroup.DELETE("/:id", h.DeleteBlockedPathById)
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
