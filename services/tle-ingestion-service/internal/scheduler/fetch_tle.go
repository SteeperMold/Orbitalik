package scheduler

import (
	"context"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"go.uber.org/zap"
	"time"
)

type TLEFetchScheduler struct {
	service        domain.FetchTLEService
	logger         *zap.Logger
	interval       time.Duration
	contextTimeout time.Duration
}

func NewTLEFetchScheduler(s domain.FetchTLEService, logger *zap.Logger, interval, contextTimeout time.Duration) *TLEFetchScheduler {
	return &TLEFetchScheduler{
		service:        s,
		logger:         logger,
		interval:       interval,
		contextTimeout: contextTimeout,
	}
}

func (s *TLEFetchScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("starting TLE fetch scheduler", zap.Duration("interval", s.interval))

	s.runFetch(ctx)

	for {
		select {
		case <-ticker.C:
			s.runFetch(ctx)
		case <-ctx.Done():
			s.logger.Info("TLE fetch scheduler stopped", zap.Error(ctx.Err()))
			return
		}
	}
}

func (s *TLEFetchScheduler) runFetch(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := s.service.FetchTLE(ctx)
	if err != nil {
		s.logger.Error("error fetching TLEs", zap.Error(err))
		return
	}

	s.logger.Info("TLE fetch completed successfully")
}
