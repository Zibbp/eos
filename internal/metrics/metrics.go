package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Store interface {
	GetVideoCount() (int64, error)
	GetChannelCount() (int64, error)
	GetChannelVideoCount() ([]ChannelVideoCount, error)
}

type ChannelVideoCount struct {
	Channel string
	Count   int
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

// Metrics
var (
	totalVideos = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_videos",
		Help: "Total number of videos",
	})
	totalChannels = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_channels",
		Help: "Total number of channels",
	})
	channelVideoCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "channel_video_count",
		Help: "Total number of videos per channel",
	}, []string{"channel"})
)

func (s *Service) GatherMetrics() *prometheus.Registry {
	// Gather metric data
	videoCount, err := s.Store.GetVideoCount()
	if err != nil {
		totalVideos.Set(0)
	}
	channelCount, err := s.Store.GetChannelCount()
	if err != nil {
		totalChannels.Set(0)
	}
	channelVideoCounts, err := s.Store.GetChannelVideoCount()
	if err != nil {
		channelVideoCount.Reset()
	}
	// Set metric data
	totalVideos.Set(float64(videoCount))
	totalChannels.Set(float64(channelCount))
	for _, cVC := range channelVideoCounts {
		channelVideoCount.WithLabelValues(cVC.Channel).Set(float64(cVC.Count))
	}
	// Remove default metrics - not needed in this tiny application
	r := prometheus.NewRegistry()
	r.MustRegister(totalVideos, totalChannels, channelVideoCount)
	return r
}
