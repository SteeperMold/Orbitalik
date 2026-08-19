package satnogs

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapper_Map(t *testing.T) {
	mapper := NewMapper()

	tests := []struct {
		name string
		row  ingestion.Row

		wantNorad       int
		wantObjectType  models.ObjectType
		wantMissionType models.MissionType
		wantOrbitRegime models.OrbitRegime
		wantStatus      models.OperationalStatus
		wantFrequencies []models.Frequency
	}{
		{
			name: "maps downlink and uplink",
			row: ingestion.Row{
				"norad_id":     "25544",
				"mode":         "FM",
				"downlink_low": "145800000",
				"uplink_low":   "435250000",
			},
			wantNorad:       25544,
			wantObjectType:  models.ObjectTypePayload,
			wantMissionType: models.MissionTypeAmateur,
			wantOrbitRegime: models.OrbitRegimeUnspecified,
			wantStatus:      models.OperationalStatusUnknown,
			wantFrequencies: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionDownlink,
					FrequencyMHz: 145.8,
					Mode:         "FM",
					Modulation:   "FM",
				},
				{
					Direction:    models.FrequencyDirectionUplink,
					FrequencyMHz: 435.25,
					Mode:         "FM",
					Modulation:   "FM",
				},
			},
		},
		{
			name: "maps only downlink",
			row: ingestion.Row{
				"norad_id":     "12345",
				"mode":         "APT",
				"downlink_low": "137100000",
			},
			wantNorad:       12345,
			wantObjectType:  models.ObjectTypePayload,
			wantMissionType: models.MissionTypeWeather,
			wantOrbitRegime: models.OrbitRegimeUnspecified,
			wantStatus:      models.OperationalStatusUnknown,
			wantFrequencies: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionDownlink,
					FrequencyMHz: 137.1,
					Mode:         "APT",
					Modulation:   "APT",
				},
			},
		},
		{
			name: "no frequencies",
			row: ingestion.Row{
				"norad_id": "12345",
				"mode":     "UNKNOWN",
			},
			wantNorad:       12345,
			wantObjectType:  models.ObjectTypePayload,
			wantMissionType: models.MissionTypeUnspecified,
			wantOrbitRegime: models.OrbitRegimeUnspecified,
			wantStatus:      models.OperationalStatusUnknown,
			wantFrequencies: []models.Frequency{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := mapper.Map(tt.row)

			require.NoError(t, err)

			var got models.SatelliteMetadataPartial
			require.NoError(t, json.Unmarshal(data, &got))

			assert.Equal(t, tt.wantNorad, got.NoradID)
			require.NotNil(t, got.ObjectType)
			assert.Equal(t, tt.wantObjectType, *got.ObjectType)

			require.NotNil(t, got.MissionType)
			assert.Equal(t, tt.wantMissionType, *got.MissionType)

			require.NotNil(t, got.OrbitRegime)
			assert.Equal(t, tt.wantOrbitRegime, *got.OrbitRegime)

			require.NotNil(t, got.OperationalStatus)
			assert.Equal(t, tt.wantStatus, *got.OperationalStatus)

			assert.Equal(t, tt.wantFrequencies, got.Frequencies)

			assert.False(t, got.FetchedAt.IsZero())
		})
	}
}

func TestMapper_Map_InvalidNoradID(t *testing.T) {
	mapper := NewMapper()

	data, err := mapper.Map(ingestion.Row{
		"norad_id": "not-a-number",
	})

	assert.Nil(t, data)
	require.Error(t, err)
}

func TestMapper_Map_NoradIDZero(t *testing.T) {
	mapper := NewMapper()

	data, err := mapper.Map(ingestion.Row{
		"norad_id": "0",
	})

	require.NoError(t, err)

	var got models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, 0, got.NoradID)
}

