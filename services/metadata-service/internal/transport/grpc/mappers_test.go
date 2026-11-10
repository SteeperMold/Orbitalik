package grpc

import (
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/gen/metadatapb"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToProtoSatelliteMetadata_Nil(t *testing.T) {
	assert.Nil(t, toProtoSatelliteMetadata(nil))
}

func TestToProtoSatelliteMetadata(t *testing.T) {
	launchDate := time.Date(
		1998, 11, 20,
		0, 0, 0, 0,
		time.UTC,
	)

	updatedAt := time.Date(
		2026, 8, 19,
		17, 33, 24, 0,
		time.UTC,
	)

	bandwidth := 12.5

	meta := &models.SatelliteMetadata{
		NoradID:       25544,
		CosparID:      stringPtr("1998-067A"),
		Name:          "ISS",
		Aliases:       []string{"ZARYA", "Alpha"},
		ObjectType:    models.ObjectTypePayload,
		MissionType:   models.MissionTypeScience,
		OrbitRegime:   models.OrbitRegimeLEO,
		Operator:      stringPtr("NASA"),
		Owner:         stringPtr("USA"),
		Constellation: stringPtr("ISS"),
		LaunchDate:    &launchDate,
		LaunchSite:    stringPtr("Baikonur"),
		LaunchVehicle: stringPtr("Proton-K"),

		OperationalStatus: models.OperationalStatusActive,

		Frequencies: []models.Frequency{
			{
				Direction:    models.FrequencyDirectionDownlink,
				FrequencyMHz: 145.8,
				BandwidthKHz: &bandwidth,
				Modulation:   "FM",
				Mode:         "Beacon",
			},
		},

		Sources: []models.FieldSource{
			{
				Field: "name",
				Sources: []models.Source{
					models.SourceUCS,
					models.SourceCelestrak,
				},
				FetchedAt: updatedAt,
			},
		},

		UpdatedAt: updatedAt,
	}

	pb := toProtoSatelliteMetadata(meta)

	require.NotNil(t, pb)

	assert.Equal(t, uint32(25544), pb.NoradId)
	assert.Equal(t, "1998-067A", *pb.CosparId)
	assert.Equal(t, "ISS", pb.Name)
	assert.Equal(t, []string{"ZARYA", "Alpha"}, pb.Aliases)

	assert.Equal(
		t,
		metadatapb.ObjectType_OBJECT_TYPE_PAYLOAD,
		pb.ObjectType,
	)

	assert.Equal(
		t,
		metadatapb.MissionType_MISSION_TYPE_SCIENCE,
		pb.MissionType,
	)

	assert.Equal(
		t,
		metadatapb.OrbitRegime_ORBIT_REGIME_LEO,
		pb.OrbitRegime,
	)

	assert.Equal(t, "NASA", *pb.Operator)
	assert.Equal(t, "USA", *pb.Owner)
	assert.Equal(t, "ISS", *pb.Constellation)

	require.NotNil(t, pb.LaunchDate)
	assert.True(
		t,
		launchDate.Equal(pb.LaunchDate.AsTime()),
	)

	assert.Equal(t, "Baikonur", *pb.LaunchSite)
	assert.Equal(t, "Proton-K", *pb.LaunchVehicle)

	assert.Equal(
		t,
		metadatapb.OperationalStatus_OPERATIONAL_STATUS_ACTIVE,
		pb.OperationalStatus,
	)

	require.Len(t, pb.Frequencies, 1)

	freq := pb.Frequencies[0]

	assert.Equal(
		t,
		metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_DOWNLINK,
		freq.Direction,
	)
	assert.Equal(t, 145.8, freq.FrequencyMhz)
	assert.Equal(t, 12.5, freq.BandwidthKhz)
	assert.Equal(t, "FM", freq.Modulation)
	assert.Equal(t, "Beacon", freq.Mode)

	require.Len(t, pb.Sources, 1)

	source := pb.Sources[0]

	assert.Equal(t, "name", source.Field)
	assert.Equal(
		t,
		[]metadatapb.Source{
			metadatapb.Source_SOURCE_UCS,
			metadatapb.Source_SOURCE_CELESTRAK,
		},
		source.Sources,
	)

	require.NotNil(t, source.FetchedAt)
	assert.True(
		t,
		updatedAt.Equal(source.FetchedAt.AsTime()),
	)

	require.NotNil(t, pb.UpdatedAt)
	assert.True(
		t,
		updatedAt.Equal(pb.UpdatedAt.AsTime()),
	)
}

