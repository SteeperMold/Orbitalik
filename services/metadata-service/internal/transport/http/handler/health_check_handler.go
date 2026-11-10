package handler

import (
	"context"
	"net/http"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/domain"
)

type HealthCheckHandler struct {
	service        domain.HealthCheckService
	logger         applog.Logger
	contextTimeout time.Duration
}

func NewHealthHandler(s domain.HealthCheckService, logger applog.Logger, timeout time.Duration) *HealthCheckHandler {
	return &HealthCheckHandler{
		service:        s,
		logger:         logger,
		contextTimeout: timeout,
	}
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.contextTimeout)
	defer cancel()

	err := h.service.HealthCheck(ctx)
	if err != nil {
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
