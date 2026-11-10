package adapter

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/gen/trajectorypb"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TrajectoryAdapter struct {
	client trajectorypb.TrajectoryServiceClient
}

func NewTrajectoryAdapter(conn *grpc.ClientConn) *TrajectoryAdapter {
	return &TrajectoryAdapter{
		client: trajectorypb.NewTrajectoryServiceClient(conn),
	}
}

func (a *TrajectoryAdapter) GetPosition(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.PositionParams,
) (*models.PositionResult, error) {
	identifier, err := satelliteIdentifier(id)
	if err != nil {
		return nil, err
	}

	req := &trajectorypb.PositionRequest{
		Identifier: identifier,
		Datetime:   timestamppb.New(params.Datetime),
		Units:      unitsToProto(params.Units),
	}

	if len(params.Fields) > 0 {
		req.OutputMask = &fieldmaskpb.FieldMask{
			Paths: params.Fields,
		}
	}

	resp, err := a.client.GetPosition(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("trajectory GetPosition: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("trajectory GetPosition: empty response")
	}

	return &models.PositionResult{
		ComputationMetadata: computationMetadataFromProto(resp.GetMetadata()),
		Time:                timestampToTime(resp.GetTime()),
		Teme:                vectorFromProto(resp.GetTeme()),
		Ecef:                vectorFromProto(resp.GetEcef()),
		Geodetic:            geodeticFromProto(resp.GetGeodetic()),
	}, nil
}

func (a *TrajectoryAdapter) GetLookAngles(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.LookAnglesParams,
) (*models.LookAnglesResult, error) {
	identifier, err := satelliteIdentifier(id)
	if err != nil {
		return nil, err
	}
	req := &trajectorypb.LookAnglesRequest{
		Identifier: identifier,
		Datetime:   timestamppb.New(params.Datetime),
		Observer:   geodeticInputToProto(params.Observer),
		Units:      unitsToProto(params.Units),
	}
	if len(params.Fields) > 0 {
		req.OutputMask = &fieldmaskpb.FieldMask{Paths: params.Fields}
	}
	resp, err := a.client.GetLookAngles(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("trajectory GetLookAngles: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("trajectory GetLookAngles: empty response")
	}
	return &models.LookAnglesResult{
		ComputationMetadata: computationMetadataFromProto(resp.GetMetadata()),
		Time:                timestampToTime(resp.GetTime()),
		Azimuth:             optionalFloat(resp.Azimuth),
		Elevation:           optionalFloat(resp.Elevation),
		Range:               optionalFloat(resp.Range),
	}, nil
}

func (a *TrajectoryAdapter) GetTrajectory(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.TrajectoryParams,
) (*models.TrajectoryResult, error) {
	identifier, err := satelliteIdentifier(id)
	if err != nil {
		return nil, err
	}

	req := &trajectorypb.TrajectoryRequest{
		Identifier: identifier,
		Range:      timeRangeToProto(params.Range),
		Sampling: &trajectorypb.SamplingOptions{
			StepSeconds: uint32(max(params.Sampling.StepSeconds, 0)),
		},
		Units: unitsToProto(params.Units),
	}

	if len(params.Fields) > 0 {
		req.OutputMask = &fieldmaskpb.FieldMask{
			Paths: params.Fields,
		}
	}

	resp, err := a.client.GetTrajectory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("trajectory GetTrajectory: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("trajectory GetTrajectory: empty response")
	}

	states := make([]models.StateVector, 0, len(resp.GetStates()))
	for _, state := range resp.GetStates() {
		if state == nil {
			continue
		}

		states = append(states, models.StateVector{
			Time:     timestampToTime(state.GetDatetime()),
			Teme:     vectorFromProto(state.GetPositionTeme()),
			Ecef:     vectorFromProto(state.GetPositionEcef()),
			Geodetic: geodeticFromProto(state.GetGeodetic()),
		})
	}

	return &models.TrajectoryResult{
		ComputationMetadata: computationMetadataFromProto(resp.GetMetadata()),
		States:              states,
		Range:               timeRangeFromProto(resp.GetRange()),
		Sampling:            samplingFromProto(resp.GetSampling()),
	}, nil
}

func (a *TrajectoryAdapter) GetObserverTrajectory(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.ObserverTrajectoryParams,
) (*models.ObserverTrajectoryResult, error) {
	identifier, err := satelliteIdentifier(id)
	if err != nil {
		return nil, err
	}
	req := &trajectorypb.ObserverTrajectoryRequest{
		Identifier: identifier,
		Range:      timeRangeToProto(params.Range),
		Sampling: &trajectorypb.SamplingOptions{
			StepSeconds: uint32(params.Sampling.StepSeconds),
		},
		Observer: geodeticInputToProto(params.Observer),
		Units:    unitsToProto(params.Units),
	}
	if len(params.Fields) > 0 {
		req.OutputMask = &fieldmaskpb.FieldMask{Paths: params.Fields}
	}
	resp, err := a.client.GetObserverTrajectory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("trajectory GetObserverTrajectory: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("trajectory GetObserverTrajectory: empty response")
	}
	points := make([]models.ObserverTrajectoryPoint, 0, len(resp.GetPoints()))
	for _, point := range resp.GetPoints() {
		if point != nil {
			points = append(points, models.ObserverTrajectoryPoint{
				Time:      timestampToTime(point.GetDatetime()),
				Azimuth:   optionalFloat(point.Azimuth),
				Elevation: optionalFloat(point.Elevation),
				Range:     optionalFloat(point.Range),
			})
		}
	}
	return &models.ObserverTrajectoryResult{
		ComputationMetadata: computationMetadataFromProto(resp.GetMetadata()),
		Points:              points,
		Range:               timeRangeFromProto(resp.GetRange()),
		Sampling:            samplingFromProto(resp.GetSampling()),
	}, nil
}

func (a *TrajectoryAdapter) GetPasses(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.PassesParams,
) (*models.PassesResult, error) {
	identifier, err := satelliteIdentifier(id)
	if err != nil {
		return nil, err
	}

	req := &trajectorypb.PassPredictionRequest{
		Satellites: []*trajectorypb.SatelliteIdentifier{identifier},
		Range:      timeRangeToProto(params.Range),
		Observer:   geodeticInputToProto(params.Observer),
		Units:      unitsToProto(params.Units),
	}

	switch {
	case params.MinPeakElevationDeg != nil:
		req.MinPeakElevation = &trajectorypb.PassPredictionRequest_MinPeakElevationDeg{
			MinPeakElevationDeg: *params.MinPeakElevationDeg,
		}
	case params.MinPeakElevationRad != nil:
		req.MinPeakElevation = &trajectorypb.PassPredictionRequest_MinPeakElevationRad{
			MinPeakElevationRad: *params.MinPeakElevationRad,
		}
	}

	switch {
	case params.MinElevationDeg != nil:
		req.MinElevation = &trajectorypb.PassPredictionRequest_MinElevationDeg{
			MinElevationDeg: *params.MinElevationDeg,
		}
	case params.MinElevationRad != nil:
		req.MinElevation = &trajectorypb.PassPredictionRequest_MinElevationRad{
			MinElevationRad: *params.MinElevationRad,
		}
	}

	if params.MaxResults != nil {
		maxResults := *params.MaxResults
		req.MaxResults = &maxResults
	}

	resp, err := a.client.GetPasses(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("trajectory GetPasses: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("trajectory GetPasses: empty response")
	}

	passes := make([]models.Pass, 0, len(resp.GetPasses()))
	for _, pass := range resp.GetPasses() {
		if pass == nil {
			continue
		}

		passes = append(passes, models.Pass{
			Satellite:           satelliteFromProto(pass.GetSatellite()),
			AOS:                 timestampToTime(pass.GetAos()),
			AOSAzimuth:          pass.GetAosAzimuth(),
			MaxElevationTime:    timestampToTime(pass.GetMaxElevationTime()),
			MaxElevation:        pass.GetMaxElevation(),
			MaxElevationAzimuth: pass.GetMaxElevationAzimuth(),
			LOS:                 timestampToTime(pass.GetLos()),
			LOSAzimuth:          pass.GetLosAzimuth(),
			DurationSeconds:     int(pass.GetDurationSeconds()),
		})
	}

	return &models.PassesResult{
		ComputationMetadata: passMetadataFromProto(resp.GetMetadata()),
		Passes:              passes,
		Range:               timeRangeFromProto(resp.GetRange()),
	}, nil
}

func (a *TrajectoryAdapter) GetNextPasses(
	ctx context.Context,
	id domain.SatelliteIdentifier,
	params domain.NextPassesParams,
) (*models.NextPassesResult, error) {
	identifier, err := satelliteIdentifier(id)
	if err != nil {
		return nil, err
	}
	req := &trajectorypb.NextPassesRequest{
		Satellites: []*trajectorypb.SatelliteIdentifier{identifier},
		Observer:   geodeticInputToProto(params.Observer),
		Units:      unitsToProto(params.Units),
		Count:      params.Count,
	}
	setNextPassElevations(req, params)
	resp, err := a.client.GetNextPasses(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("trajectory GetNextPasses: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("trajectory GetNextPasses: empty response")
	}
	return &models.NextPassesResult{
		ComputationMetadata: passMetadataFromProto(resp.GetMetadata()),
		Passes:              passesFromProto(resp.GetPasses()),
	}, nil
}

func setNextPassElevations(req *trajectorypb.NextPassesRequest, params domain.NextPassesParams) {
	if params.MinPeakElevationDeg != nil {
		req.MinPeakElevation =
			&trajectorypb.NextPassesRequest_MinPeakElevationDeg{
				MinPeakElevationDeg: *params.MinPeakElevationDeg,
			}
	} else if params.MinPeakElevationRad != nil {
		req.MinPeakElevation =
			&trajectorypb.NextPassesRequest_MinPeakElevationRad{
				MinPeakElevationRad: *params.MinPeakElevationRad,
			}
	}
	if params.MinElevationDeg != nil {
		req.MinElevation = &trajectorypb.NextPassesRequest_MinElevationDeg{MinElevationDeg: *params.MinElevationDeg}
	} else if params.MinElevationRad != nil {
		req.MinElevation = &trajectorypb.NextPassesRequest_MinElevationRad{MinElevationRad: *params.MinElevationRad}
	}
}

func passesFromProto(values []*trajectorypb.Pass) []models.Pass {
	passes := make([]models.Pass, 0, len(values))
	for _, pass := range values {
		if pass != nil {
			passes = append(passes, models.Pass{
				Satellite:           satelliteFromProto(pass.GetSatellite()),
				AOS:                 timestampToTime(pass.GetAos()),
				AOSAzimuth:          pass.GetAosAzimuth(),
				MaxElevationTime:    timestampToTime(pass.GetMaxElevationTime()),
				MaxElevation:        pass.GetMaxElevation(),
				MaxElevationAzimuth: pass.GetMaxElevationAzimuth(),
				LOS:                 timestampToTime(pass.GetLos()),
				LOSAzimuth:          pass.GetLosAzimuth(),
				DurationSeconds:     int(pass.GetDurationSeconds()),
			})
		}
	}
	return passes
}

func optionalFloat(value *float64) *float64 {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func satelliteIdentifier(id domain.SatelliteIdentifier) (*trajectorypb.SatelliteIdentifier, error) {
	switch {
	case id.NORADID != nil && id.Name == nil:
		if *id.NORADID <= 0 || uint64(*id.NORADID) > uint64(math.MaxUint32) {
			return nil, fmt.Errorf("invalid NORAD ID: %d", *id.NORADID)
		}

		return &trajectorypb.SatelliteIdentifier{
			Kind: &trajectorypb.SatelliteIdentifier_NoradId{
				NoradId: uint32(*id.NORADID),
			},
		}, nil

	case id.Name != nil && id.NORADID == nil:
		name := strings.TrimSpace(*id.Name)
		if name == "" {
			return nil, fmt.Errorf("satellite name cannot be empty")
		}

		return &trajectorypb.SatelliteIdentifier{
			Kind: &trajectorypb.SatelliteIdentifier_SatelliteName{
				SatelliteName: name,
			},
		}, nil

	case id.NORADID != nil && id.Name != nil:
		return nil, fmt.Errorf("provide either NORAD ID or satellite name, not both")

	default:
		return nil, fmt.Errorf("satellite NORAD ID or name is required")
	}
}

func unitsToProto(units models.UnitSettings) *trajectorypb.UnitSettings {
	result := &trajectorypb.UnitSettings{}

	switch units.DistanceUnit {
	case "meters":
		result.DistanceUnit = trajectorypb.UnitSettings_DISTANCE_UNIT_METERS
	case "kilometers":
		result.DistanceUnit = trajectorypb.UnitSettings_DISTANCE_UNIT_KILOMETERS
	case "miles":
		result.DistanceUnit = trajectorypb.UnitSettings_DISTANCE_UNIT_MILES
	}

	switch units.AngleUnit {
	case "degrees":
		result.AngleUnit = trajectorypb.UnitSettings_ANGLE_UNIT_DEGREES
	case "radians":
		result.AngleUnit = trajectorypb.UnitSettings_ANGLE_UNIT_RADIANS
	}

	return result
}

func timeRangeToProto(r models.TimeRange) *trajectorypb.TimeRange {
	return &trajectorypb.TimeRange{
		Start: timestamppb.New(r.Start),
		End:   timestamppb.New(r.End),
	}
}

func geodeticInputToProto(input domain.GeodeticInput) *trajectorypb.GeodeticInput {
	result := &trajectorypb.GeodeticInput{}

	switch {
	case input.LatDeg != nil:
		result.Lat = &trajectorypb.GeodeticInput_LatDeg{LatDeg: *input.LatDeg}
	case input.LatRad != nil:
		result.Lat = &trajectorypb.GeodeticInput_LatRad{LatRad: *input.LatRad}
	}

	switch {
	case input.LonDeg != nil:
		result.Lon = &trajectorypb.GeodeticInput_LonDeg{LonDeg: *input.LonDeg}
	case input.LonRad != nil:
		result.Lon = &trajectorypb.GeodeticInput_LonRad{LonRad: *input.LonRad}
	}

	switch {
	case input.AltM != nil:
		result.Alt = &trajectorypb.GeodeticInput_AltM{AltM: *input.AltM}
	case input.AltKm != nil:
		result.Alt = &trajectorypb.GeodeticInput_AltKm{AltKm: *input.AltKm}
	}

	return result
}

func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

func vectorFromProto(v *trajectorypb.Vector3) *models.Vector3 {
	if v == nil {
		return nil
	}
	return &models.Vector3{
		X: v.GetX(),
		Y: v.GetY(),
		Z: v.GetZ(),
	}
}

func geodeticFromProto(g *trajectorypb.GeodeticOutput) *models.Geodetic {
	if g == nil {
		return nil
	}
	return &models.Geodetic{
		Lat: g.GetLat(),
		Lon: g.GetLon(),
		Alt: g.GetAlt(),
	}
}

func computationMetadataFromProto(m *trajectorypb.TrajectoryComputationMetadata) models.TrajectoryComputationMetadata {
	if m == nil {
		return models.TrajectoryComputationMetadata{}
	}

	return models.TrajectoryComputationMetadata{
		PropagationModel: m.GetPropagationModel(),
		ComputationTime:  timestampToTime(m.GetComputationTime()),
		NoradID:          int64(m.GetNoradId()),
		SatelliteName:    m.GetSatelliteName(),
		TleEpoch:         timestampToTime(m.GetTleEpoch()),
		Units:            unitsFromProto(m.GetUnits()),
	}
}

func unitsFromProto(u *trajectorypb.UnitSettings) models.UnitSettings {
	if u == nil {
		return models.UnitSettings{}
	}

	var result models.UnitSettings

	switch u.GetDistanceUnit() {
	case trajectorypb.UnitSettings_DISTANCE_UNIT_METERS:
		result.DistanceUnit = "meters"
	case trajectorypb.UnitSettings_DISTANCE_UNIT_KILOMETERS:
		result.DistanceUnit = "kilometers"
	case trajectorypb.UnitSettings_DISTANCE_UNIT_MILES:
		result.DistanceUnit = "miles"
	}

	switch u.GetAngleUnit() {
	case trajectorypb.UnitSettings_ANGLE_UNIT_DEGREES:
		result.AngleUnit = "degrees"
	case trajectorypb.UnitSettings_ANGLE_UNIT_RADIANS:
		result.AngleUnit = "radians"
	}

	return result
}

func timeRangeFromProto(r *trajectorypb.TimeRange) models.TimeRange {
	if r == nil {
		return models.TimeRange{}
	}

	return models.TimeRange{
		Start: timestampToTime(r.GetStart()),
		End:   timestampToTime(r.GetEnd()),
	}
}

func samplingFromProto(s *trajectorypb.SamplingOptions) models.SamplingOptions {
	if s == nil {
		return models.SamplingOptions{}
	}

	return models.SamplingOptions{
		StepSeconds: int(s.GetStepSeconds()),
	}
}

func satelliteFromProto(s *trajectorypb.SatelliteIdentifier) models.SatelliteIdentifier {
	if s == nil {
		return models.SatelliteIdentifier{}
	}

	result := models.SatelliteIdentifier{}

	switch value := s.GetKind().(type) {
	case *trajectorypb.SatelliteIdentifier_NoradId:
		id := int64(value.NoradId)
		result.NoradID = &id
	case *trajectorypb.SatelliteIdentifier_SatelliteName:
		name := value.SatelliteName
		result.Name = &name
	}

	return result
}

func passMetadataFromProto(m *trajectorypb.PassesComputationMetadata) models.PassComputationMetadata {
	if m == nil {
		return models.PassComputationMetadata{}
	}

	noradIDs := make([]int64, 0, len(m.GetNoradIds()))
	for _, id := range m.GetNoradIds() {
		noradIDs = append(noradIDs, int64(id))
	}

	return models.PassComputationMetadata{
		PropagationModel:    m.GetPropagationModel(),
		ComputationTime:     timestampToTime(m.GetComputationTime()),
		NoradIDs:            noradIDs,
		SatelliteNames:      append([]string(nil), m.GetSatelliteNames()...),
		TleEpoch:            timestampToTime(m.GetTleEpoch()),
		Units:               unitsFromProto(m.GetUnits()),
		SatellitesEvaluated: int(m.GetSatellitesEvaluated()),
		PassesFound:         int(m.GetPassesFound()),
	}
}