func TestToProtoSatelliteMetadata_OptionalFields(t *testing.T) {
	meta := &models.SatelliteMetadata{
		NoradID:           1,
		Name:              "Satellite",
		Aliases:           []string{},
		ObjectType:        models.ObjectTypeUnspecified,
		MissionType:       models.MissionTypeUnspecified,
		OrbitRegime:       models.OrbitRegimeUnspecified,
		OperationalStatus: models.OperationalStatusUnspecified,
		Frequencies:       []models.Frequency{},
		Sources:           []models.FieldSource{},
	}

	pb := toProtoSatelliteMetadata(meta)

	require.NotNil(t, pb)

	assert.Nil(t, pb.CosparId)
	assert.Nil(t, pb.Operator)
	assert.Nil(t, pb.Owner)
	assert.Nil(t, pb.Constellation)
	assert.Nil(t, pb.LaunchDate)
	assert.Nil(t, pb.LaunchSite)
	assert.Nil(t, pb.LaunchVehicle)
	assert.Nil(t, pb.UpdatedAt)

	assert.Empty(t, pb.Aliases)
	assert.Empty(t, pb.Frequencies)
	assert.Empty(t, pb.Sources)
}

func TestToProtoSatelliteMetadata_ZeroUpdatedAtIsOmitted(t *testing.T) {
	meta := &models.SatelliteMetadata{
		NoradID: 1,
		Name:    "Satellite",
	}

	pb := toProtoSatelliteMetadata(meta)

	require.NotNil(t, pb)
	assert.Nil(t, pb.UpdatedAt)
}

func TestToProtoSatelliteMetadata_LaunchDate(t *testing.T) {
	launchDate := time.Date(
		2020, 1, 2,
		3, 4, 5, 0,
		time.FixedZone("UTC+4", 4*60*60),
	)

	meta := &models.SatelliteMetadata{
		NoradID:    1,
		LaunchDate: &launchDate,
	}

	pb := toProtoSatelliteMetadata(meta)

	require.NotNil(t, pb)
	require.NotNil(t, pb.LaunchDate)

	assert.True(
		t,
		launchDate.Equal(pb.LaunchDate.AsTime()),
	)

	assert.True(
		t,
		timestamppb.New(launchDate).AsTime().Equal(pb.LaunchDate.AsTime()),
	)
}

func TestMapObjectType(t *testing.T) {
	tests := []struct {
		name string
		in   models.ObjectType
		want metadatapb.ObjectType
	}{
		{
			name: "payload",
			in:   models.ObjectTypePayload,
			want: metadatapb.ObjectType_OBJECT_TYPE_PAYLOAD,
		},
		{
			name: "rocket body",
			in:   models.ObjectTypeRocketBody,
			want: metadatapb.ObjectType_OBJECT_TYPE_ROCKET_BODY,
		},
		{
			name: "debris",
			in:   models.ObjectTypeDebris,
			want: metadatapb.ObjectType_OBJECT_TYPE_DEBRIS,
		},
		{
			name: "unspecified",
			in:   models.ObjectTypeUnspecified,
			want: metadatapb.ObjectType_OBJECT_TYPE_UNSPECIFIED,
		},
		{
			name: "unknown value",
			in:   models.ObjectType("unknown"),
			want: metadatapb.ObjectType_OBJECT_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapObjectType(tt.in))
		})
	}
}

