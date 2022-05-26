package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Services struct {
	VideoService   VideoService
	ChannelService ChannelService
	CommentService CommentService
	ScannerService ScannerService
	MetricsService MetricsService
}

type Handler struct {
	Server  *echo.Echo
	Service Services
}

func NewHandler(videoService VideoService, channelService ChannelService, commentService CommentService, scannerService ScannerService, metricsService MetricsService) *Handler {
	h := &Handler{
		Server: echo.New(),
		Service: Services{
			VideoService:   videoService,
			ChannelService: channelService,
			CommentService: commentService,
			ScannerService: scannerService,
			MetricsService: metricsService,
		},
	}

	h.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	h.Server.Debug = true
	h.Server.Logger.SetLevel(log.DEBUG)

	h.mapRoutes()

	// Middleware
	h.Server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	return h
}

func (h *Handler) mapRoutes() {
	h.Server.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World")
	})

	h.Server.GET("/metrics", func(c echo.Context) error {
		r := h.GatherMetrics()

		handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
		handler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	h.Server.GET("/channels", h.GetChannels)
	h.Server.GET("/channels/:id", h.GetChannel)
	h.Server.GET("/channels/name/:name", h.GetChannelByName)
	h.Server.POST("/channels", h.CreateChannel)

	h.Server.GET("/videos/:vid_id", h.GetVideo)
	h.Server.GET("/videos/channel/:channel_id", h.GetChannelVideos)
	h.Server.GET("/videos/random", h.GetRandomVideos)
	h.Server.GET("/videos/search", h.SearchVideos)
	h.Server.POST("/videos", h.CreateVideo)

	h.Server.POST("/comments", h.CreateComment)
	h.Server.GET("/comments/:vid_id", h.GetVideoComments)

	h.Server.POST("/scanner", h.StartScanner)

}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.Start(":3001"); err != nil && err != http.ErrServerClosed {
			h.Server.Logger.Fatal("Shutting down the server")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Server.Shutdown(ctx); err != nil {
		h.Server.Logger.Fatal(err)
	}

	return nil
}
