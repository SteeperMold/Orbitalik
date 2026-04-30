package scheduler

import (
	"context"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/domain"
)

type FetchMetadataScheduler struct {
	service        domain.IngestionService
	logger         applog.Logger
	interval       time.Duration
	contextTimeout time.Duration
}

func NewFetchMetadataScheduler(
	s domain.IngestionService,
	logger applog.Logger,
	interval,
	contextTimeout time.Duration,
) *FetchMetadataScheduler {

	return &FetchMetadataScheduler{
		service:        s,
		logger:         logger,
		interval:       interval,
		contextTimeout: contextTimeout,
	}
}

func (s *FetchMetadataScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("starting metadata fetch scheduler", applog.NewField("interval", s.interval))

	s.runFetch(ctx)

	for {
		select {
		case <-ticker.C:
			s.runFetch(ctx)
		case <-ctx.Done():
			s.logger.Info("metadata fetch scheduler stopped", applog.NewErrorField(ctx.Err()))
			return
		}
	}
}

func (s *FetchMetadataScheduler) runFetch(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := s.service.IngestMetadata(ctx)
	if err != nil {
		s.logger.Error("error ingesting metadata", applog.NewErrorField(err))
		return
	}

	s.logger.Info("ingested metadata")
}
