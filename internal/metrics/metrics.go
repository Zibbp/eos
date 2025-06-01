package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	db "github.com/zibbp/eos/internal/db/sqlc"
	jobs_client "github.com/zibbp/eos/internal/jobs/client"
)

type Service struct {
	Store       db.Store
	riverClient *jobs_client.RiverClient
	metrics     *Metrics
	Registry    *prometheus.Registry
}

type Metrics struct {
	totalVideos             prometheus.Gauge
	totalChannels           prometheus.Gauge
	totalBlockedPaths       prometheus.Gauge
	channelVideoCount       *prometheus.GaugeVec
	riverTotalPendingJobs   *prometheus.GaugeVec
	riverTotalScheduledJobs *prometheus.GaugeVec
	riverTotalAvailableJobs *prometheus.GaugeVec
	riverTotalRunningJobs   *prometheus.GaugeVec
	riverTotalRetryableJobs *prometheus.GaugeVec
	riverTotalCancelledJobs *prometheus.GaugeVec
	riverTotalDiscardedJobs *prometheus.GaugeVec
	riverTotalCompletedJobs *prometheus.GaugeVec
}

func NewService(store db.Store, riverClient *jobs_client.RiverClient) *Service {
	registry := prometheus.NewRegistry()
	metrics := &Metrics{
		totalVideos: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_videos",
			Help: "Total number of videos",
		}),
		totalChannels: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_channels",
			Help: "Total number of channels",
		}),
		totalBlockedPaths: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "total_blocked_paths",
			Help: "Total number of blocked paths",
		}),
		channelVideoCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "channel_video_count",
			Help: "Number of videos per channel",
		}, []string{"channel"}),
		riverTotalPendingJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_pending_jobs",
			Help: "Total number of pending jobs",
		}, []string{"kind"}),
		riverTotalScheduledJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_scheduled_jobs",
			Help: "Total number of scheduled jobs",
		}, []string{"kind"}),
		riverTotalAvailableJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_available_jobs",
			Help: "Total number of available jobs",
		}, []string{"kind"}),
		riverTotalRunningJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_running_jobs",
			Help: "Total number of running jobs",
		}, []string{"kind"}),
		riverTotalRetryableJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_retryable_jobs",
			Help: "Total number of retryable jobs",
		}, []string{"kind"}),
		riverTotalCancelledJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_cancelled_jobs",
			Help: "Total number of cancelled jobs",
		}, []string{"kind"}),
		riverTotalDiscardedJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_discarded_jobs",
			Help: "Total number of discarded jobs",
		}, []string{"kind"}),
		riverTotalCompletedJobs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "river_total_completed_jobs",
			Help: "Total number of completed jobs",
		}, []string{"kind"}),
	}

	registry.MustRegister(
		metrics.totalVideos,
		metrics.totalChannels,
		metrics.totalBlockedPaths,
		metrics.channelVideoCount,
		metrics.riverTotalPendingJobs,
		metrics.riverTotalScheduledJobs,
		metrics.riverTotalAvailableJobs,
		metrics.riverTotalRunningJobs,
		metrics.riverTotalRetryableJobs,
		metrics.riverTotalCancelledJobs,
		metrics.riverTotalDiscardedJobs,
		metrics.riverTotalCompletedJobs,
	)

	return &Service{Store: store, riverClient: riverClient, metrics: metrics, Registry: registry}
}