func TestMapMissionType(t *testing.T) {
	tests := []struct {
		name string
		in   models.MissionType
		want metadatapb.MissionType
	}{
		{
			name: "communications",
			in:   models.MissionTypeCommunications,
			want: metadatapb.MissionType_MISSION_TYPE_COMMUNICATIONS,
		},
		{
			name: "earth observation",
			in:   models.MissionTypeEarthObservation,
			want: metadatapb.MissionType_MISSION_TYPE_EARTH_OBSERVATION,
		},
		{
			name: "navigation",
			in:   models.MissionTypeNavigation,
			want: metadatapb.MissionType_MISSION_TYPE_NAVIGATION,
		},
		{
			name: "science",
			in:   models.MissionTypeScience,
			want: metadatapb.MissionType_MISSION_TYPE_SCIENCE,
		},
		{
			name: "weather",
			in:   models.MissionTypeWeather,
			want: metadatapb.MissionType_MISSION_TYPE_WEATHER,
		},
		{
			name: "amateur",
			in:   models.MissionTypeAmateur,
			want: metadatapb.MissionType_MISSION_TYPE_AMATEUR,
		},
		{
			name: "tech demo",
			in:   models.MissionTypeTechDemo,
			want: metadatapb.MissionType_MISSION_TYPE_TECH_DEMO,
		},
		{
			name: "unspecified",
			in:   models.MissionTypeUnspecified,
			want: metadatapb.MissionType_MISSION_TYPE_UNSPECIFIED,
		},
		{
			name: "unknown value",
			in:   models.MissionType("unknown"),
			want: metadatapb.MissionType_MISSION_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapMissionType(tt.in))
		})
	}
}

func TestMapOrbitRegime(t *testing.T) {
	tests := []struct {
		name string
		in   models.OrbitRegime
		want metadatapb.OrbitRegime
	}{
		{
			name: "LEO",
			in:   models.OrbitRegimeLEO,
			want: metadatapb.OrbitRegime_ORBIT_REGIME_LEO,
		},
		{
			name: "MEO",
			in:   models.OrbitRegimeMEO,
			want: metadatapb.OrbitRegime_ORBIT_REGIME_MEO,
		},
		{
			name: "GEO",
			in:   models.OrbitRegimeGEO,
			want: metadatapb.OrbitRegime_ORBIT_REGIME_GEO,
		},
		{
			name: "HEO",
			in:   models.OrbitRegimeHEO,
			want: metadatapb.OrbitRegime_ORBIT_REGIME_HEO,
		},
		{
			name: "unspecified",
			in:   models.OrbitRegimeUnspecified,
			want: metadatapb.OrbitRegime_ORBIT_REGIME_UNSPECIFIED,
		},
		{
			name: "unknown value",
			in:   models.OrbitRegime("unknown"),
			want: metadatapb.OrbitRegime_ORBIT_REGIME_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapOrbitRegime(tt.in))
		})
	}
}

func TestMapOperationalStatus(t *testing.T) {
	tests := []struct {
		name string
		in   models.OperationalStatus
		want metadatapb.OperationalStatus
	}{
		{
			name: "active",
			in:   models.OperationalStatusActive,
			want: metadatapb.OperationalStatus_OPERATIONAL_STATUS_ACTIVE,
		},
		{
			name: "inactive",
			in:   models.OperationalStatusInactive,
			want: metadatapb.OperationalStatus_OPERATIONAL_STATUS_INACTIVE,
		},
		{
			name: "decayed",
			in:   models.OperationalStatusDecayed,
			want: metadatapb.OperationalStatus_OPERATIONAL_STATUS_DECAYED,
		},
		{
			name: "unspecified",
			in:   models.OperationalStatusUnspecified,
			want: metadatapb.OperationalStatus_OPERATIONAL_STATUS_UNSPECIFIED,
		},
		{
			name: "unknown value",
			in:   models.OperationalStatus("unknown"),
			want: metadatapb.OperationalStatus_OPERATIONAL_STATUS_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapOperationalStatus(tt.in))
		})
	}
}

