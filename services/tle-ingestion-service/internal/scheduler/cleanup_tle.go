package scheduler

import (
	"context"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
)

type TLECleanupScheduler struct {
	tleRepo         domain.TLERepository
	logger          applog.Logger
	interval        time.Duration
	cleanupTimeout  time.Duration
	retentionPeriod time.Duration
}

func NewTLECleanupScheduler(
	r domain.TLERepository,
	logger applog.Logger,
	interval,
	contextTimeout,
	retentionPeriod time.Duration,
) *TLECleanupScheduler {

	return &TLECleanupScheduler{
		tleRepo:         r,
		logger:          logger,
		interval:        interval,
		cleanupTimeout:  contextTimeout,
		retentionPeriod: retentionPeriod,
	}
}

func (s *TLECleanupScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("starting TLE cleanup scheduler", applog.NewField("interval", s.interval))

	s.cleanup(ctx)

	for {
		select {
		case <-ticker.C:
			s.cleanup(ctx)
		case <-ctx.Done():
			s.logger.Info("TLE cleanup scheduler stopped", applog.NewErrorField(ctx.Err()))
			return
		}
	}
}

func (s *TLECleanupScheduler) cleanup(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, s.cleanupTimeout)
	defer cancel()

	err := s.tleRepo.DeleteOlderThan(ctx, s.retentionPeriod)
	if err != nil {
		s.logger.Error("failed to delete old TLEs", applog.NewErrorField(err))
		return
	}

	s.logger.Info("deleted old TLEs successfully")
}
