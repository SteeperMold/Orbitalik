package domain

import (
	"context"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
)

type SatelliteIdentifier struct {
	NORADID *int64
	Name    *string
}

type PositionParams struct {
	Datetime time.Time
	Fields   []string
	Units    models.UnitSettings
}

type TrajectoryParams struct {
	Range    models.TimeRange
	Sampling models.SamplingOptions
	Fields   []string
	Units    models.UnitSettings
}

type LookAnglesParams struct {
	Datetime time.Time
	Observer GeodeticInput
	Fields   []string
	Units    models.UnitSettings
}

type ObserverTrajectoryParams struct {
	Range    models.TimeRange
	Sampling models.SamplingOptions
	Observer GeodeticInput
	Fields   []string
	Units    models.UnitSettings
}

type PassesParams struct {
	Range    models.TimeRange
	Observer GeodeticInput
	Units    models.UnitSettings

	MinPeakElevationDeg *float64
	MinPeakElevationRad *float64

	MinElevationDeg *float64
	MinElevationRad *float64

	MaxResults *uint32
}

type NextPassesParams struct {
	Observer            GeodeticInput
	Units               models.UnitSettings
	MinPeakElevationDeg *float64
	MinPeakElevationRad *float64
	MinElevationDeg     *float64
	MinElevationRad     *float64
	Count               uint32
}

// GeodeticInput matches the protobuf GeodeticInput.
// Only one option should be present for each coordinate
type GeodeticInput struct {
	LatDeg *float64
	LatRad *float64

	LonDeg *float64
	LonRad *float64

	AltM  *float64
	AltKm *float64
}

type SatelliteService interface {
	List(ctx context.Context, pageSize uint32, pageToken string) ([]models.Satellite, string, error)
	Get(ctx context.Context, id SatelliteIdentifier) (*models.Satellite, error)

	GetMetadata(ctx context.Context, id SatelliteIdentifier) (*models.SatelliteMetadata, error)
	GetFrequencies(ctx context.Context, id SatelliteIdentifier) (models.Frequencies, error)
	GetSources(ctx context.Context, id SatelliteIdentifier) (models.Sources, error)

	GetTle(ctx context.Context, id SatelliteIdentifier) (*models.Tle, error)
	ListTles(ctx context.Context) ([]models.Tle, error)

	GetPosition(ctx context.Context, id SatelliteIdentifier, params PositionParams) (*models.PositionResult, error)
	GetLookAngles(ctx context.Context, id SatelliteIdentifier, params LookAnglesParams) (*models.LookAnglesResult, error)
	GetTrajectory(ctx context.Context, id SatelliteIdentifier, params TrajectoryParams) (*models.TrajectoryResult, error)
	GetObserverTrajectory(ctx context.Context, id SatelliteIdentifier, params ObserverTrajectoryParams) (*models.ObserverTrajectoryResult, error)
	GetPasses(ctx context.Context, id SatelliteIdentifier, params PassesParams) (*models.PassesResult, error)
	GetNextPasses(ctx context.Context, id SatelliteIdentifier, params NextPassesParams) (*models.NextPassesResult, error)
}
