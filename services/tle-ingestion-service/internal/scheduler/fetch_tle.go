package scheduler

import (
	"context"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
)

type TLEFetchScheduler struct {
	service        domain.FetchTLEService
	logger         applog.Logger
	interval       time.Duration
	contextTimeout time.Duration
}

func NewTLEFetchScheduler(
	s domain.FetchTLEService,
	logger applog.Logger,
	interval,
	contextTimeout time.Duration,
) *TLEFetchScheduler {

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

	s.logger.Info("starting TLE fetch scheduler", applog.NewField("interval", s.interval))

	s.runFetch(ctx)

	for {
		select {
		case <-ticker.C:
			s.runFetch(ctx)
		case <-ctx.Done():
			s.logger.Info("TLE fetch scheduler stopped", applog.NewErrorField(ctx.Err()))
			return
		}
	}
}

func (s *TLEFetchScheduler) runFetch(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err := s.service.FetchTLE(ctx)
	if err != nil {
		s.logger.Error("error fetching TLEs", applog.NewErrorField(err))
		return
	}

	s.logger.Info("TLE fetch completed successfully")
}