func (s *Service) gatherRiverJobMetrics() error {
	// Helper function to count jobs by kind
	countJobsByKind := func(jobs []*rivertype.JobRow) map[string]int {
		counts := make(map[string]int)
		for _, job := range jobs {
			counts[job.Kind]++
		}
		return counts
	}

	// Reset all metrics first
	s.metrics.riverTotalPendingJobs.Reset()
	s.metrics.riverTotalScheduledJobs.Reset()
	s.metrics.riverTotalAvailableJobs.Reset()
	s.metrics.riverTotalRunningJobs.Reset()
	s.metrics.riverTotalRetryableJobs.Reset()
	s.metrics.riverTotalCancelledJobs.Reset()
	s.metrics.riverTotalDiscardedJobs.Reset()
	s.metrics.riverTotalCompletedJobs.Reset()

	// Pending Jobs
	pendingJobsParams := river.NewJobListParams().States(rivertype.JobStatePending).First(10000)
	pendingJobs, err := s.riverClient.Client.JobList(context.Background(), pendingJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(pendingJobs.Jobs) {
		s.metrics.riverTotalPendingJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Scheduled Jobs
	scheduledJobsParams := river.NewJobListParams().States(rivertype.JobStateScheduled).First(10000)
	scheduledJobs, err := s.riverClient.Client.JobList(context.Background(), scheduledJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(scheduledJobs.Jobs) {
		s.metrics.riverTotalScheduledJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Available Jobs
	availableJobsParams := river.NewJobListParams().States(rivertype.JobStateAvailable).First(10000)
	availableJobs, err := s.riverClient.Client.JobList(context.Background(), availableJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(availableJobs.Jobs) {
		s.metrics.riverTotalAvailableJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Running Jobs
	runningJobsParams := river.NewJobListParams().States(rivertype.JobStateRunning).First(10000)
	runningJobs, err := s.riverClient.Client.JobList(context.Background(), runningJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(runningJobs.Jobs) {
		s.metrics.riverTotalRunningJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Retryable Jobs
	retryableJobsParams := river.NewJobListParams().States(rivertype.JobStateRetryable).First(10000)
	retryableJobs, err := s.riverClient.Client.JobList(context.Background(), retryableJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(retryableJobs.Jobs) {
		s.metrics.riverTotalRetryableJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Cancelled Jobs
	cancelledJobsParams := river.NewJobListParams().States(rivertype.JobStateCancelled).First(10000)
	cancelledJobs, err := s.riverClient.Client.JobList(context.Background(), cancelledJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(cancelledJobs.Jobs) {
		s.metrics.riverTotalCancelledJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Discarded Jobs
	discardedJobsParams := river.NewJobListParams().States(rivertype.JobStateDiscarded).First(10000)
	discardedJobs, err := s.riverClient.Client.JobList(context.Background(), discardedJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(discardedJobs.Jobs) {
		s.metrics.riverTotalDiscardedJobs.WithLabelValues(kind).Set(float64(count))
	}

	// Completed Jobs
	completedJobsParams := river.NewJobListParams().States(rivertype.JobStateCompleted).First(10000)
	completedJobs, err := s.riverClient.Client.JobList(context.Background(), completedJobsParams)
	if err != nil {
		return err
	}
	for kind, count := range countJobsByKind(completedJobs.Jobs) {
		s.metrics.riverTotalCompletedJobs.WithLabelValues(kind).Set(float64(count))
	}

	return nil
}

func (s *Service) GatherMetrics() (*prometheus.Registry, error) {
	ctx := context.Background()

	// Gather River job metrics
	if err := s.gatherRiverJobMetrics(); err != nil {
		return nil, err
	}

	// Gather total videos
	totalVideos, err := s.Store.GetTotalVideos(ctx)
	if err != nil {
		return nil, err
	}
	s.metrics.totalVideos.Set(float64(totalVideos))

	// Gather total channels
	totalChannels, err := s.Store.GetTotalChannels(ctx)
	if err != nil {
		return nil, err
	}
	s.metrics.totalChannels.Set(float64(totalChannels))

	// Gather total blocked paths
	totalBlockedPaths, err := s.Store.GetTotalBlockedPaths(ctx)
	if err != nil {
		return nil, err
	}
	s.metrics.totalBlockedPaths.Set(float64(totalBlockedPaths))

	// Gather channel video counts
	channels, err := s.Store.GetChannels(ctx)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		videoCount, err := s.Store.GetTotalVideosByChannelId(ctx, channel.ID)
		if err != nil {
			return nil, err
		}
		s.metrics.channelVideoCount.WithLabelValues(channel.Name).Set(float64(videoCount))
	}

	return s.Registry, nil
}