func TestMapper_Map_FetchedAtIsSet(t *testing.T) {
	mapper := NewMapper()

	before := time.Now()

	data, err := mapper.Map(ingestion.Row{
		"norad_id": "25544",
	})

	require.NoError(t, err)

	after := time.Now()

	var got models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(data, &got))

	assert.False(t, got.FetchedAt.Before(before))
	assert.False(t, got.FetchedAt.After(after))
}

func TestBuildFrequencies(t *testing.T) {
	tests := []struct {
		name string
		row  ingestion.Row
		want []models.Frequency
	}{
		{
			name: "downlink and uplink",
			row: ingestion.Row{
				"downlink_low": "145800000",
				"uplink_low":   "435250000",
				"mode":         "FM",
			},
			want: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionDownlink,
					FrequencyMHz: 145.8,
					Mode:         "FM",
					Modulation:   "FM",
				},
				{
					Direction:    models.FrequencyDirectionUplink,
					FrequencyMHz: 435.25,
					Mode:         "FM",
					Modulation:   "FM",
				},
			},
		},
		{
			name: "downlink only",
			row: ingestion.Row{
				"downlink_low": "145800000",
				"mode":         "BPSK",
			},
			want: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionDownlink,
					FrequencyMHz: 145.8,
					Mode:         "BPSK",
					Modulation:   "BPSK",
				},
			},
		},
		{
			name: "uplink only",
			row: ingestion.Row{
				"uplink_low": "435250000",
				"mode":       "FM",
			},
			want: []models.Frequency{
				{
					Direction:    models.FrequencyDirectionUplink,
					FrequencyMHz: 435.25,
					Mode:         "FM",
					Modulation:   "FM",
				},
			},
		},
		{
			name: "missing frequencies",
			row:  ingestion.Row{},
			want: []models.Frequency{},
		},
		{
			name: "invalid frequencies",
			row: ingestion.Row{
				"downlink_low": "not-a-number",
				"uplink_low":   "invalid",
				"mode":         "FM",
			},
			want: []models.Frequency{},
		},
		{
			name: "zero frequencies",
			row: ingestion.Row{
				"downlink_low": "0",
				"uplink_low":   "0",
				"mode":         "FM",
			},
			want: []models.Frequency{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildFrequencies(tt.row)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseHz(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{
			name:  "145.8 MHz in Hz",
			input: "145800000",
			want:  145.8,
		},
		{
			name:  "435.25 MHz in Hz",
			input: "435250000",
			want:  435.25,
		},
		{
			name:  "decimal input",
			input: "145800000.0",
			want:  145.8,
		},
		{
			name:  "zero",
			input: "0",
			want:  0,
		},
		{
			name:  "invalid",
			input: "not-a-number",
			want:  0,
		},
		{
			name:  "empty",
			input: "",
			want:  0,
		},
		{
			name:  "negative value",
			input: "-1000000",
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseHz(tt.input))
		})
	}
}

func TestDetectMissionType(t *testing.T) {
	tests := []struct {
		mode string
		want models.MissionType
	}{
		{
			mode: "AFSK",
			want: models.MissionTypeAmateur,
		},
		{
			mode: "FM",
			want: models.MissionTypeAmateur,
		},
		{
			mode: "BPSK",
			want: models.MissionTypeAmateur,
		},
		{
			mode: "APT",
			want: models.MissionTypeWeather,
		},
		{
			mode: "LRPT",
			want: models.MissionTypeWeather,
		},
		{
			mode: "UNKNOWN",
			want: models.MissionTypeUnspecified,
		},
		{
			mode: "",
			want: models.MissionTypeUnspecified,
		},
		{
			mode: "fm",
			want: models.MissionTypeUnspecified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			assert.Equal(t, tt.want, detectMissionType(tt.mode))
		})
	}
}

func TestGetPtr(t *testing.T) {
	value := getPtr("FM")

	require.NotNil(t, value)
	assert.Equal(t, "FM", *value)
}

func TestGetPtr_Empty(t *testing.T) {
	value := getPtr("")

	assert.Nil(t, value)
}
