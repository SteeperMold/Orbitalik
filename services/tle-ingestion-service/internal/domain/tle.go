package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"
)

type FetchTLEService interface {
	FetchTLE(ctx context.Context) error
}

type TLEService interface {
	GetAllTLEs(ctx context.Context) ([]*models.TLE, error)
	GetTLEByNoradID(ctx context.Context, noradID int) (*models.TLE, error)
	GetTLEBySatelliteName(ctx context.Context, name string) (*models.TLE, error)
}

type TLERepository interface {
	SaveBatch(ctx context.Context, tles []*models.TLE) error
	GetAllTLEs(ctx context.Context) ([]*models.TLE, error)
	GetTLEByNoradID(ctx context.Context, noradID int) (*models.TLE, error)
	GetTLEBySatelliteName(ctx context.Context, name string) (*models.TLE, error)
}
