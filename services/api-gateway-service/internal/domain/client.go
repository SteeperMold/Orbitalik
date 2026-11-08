package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
)

type TrajectoryClient interface {
	GetPosition(ctx context.Context, id SatelliteIdentifier, params PositionParams) (*models.PositionResult, error)
	GetLookAngles(ctx context.Context, id SatelliteIdentifier, params LookAnglesParams) (*models.LookAnglesResult, error)
	GetTrajectory(ctx context.Context, id SatelliteIdentifier, params TrajectoryParams) (*models.TrajectoryResult, error)
	GetObserverTrajectory(ctx context.Context, id SatelliteIdentifier, params ObserverTrajectoryParams) (*models.ObserverTrajectoryResult, error)
	GetPasses(ctx context.Context, id SatelliteIdentifier, params PassesParams) (*models.PassesResult, error)
	GetNextPasses(ctx context.Context, id SatelliteIdentifier, params NextPassesParams) (*models.NextPassesResult, error)
}

type TleClient interface {
	GetTle(ctx context.Context, id SatelliteIdentifier) (*models.Tle, error)
	ListTles(ctx context.Context) ([]models.Tle, error)
}

type MetadataClient interface {
	GetMetadata(ctx context.Context, id SatelliteIdentifier) (*models.SatelliteMetadata, error)
	List(ctx context.Context, pageSize uint32, pageToken string) ([]models.Satellite, string, error)
}
