package grpc

import (
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/gen/metadatapb"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoSatelliteMetadata(m *models.SatelliteMetadata) *metadatapb.SatelliteMetadata {
	if m == nil {
		return nil
	}

	pb := &metadatapb.SatelliteMetadata{
		// #nosec G115 -- user.ID is a postgres INTEGER and fits in uint32
		NoradId:  uint32(m.NoradID),
		CosparId: m.CosparID,

		Name:    m.Name,
		Aliases: m.Aliases,

		ObjectType:  mapObjectType(m.ObjectType),
		MissionType: mapMissionType(m.MissionType),
		OrbitRegime: mapOrbitRegime(m.OrbitRegime),

		Operator:      m.Operator,
		Owner:         m.Owner,
		Constellation: m.Constellation,

		LaunchSite:    m.LaunchSite,
		LaunchVehicle: m.LaunchVehicle,

		OperationalStatus: mapOperationalStatus(m.OperationalStatus),
		Frequencies:       mapFrequencies(m.Frequencies),
		Sources:           mapSources(m.Sources),
	}

	if m.LaunchDate != nil {
		pb.LaunchDate = timestamppb.New(*m.LaunchDate)
	}
	if !m.UpdatedAt.IsZero() {
		pb.UpdatedAt = timestamppb.New(m.UpdatedAt)
	}

	return pb
}

func mapObjectType(t models.ObjectType) metadatapb.ObjectType {
	switch t {
	case models.ObjectTypePayload:
		return metadatapb.ObjectType_OBJECT_TYPE_PAYLOAD
	case models.ObjectTypeRocketBody:
		return metadatapb.ObjectType_OBJECT_TYPE_ROCKET_BODY
	case models.ObjectTypeDebris:
		return metadatapb.ObjectType_OBJECT_TYPE_DEBRIS
	default:
		return metadatapb.ObjectType_OBJECT_TYPE_UNSPECIFIED
	}
}

func mapMissionType(t models.MissionType) metadatapb.MissionType {
	switch t {
	case models.MissionTypeCommunications:
		return metadatapb.MissionType_MISSION_TYPE_COMMUNICATIONS
	case models.MissionTypeEarthObservation:
		return metadatapb.MissionType_MISSION_TYPE_EARTH_OBSERVATION
	case models.MissionTypeNavigation:
		return metadatapb.MissionType_MISSION_TYPE_NAVIGATION
	case models.MissionTypeScience:
		return metadatapb.MissionType_MISSION_TYPE_SCIENCE
	case models.MissionTypeWeather:
		return metadatapb.MissionType_MISSION_TYPE_WEATHER
	case models.MissionTypeAmateur:
		return metadatapb.MissionType_MISSION_TYPE_AMATEUR
	case models.MissionTypeTechDemo:
		return metadatapb.MissionType_MISSION_TYPE_TECH_DEMO
	default:
		return metadatapb.MissionType_MISSION_TYPE_UNSPECIFIED
	}
}

func mapOrbitRegime(r models.OrbitRegime) metadatapb.OrbitRegime {
	switch r {
	case models.OrbitRegimeLEO:
		return metadatapb.OrbitRegime_ORBIT_REGIME_LEO
	case models.OrbitRegimeMEO:
		return metadatapb.OrbitRegime_ORBIT_REGIME_MEO
	case models.OrbitRegimeGEO:
		return metadatapb.OrbitRegime_ORBIT_REGIME_GEO
	case models.OrbitRegimeHEO:
		return metadatapb.OrbitRegime_ORBIT_REGIME_HEO
	default:
		return metadatapb.OrbitRegime_ORBIT_REGIME_UNSPECIFIED
	}
}

func mapOperationalStatus(s models.OperationalStatus) metadatapb.OperationalStatus {
	switch s {
	case models.OperationalStatusActive:
		return metadatapb.OperationalStatus_OPERATIONAL_STATUS_ACTIVE
	case models.OperationalStatusInactive:
		return metadatapb.OperationalStatus_OPERATIONAL_STATUS_INACTIVE
	case models.OperationalStatusDecayed:
		return metadatapb.OperationalStatus_OPERATIONAL_STATUS_DECAYED
	default:
		return metadatapb.OperationalStatus_OPERATIONAL_STATUS_UNSPECIFIED
	}
}

func mapFrequencies(freqs []models.Frequency) []*metadatapb.Frequency {
	out := make([]*metadatapb.Frequency, 0, len(freqs))

	for _, f := range freqs {
		pb := &metadatapb.Frequency{
			FrequencyMhz: f.FrequencyMHz,
			Modulation:   f.Modulation,
			Mode:         f.Mode,
		}

		if f.BandwidthKHz != nil {
			pb.BandwidthKhz = *f.BandwidthKHz
		}

		switch f.Direction {
		case models.FrequencyDirectionUplink:
			pb.Direction = metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_UPLINK
		case models.FrequencyDirectionDownlink:
			pb.Direction = metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_DOWNLINK
		default:
			pb.Direction = metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_UNSPECIFIED
		}

		out = append(out, pb)
	}

	return out
}

func mapSources(src []models.FieldSource) []*metadatapb.FieldSource {
	out := make([]*metadatapb.FieldSource, 0, len(src))

	for _, s := range src {
		pb := &metadatapb.FieldSource{
			Field:     s.Field,
			FetchedAt: timestamppb.New(s.FetchedAt),
		}

		pb.Sources = make([]metadatapb.Source, 0, len(s.Sources))

		for _, source := range s.Sources {
			pb.Sources = append(pb.Sources, mapSource(source))
		}

		out = append(out, pb)
	}

	return out
}

func mapSource(source models.Source) metadatapb.Source {
	switch source {
	case models.SourceCelestrak:
		return metadatapb.Source_SOURCE_CELESTRAK

	case models.SourceAMSAT:
		return metadatapb.Source_SOURCE_AMSAT

	case models.SourceUCS:
		return metadatapb.Source_SOURCE_UCS

	case models.SourceManual:
		return metadatapb.Source_SOURCE_MANUAL

	default:
		return metadatapb.Source_SOURCE_UNSPECIFIED
	}
}
