package aggregation

import (
	"context"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
)

type RawRepository interface {
	GetByNoradID(ctx context.Context, noradID int) ([]models.SatelliteIngestRecord, error)
}

type MetadataRepository interface {
	Upsert(ctx context.Context, meta *models.SatelliteMetadata) error
}

type Service struct {
	rawRepo  RawRepository
	metaRepo MetadataRepository
}

func NewService(rawRepo RawRepository, metaRepo MetadataRepository) *Service {
	return &Service{
		rawRepo:  rawRepo,
		metaRepo: metaRepo,
	}
}

func (s *Service) RebuildSatellite(ctx context.Context, noradID int) error {
	records, err := s.rawRepo.GetByNoradID(ctx, noradID)
	if err != nil {
		return err
	}

	meta := Aggregate(records)

	if meta == nil {
		return nil
	}

	return s.metaRepo.Upsert(ctx, meta)
}
