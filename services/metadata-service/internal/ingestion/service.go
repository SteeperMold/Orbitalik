package ingestion

import (
	"context"
	"encoding/json"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type MetadataSource interface {
	StreamBatch(ctx context.Context, fn func([]*SatelliteSourceRecord) error) error
	Name() string
}

type RawMetadataRepository interface {
	SaveRawBatch(ctx context.Context, data []*models.SatelliteIngestRecord) error
}

type DirtyMarker interface {
	MarkDirty(ctx context.Context, noradIDs []int) error
}

type Service struct {
	sources []MetadataSource
	repo    RawMetadataRepository
	dirty   DirtyMarker
	logger  applog.Logger
}

func NewService(sources []MetadataSource, repo RawMetadataRepository, dirty DirtyMarker, logger applog.Logger) *Service {
	return &Service{
		sources: sources,
		repo:    repo,
		dirty:   dirty,
		logger:  logger,
	}
}

func (s *Service) IngestMetadata(ctx context.Context) error {
	for _, src := range s.sources {
		if err := s.ingestSource(ctx, src); err != nil {
			s.logger.Error("failed to ingest source",
				applog.NewField("source", src.Name()),
				applog.NewErrorField(err),
			)
		}
	}
	return nil
}

func (s *Service) ingestSource(ctx context.Context, src MetadataSource) error {
	var total int

	s.logger.Info("starting ingestion", applog.NewField("source", src.Name()))

	err := src.StreamBatch(ctx, func(batch []*SatelliteSourceRecord) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if len(batch) == 0 {
			return nil
		}

		ingestBatch := make([]*models.SatelliteIngestRecord, 0, len(batch))
		noradIDs := make([]int, 0, len(batch))

		for _, r := range batch {
			rec, err := toIngestRecord(r)
			if err != nil {
				continue
			}

			ingestBatch = append(ingestBatch, rec)
			noradIDs = append(noradIDs, rec.NoradID)
		}

		total += len(batch)

		if err := s.repo.SaveRawBatch(ctx, ingestBatch); err != nil {
			return err
		}

		if err := s.dirty.MarkDirty(ctx, noradIDs); err != nil {
			s.logger.Error("failed to mark satellites dirty", applog.NewErrorField(err))
		}

		return nil
	})

	if err != nil {
		s.logger.Error("ingestion failed",
			applog.NewField("source", src.Name()),
			applog.NewErrorField(err),
		)
		return err
	}

	s.logger.Info("ingestion completed",
		applog.NewField("source", src.Name()),
		applog.NewField("records", total),
	)

	return nil
}

func toIngestRecord(src *SatelliteSourceRecord) (*models.SatelliteIngestRecord, error) {
	var payload struct {
		NoradID  int
		CosparID *string
	}

	if err := json.Unmarshal(src.Raw, &payload); err != nil {
		return nil, err
	}

	now := time.Now()

	return &models.SatelliteIngestRecord{
		NoradID:  payload.NoradID,
		CosparID: payload.CosparID,

		Source:  src.Source,
		Payload: src.Raw,

		FetchedAt: src.FetchedAt,
		StoredAt:  now,
	}, nil
}
