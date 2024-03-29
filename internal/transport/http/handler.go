package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	asynqMetrics "github.com/hibiken/asynq/x/metrics"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/eos/internal/metrics"
	"github.com/zibbp/eos/internal/redis"
	"github.com/zibbp/eos/internal/utils"
)

type Services struct {
	VideoService   VideoService
	ChannelService ChannelService
	CommentService CommentService
}

type Handler struct {
	Server  *echo.Echo
	Service Services
}

func NewHandler(videoService VideoService, channelService ChannelService, commentService CommentService) *Handler {
	log.Debug().Msg("initializing http handler")
	h := &Handler{
		Server: echo.New(),
		Service: Services{
			VideoService:   videoService,
			ChannelService: channelService,
			CommentService: commentService,
		},
	}

	h.Server.HideBanner = true

	// Middleware
	h.Server.Validator = &utils.CustomValidator{Validator: validator.New()}

	// CORS
	h.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	h.mapRoutes()

	return h
}

func (h *Handler) mapRoutes() {
	log.Debug().Msg("mapping routes")

	h.Server.GET("/", func(c echo.Context) error {
		return c.String(200, "EOS API")
	})

	// Metrics
	inspector := redis.GetAsynqInspector()
	tmp := asynqMetrics.NewQueueMetricsCollector(inspector.Inspector)
	eosMetrics := metrics.NewEosMetricsCollector()
	prometheus.MustRegister(tmp)
	prometheus.MustRegister(eosMetrics)

	h.Server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	v1 := h.Server.Group("/api/v1")
	groupV1Routes(v1, h)
}

func groupV1Routes(e *echo.Group, h *Handler) {
	log.Debug().Msg("mapping v1 routes")

	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// Channel group
	channelGroup := e.Group("/channels")
	channelGroup.GET("", h.GetChannels)
	channelGroup.GET("/:id", h.GetChannel)
	channelGroup.POST("", h.CreateChannel)
	channelGroup.PUT("/:id", h.UpdateChannel)
	channelGroup.DELETE("/:id", h.DeleteChannel)
	channelGroup.GET("/name/:name", h.GetChannelByName)

	// Video group
	videoGroup := e.Group("/videos")
	videoGroup.GET("", h.GetVideos)
	videoGroup.GET("/:id", h.GetVideo)
	videoGroup.POST("", h.CreateVideo)
	videoGroup.PUT("/:id", h.UpdateVideo)
	videoGroup.DELETE("/:id", h.DeleteVideo)
	videoGroup.GET("/channel/:id", h.GetVideosByChannelID)
	videoGroup.GET("/random", h.GetRandomVideos)
	videoGroup.GET("/search", h.SearchVideos)
	videoGroup.POST("/generate_thumbnails_vtt", h.GenerateThumbnailsVTT)

	// Comment group
	commentGroup := e.Group("/comments")
	commentGroup.GET("", h.GetComments)

	// Tasks group
	tasksGroup := e.Group("/tasks")
	tasksGroup.POST("/video/start_scanner", h.StartVideoScannerTask)
	tasksGroup.POST("/video/generate_thumbnails", h.StartVideoGenerateAllThumbnailsTask)
	tasksGroup.POST("/video/generate_thumbnails/:id", h.StartVideoGenerateThumbnailsTask)
	tasksGroup.POST("/video/download_thumbnails", h.StartVideoDownloadThumbnailsTask)

	// Playback group
	playbackGroup := e.Group("/playback")
	playbackGroup.POST("/:video_id", UpdateProgress)
	playbackGroup.GET("/:video_id", GetProgress)
	playbackGroup.GET("", GetAllProgress)
	playbackGroup.DELETE("/:video_id", DeleteProgress)

}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.Start(":4000"); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
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
		log.Fatal().Err(err).Msg("failed to shutdown server")
	}

	return nil
}
