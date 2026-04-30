package fallback

import (
	"context"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
)

type Source struct {
	primary   ingestion.MetadataSource
	secondary ingestion.MetadataSource
	logger    applog.Logger
}

func NewFallbackSource(primary, secondary ingestion.MetadataSource, logger applog.Logger) *Source {
	return &Source{
		primary:   primary,
		secondary: secondary,
		logger:    logger,
	}
}

func (s *Source) Name() string {
	return s.primary.Name() + "_fallback"
}

func (s *Source) StreamBatch(ctx context.Context, fn func([]*ingestion.SatelliteSourceRecord) error) error {
	err := s.primary.StreamBatch(ctx, fn)
	if err == nil {
		return nil
	}

	s.logger.Error("primary source failed, switching to fallback",
		applog.NewField("primary", s.primary.Name()),
		applog.NewField("fallback", s.secondary.Name()),
		applog.NewErrorField(err),
	)

	return s.secondary.StreamBatch(ctx, fn)
}
