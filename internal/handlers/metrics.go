package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type MetricsService interface {
	GatherMetrics() (*prometheus.Registry, error)
}

func (h *Handler) GatherMetrics() (*prometheus.Registry, error) {
	r, err := h.Services.MetricsService.GatherMetrics()
	if err != nil {
		log.Error().Err(err).Msg("error gathering metrics")
		return nil, err
	}
	return r, nil
}
