package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type IngestionService interface {
	IngestMetadata(ctx context.Context) error
}

type MetadataService interface {
	GetMetadataByNoradID(ctx context.Context, noradID int) (*models.SatelliteMetadata, error)
	GetMetadataByName(ctx context.Context, name string) (*models.SatelliteMetadata, error)
	ListSatellites(ctx context.Context, filter *models.ListFilter) ([]*models.SatelliteMetadata, string, uint32, error)
}

type MetadataRepository interface {
	GetMetadataByNoradID(ctx context.Context, noradID int) (*models.SatelliteMetadata, error)
	GetMetadataByName(ctx context.Context, name string) (*models.SatelliteMetadata, error)
	ListSatellites(ctx context.Context, filter *models.ListFilter) ([]*models.SatelliteMetadata, string, uint32, error)
}
