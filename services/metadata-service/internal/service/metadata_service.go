package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type MetadataRepository interface {
	GetMetadataByNoradID(ctx context.Context, noradID int) (*models.SatelliteMetadata, error)
	GetMetadataByName(ctx context.Context, name string) (*models.SatelliteMetadata, error)
	ListSatellites(ctx context.Context, filter *models.ListFilter) ([]*models.SatelliteMetadata, string, uint32, error)
}

type MetadataService struct {
	repo MetadataRepository
}

func NewMetadataService(repo MetadataRepository) *MetadataService {
	return &MetadataService{
		repo: repo,
	}
}

func (s *MetadataService) GetMetadataByNoradID(ctx context.Context, noradID int) (*models.SatelliteMetadata, error) {
	return s.repo.GetMetadataByNoradID(ctx, noradID)
}

func (s *MetadataService) GetMetadataByName(ctx context.Context, name string) (*models.SatelliteMetadata, error) {
	return s.repo.GetMetadataByName(ctx, name)
}

func (s *MetadataService) ListSatellites(
	ctx context.Context,
	filter *models.ListFilter,
) ([]*models.SatelliteMetadata, string, uint32, error) {

	return s.repo.ListSatellites(ctx, filter)
}