func TestMapFrequencies(t *testing.T) {
	bandwidth := 25.0

	tests := []struct {
		name string
		in   []models.Frequency
		want []*metadatapb.Frequency
	}{
		{
			name: "empty",
			in:   []models.Frequency{},
			want: []*metadatapb.Frequency{},
		},
		{
			name: "downlink with bandwidth",
			in: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionDownlink,
					FrequencyMHz: 145.8,
					BandwidthKHz: &bandwidth,
					Modulation:   "FM",
					Mode:         "Beacon",
				},
			},
			want: []*metadatapb.Frequency{
				{
					Direction:    metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_DOWNLINK,
					FrequencyMhz: 145.8,
					BandwidthKhz: 25.0,
					Modulation:   "FM",
					Mode:         "Beacon",
				},
			},
		},
		{
			name: "uplink without bandwidth",
			in: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionUplink,
					FrequencyMHz: 435.25,
					Modulation:   "FM",
					Mode:         "Voice",
				},
			},
			want: []*metadatapb.Frequency{
				{
					Direction:    metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_UPLINK,
					FrequencyMhz: 435.25,
					BandwidthKhz: 0,
					Modulation:   "FM",
					Mode:         "Voice",
				},
			},
		},
		{
			name: "unspecified direction",
			in: []models.Frequency{
				{
					FrequencyMHz: 10,
				},
			},
			want: []*metadatapb.Frequency{
				{
					Direction:    metadatapb.FrequencyDirection_FREQUENCY_DIRECTION_UNSPECIFIED,
					FrequencyMhz: 10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapFrequencies(tt.in))
		})
	}
}

func TestMapSources(t *testing.T) {
	fetchedAt := time.Date(
		2026, 8, 19,
		10, 0, 0, 0,
		time.UTC,
	)

	got := mapSources([]models.FieldSource{
		{
			Field: "name",
			Sources: []models.Source{
				models.SourceCelestrak,
				models.SourceAMSAT,
				models.SourceUCS,
				models.SourceManual,
			},
			FetchedAt: fetchedAt,
		},
	})

	require.Len(t, got, 1)

	assert.Equal(t, "name", got[0].Field)

	assert.Equal(
		t,
		[]metadatapb.Source{
			metadatapb.Source_SOURCE_CELESTRAK,
			metadatapb.Source_SOURCE_AMSAT,
			metadatapb.Source_SOURCE_UCS,
			metadatapb.Source_SOURCE_MANUAL,
		},
		got[0].Sources,
	)

	require.NotNil(t, got[0].FetchedAt)
	assert.True(t, fetchedAt.Equal(got[0].FetchedAt.AsTime()))
}

func TestMapSources_Empty(t *testing.T) {
	got := mapSources([]models.FieldSource{})

	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestMapSource(t *testing.T) {
	tests := []struct {
		name string
		in   models.Source
		want metadatapb.Source
	}{
		{
			name: "Celestrak",
			in:   models.SourceCelestrak,
			want: metadatapb.Source_SOURCE_CELESTRAK,
		},
		{
			name: "AMSAT",
			in:   models.SourceAMSAT,
			want: metadatapb.Source_SOURCE_AMSAT,
		},
		{
			name: "UCS",
			in:   models.SourceUCS,
			want: metadatapb.Source_SOURCE_UCS,
		},
		{
			name: "manual",
			in:   models.SourceManual,
			want: metadatapb.Source_SOURCE_MANUAL,
		},
		{
			name: "unknown",
			in:   models.Source("unknown"),
			want: metadatapb.Source_SOURCE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapSource(tt.in))
		})
	}
}

func stringPtr(v string) *string {
	return &v
}
