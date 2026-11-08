package adapter

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/gen/metadatapb"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetadataAdapter struct {
	client metadatapb.SatelliteMetadataServiceClient
}

func NewMetadataAdapter(conn *grpc.ClientConn) *MetadataAdapter {
	return &MetadataAdapter{
		client: metadatapb.NewSatelliteMetadataServiceClient(conn),
	}
}

func (a *MetadataAdapter) GetMetadata(
	ctx context.Context,
	id domain.SatelliteIdentifier,
) (*models.SatelliteMetadata, error) {
	identifier, err := satelliteIdentifierToMetadataProto(id)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.GetSatelliteMetadata(
		ctx,
		&metadatapb.GetMetadataRequest{
			Identifier: identifier,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get satellite metadata (NORAD %d): %w", id, err)
	}
	if resp == nil || resp.GetMetadata() == nil {
		return nil, fmt.Errorf("metadata not found for NORAD %d", id)
	}

	return mapSatelliteMetadata(resp.GetMetadata()), nil
}

func (a *MetadataAdapter) List(ctx context.Context, pageSize uint32, pageToken string) ([]models.Satellite, string, error) {
	var satellites []models.Satellite

	resp, err := a.client.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{
			PageSize:  pageSize,
			PageToken: pageToken,
		},
	)
	if err != nil {
		return nil, "", fmt.Errorf("list satellite metadata: %w", err)
	}
	if resp == nil {
		return nil, "", fmt.Errorf("list satellite metadata: empty response")
	}

	for _, item := range resp.GetItems() {
		if item == nil {
			continue
		}

		satellites = append(satellites, models.Satellite{
			NoradID:           int64(item.GetNoradId()),
			CosparID:          item.GetCosparId(),
			Name:              item.GetName(),
			ObjectType:        enumValue(item.GetObjectType().String(), "OBJECT_TYPE_"),
			MissionType:       enumValue(item.GetMissionType().String(), "MISSION_TYPE_"),
			OrbitRegime:       enumValue(item.GetOrbitRegime().String(), "ORBIT_REGIME_"),
			Operator:          item.GetOperator(),
			OperationalStatus: enumValue(item.GetOperationalStatus().String(), "OPERATIONAL_STATUS_"),
		})
	}

	next := resp.GetNextPageToken()

	return satellites, next, nil
}

func satelliteIdentifierToMetadataProto(id domain.SatelliteIdentifier) (*metadatapb.SatelliteIdentifier, error) {
	switch {
	case id.NORADID != nil && id.Name == nil:
		if *id.NORADID <= 0 || uint64(*id.NORADID) > uint64(math.MaxUint32) {
			return nil, fmt.Errorf("invalid NORAD ID: %d", *id.NORADID)
		}

		return &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_NoradId{
				NoradId: uint32(*id.NORADID),
			},
		}, nil

	case id.Name != nil && id.NORADID == nil:
		name := strings.TrimSpace(*id.Name)
		if name == "" {
			return nil, fmt.Errorf("satellite name cannot be empty")
		}

		return &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_SatelliteName{
				SatelliteName: name,
			},
		}, nil

	case id.NORADID != nil && id.Name != nil:
		return nil, fmt.Errorf(
			"provide either NORAD ID or satellite name, not both",
		)

	default:
		return nil, fmt.Errorf("satellite NORAD ID or name is required")
	}
}

func mapSatelliteMetadata(src *metadatapb.SatelliteMetadata) *models.SatelliteMetadata {
	dst := &models.SatelliteMetadata{
		NoradID:           int64(src.GetNoradId()),
		CosparID:          optionalString(src.CosparId),
		Name:              src.GetName(),
		Aliases:           append([]string(nil), src.GetAliases()...),
		ObjectType:        enumValue(src.GetObjectType().String(), "OBJECT_TYPE_"),
		MissionType:       enumValue(src.GetMissionType().String(), "MISSION_TYPE_"),
		OrbitRegime:       enumValue(src.GetOrbitRegime().String(), "ORBIT_REGIME_"),
		Operator:          optionalString(src.Operator),
		Owner:             optionalString(src.Owner),
		Constellation:     optionalString(src.Constellation),
		LaunchDate:        optionalTime(src.LaunchDate),
		LaunchSite:        optionalString(src.LaunchSite),
		LaunchVehicle:     optionalString(src.LaunchVehicle),
		OperationalStatus: enumValue(src.GetOperationalStatus().String(), "OPERATIONAL_STATUS_"),
		UpdatedAt:         timestamp(src.GetUpdatedAt()),
	}

	for _, frequency := range src.GetFrequencies() {
		if frequency == nil {
			continue
		}

		bandwidth := frequency.GetBandwidthKhz()
		dst.Frequencies = append(dst.Frequencies, models.Frequency{
			Direction:    enumValue(frequency.GetDirection().String(), "FREQUENCY_DIRECTION_"),
			FrequencyMHz: frequency.GetFrequencyMhz(),
			BandwidthKHz: &bandwidth,
			Modulation:   frequency.GetModulation(),
			Mode:         frequency.GetMode(),
		})
	}

	for _, source := range src.GetSources() {
		if source == nil {
			continue
		}

		names := make([]string, 0, len(source.GetSources()))
		for _, value := range source.GetSources() {
			names = append(names, enumValue(value.String(), "SOURCE_"))
		}

		dst.Sources = append(dst.Sources, models.FieldSource{
			Field:     source.GetField(),
			Sources:   names,
			FetchedAt: timestamp(source.GetFetchedAt()),
		})
	}

	return dst
}

func enumValue(value, prefix string) string {
	value = strings.TrimPrefix(value, prefix)
	value = strings.ToLower(value)
	if value == "unspecified" {
		return ""
	}
	return value
}

func optionalString(value *string) *string {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}

func optionalTime(value *timestamppb.Timestamp) *time.Time {
	if value == nil || !value.IsValid() {
		return nil
	}
	result := value.AsTime()
	return &result
}

func timestamp(value *timestamppb.Timestamp) time.Time {
	if value == nil || !value.IsValid() {
		return time.Time{}
	}
	return value.AsTime()
}
