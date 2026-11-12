package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
)

type SatelliteService struct {
	metadataClient   domain.MetadataClient
	tleClient        domain.TleClient
	trajectoryClient domain.TrajectoryClient
}

func NewSatelliteService(
	metadataClient domain.MetadataClient,
	tleClient domain.TleClient,
	trajectoryClient domain.TrajectoryClient,
) *SatelliteService {
	return &SatelliteService{
		metadataClient:   metadataClient,
		tleClient:        tleClient,
		trajectoryClient: trajectoryClient,
	}
}

func (s *SatelliteService) List(
	ctx context.Context,
	pageSize uint32,
	pageToken string,
) ([]models.Satellite, string, error) {
	if pageSize == 0 {
		return nil, "", domain.NewValidationError(
			"page size must be greater than zero",
		)
	}

	return s.metadataClient.List(ctx, pageSize, pageToken)
}

func (s *SatelliteService) Get(
	ctx context.Context,
	id domain.SatelliteIdentifier,
) (*models.Satellite, error) {
	metadata, err := s.GetMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Satellite{
		NoradID:           metadata.NoradID,
		CosparID:          stringValue(metadata.CosparID),
		Name:              metadata.Name,
		ObjectType:        metadata.ObjectType,
		MissionType:       metadata.MissionType,
		OrbitRegime:       metadata.OrbitRegime,
		Operator:          stringValue(metadata.Operator),
		OperationalStatus: metadata.OperationalStatus,
	}, nil
}

func (s *SatelliteService) GetMetadata(
	ctx context.Context,
	id domain.SatelliteIdentifier,
) (*models.SatelliteMetadata, error) {
	return s.metadataClient.GetMetadata(ctx, id)
}

func (s *SatelliteService) GetFrequencies(
	ctx context.Context,
	id domain.SatelliteIdentifier,
) (models.Frequencies, error) {
	metadata, err := s.GetMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	return metadata.Frequencies, nil
}

func (s *SatelliteService) GetSources(
	ctx context.Context,
	id domain.SatelliteIdentifier,
) (models.Sources, error) {
	metadata, err := s.GetMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	return metadata.Sources, nil
}

func (s *SatelliteService) GetTle(
	ctx context.Context,
	id domain.SatelliteIdentifier,
) (*models.Tle, error) {
	return s.tleClient.GetTle(ctx, id)
}

func (s *SatelliteService) ListTles(ctx context.Context) ([]models.Tle, error) {
	return s.tleClient.ListTles(ctx)
}

func (s *SatelliteService) GetPosition(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.PositionParams,
) (*models.PositionResult, error) {
	if err := validatePositionParams(params); err != nil {
		return nil, err
	}

	return s.trajectoryClient.GetPosition(ctx, id, params)
}

func (s *SatelliteService) GetLookAngles(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.LookAnglesParams,
) (*models.LookAnglesResult, error) {
	if err := validateLookAnglesParams(params); err != nil {
		return nil, err
	}
	return s.trajectoryClient.GetLookAngles(ctx, id, params)
}

func (s *SatelliteService) GetTrajectory(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.TrajectoryParams,
) (*models.TrajectoryResult, error) {
	if params.Sampling.StepSeconds == 0 {
		params.Sampling.StepSeconds = defaultTrajectoryStepSeconds
	}

	if err := validateTrajectoryParams(params); err != nil {
		return nil, err
	}

	return s.trajectoryClient.GetTrajectory(ctx, id, params)
}

func (s *SatelliteService) GetObserverTrajectory(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.ObserverTrajectoryParams,
) (*models.ObserverTrajectoryResult, error) {
	if params.Sampling.StepSeconds == 0 {
		params.Sampling.StepSeconds = defaultTrajectoryStepSeconds
	}
	if err := validateObserverTrajectoryParams(params); err != nil {
		return nil, err
	}
	return s.trajectoryClient.GetObserverTrajectory(ctx, id, params)
}

func (s *SatelliteService) GetPasses(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.PassesParams,
) (*models.PassesResult, error) {
	if err := validatePassesParams(params); err != nil {
		return nil, err
	}

	return s.trajectoryClient.GetPasses(ctx, id, params)
}

func (s *SatelliteService) GetNextPasses(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.NextPassesParams,
) (*models.NextPassesResult, error) {
	if err := validateNextPassesParams(params); err != nil {
		return nil, err
	}
	return s.trajectoryClient.GetNextPasses(ctx, id, params)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
