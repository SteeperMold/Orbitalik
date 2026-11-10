package ucs

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

		wantNoradID           int
		wantName              *string
		wantOperator          *string
		wantOwner             *string
		wantCosparID          *string
		wantLaunchSite        *string
		wantLaunchVehicle     *string
		wantObjectType        models.ObjectType
		wantMissionType       models.MissionType
		wantOrbitRegime       models.OrbitRegime
		wantOperationalStatus models.OperationalStatus
		wantAliases           []string
	}{
		{
			name: "maps all fields",
			row: ingestion.Row{
				"norad_id":       " 25544 ",
				"name":           " ISS ",
				"operator":       " NASA ",
				"owner":          " USA ",
				"cospar":         " 1998-067A ",
				"launch_site":    " Baikonur ",
				"launch_vehicle": " Proton-K ",
				"launch_date":    "11/20/98",
				"users":          "Commercial",
				"purpose":        "Communications",
				"orbit_class":    "LEO",
				"aliases":        "ISS (ZARYA, Alpha)",
			},
			wantNoradID:           25544,
			wantName:              getPtr("ISS"),
			wantOperator:          getPtr(" NASA "),
			wantOwner:             getPtr(" USA "),
			wantCosparID:          getPtr(" 1998-067A "),
			wantLaunchSite:        getPtr(" Baikonur "),
			wantLaunchVehicle:     getPtr(" Proton-K "),
			wantObjectType:        models.ObjectTypePayload,
			wantMissionType:       models.MissionTypeCommunications,
			wantOrbitRegime:       models.OrbitRegimeLEO,
			wantOperationalStatus: models.OperationalStatusUnknown,
			wantAliases:           []string{"ZARYA", "Alpha"},
		},
		{
			name: "missing optional fields",
			row: ingestion.Row{
				"norad_id": "25544",
			},
			wantNoradID:           25544,
			wantObjectType:        models.ObjectTypePayload,
			wantMissionType:       models.MissionTypeUnspecified,
			wantOrbitRegime:       models.OrbitRegimeUnspecified,
			wantOperationalStatus: models.OperationalStatusUnknown,
			wantAliases:           []string{},
		},
		{
			name: "maps weather mission",
			row: ingestion.Row{
				"norad_id": "12345",
				"users":    "Government",
				"purpose":  "Earth weather observation",
			},
			wantNoradID:           12345,
			wantObjectType:        models.ObjectTypePayload,
			wantMissionType:       models.MissionTypeWeather,
			wantOrbitRegime:       models.OrbitRegimeUnspecified,
			wantOperationalStatus: models.OperationalStatusUnknown,
			wantAliases:           []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := mapper.Map(tt.row)

			require.NoError(t, err)
			require.NotNil(t, raw)

			var got models.SatelliteMetadataPartial
			require.NoError(t, json.Unmarshal(raw, &got))

			assert.Equal(t, tt.wantNoradID, got.NoradID)
			assert.Equal(t, tt.wantName, got.Name)
			assert.Equal(t, tt.wantOperator, got.Operator)
			assert.Equal(t, tt.wantOwner, got.Owner)
			assert.Equal(t, tt.wantCosparID, got.CosparID)
			assert.Equal(t, tt.wantLaunchSite, got.LaunchSite)
			assert.Equal(t, tt.wantLaunchVehicle, got.LaunchVehicle)

			require.NotNil(t, got.ObjectType)
			assert.Equal(t, tt.wantObjectType, *got.ObjectType)

			require.NotNil(t, got.MissionType)
			assert.Equal(t, tt.wantMissionType, *got.MissionType)

			require.NotNil(t, got.OrbitRegime)
			assert.Equal(t, tt.wantOrbitRegime, *got.OrbitRegime)

			require.NotNil(t, got.OperationalStatus)
			assert.Equal(
				t,
				tt.wantOperationalStatus,
				*got.OperationalStatus,
			)

			assert.Equal(t, tt.wantAliases, got.Aliases)
			assert.False(t, got.FetchedAt.IsZero())

			if tt.row["launch_date"] != "" {
				require.NotNil(t, got.LaunchDate)
			} else {
				assert.Nil(t, got.LaunchDate)
			}
		})
	}
}

func TestMapper_Map_InvalidNoradID(t *testing.T) {
	mapper := NewMapper()

	raw, err := mapper.Map(ingestion.Row{
		"norad_id": "not-a-number",
	})

	assert.Nil(t, raw)
	require.Error(t, err)
}

func TestMapper_Map_FetchedAtIsSet(t *testing.T) {
	mapper := NewMapper()

	before := time.Now()

	raw, err := mapper.Map(ingestion.Row{
		"norad_id": "25544",
	})

	require.NoError(t, err)

	after := time.Now()

	var got models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(raw, &got))

	assert.False(t, got.FetchedAt.Before(before))
	assert.False(t, got.FetchedAt.After(after))
}

func TestGetPtr(t *testing.T) {
	value := getPtr("NASA")

	require.NotNil(t, value)
	assert.Equal(t, "NASA", *value)
}

func TestGetPtr_Empty(t *testing.T) {
	assert.Nil(t, getPtr(""))
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *time.Time
	}{
		{
			name:  "short date",
			input: "1/2/06",
			want:  timePtr(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:  "full date",
			input: "11/20/1998",
			want:  timePtr(time.Date(1998, 11, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:  "trimmed date",
			input: " 11/20/1998 ",
			want:  timePtr(time.Date(1998, 11, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:  "empty",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace",
			input: "   ",
			want:  nil,
		},
		{
			name:  "invalid",
			input: "1998-11-20",
			want:  nil,
		},
		{
			name:  "invalid date",
			input: "99/99/9999",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDate(tt.input)

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			assert.True(t, tt.want.Equal(*got))
		})
	}
}

func TestExtractAliases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "multiple aliases",
			input: "ISS (ZARYA, Alpha, Beta)",
			want:  []string{"ZARYA", "Alpha", "Beta"},
		},
		{
			name:  "single alias",
			input: "ISS (ZARYA)",
			want:  []string{"ZARYA"},
		},
		{
			name:  "spaces are trimmed",
			input: "ISS ( ZARYA , Alpha )",
			want:  []string{"ZARYA", "Alpha"},
		},
		{
			name:  "empty alias section",
			input: "ISS ()",
			want:  []string{},
		},
		{
			name:  "no parentheses",
			input: "ISS",
			want:  []string{},
		},
		{
			name:  "missing closing parenthesis",
			input: "ISS (ZARYA",
			want:  []string{},
		},
		{
			name:  "missing opening parenthesis",
			input: "ISS ZARYA)",
			want:  []string{},
		},
		{
			name:  "empty input",
			input: "",
			want:  []string{},
		},
		{
			name:  "empty items are ignored",
			input: "ISS (Alpha,, Beta, )",
			want:  []string{"Alpha", "Beta"},
		},
		{
			name:  "outer whitespace",
			input: "  ISS (ZARYA, Alpha)  ",
			want:  []string{"ZARYA", "Alpha"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractAliases(tt.input))
		})
	}
}

func TestMapMissionType(t *testing.T) {
	tests := []struct {
		name    string
		users   string
		purpose string
		want    models.MissionType
	}{
		{
			name:    "communications",
			users:   "Commercial",
			purpose: "Communications",
			want:    models.MissionTypeCommunications,
		},
		{
			name:    "earth observation",
			users:   "",
			purpose: "Earth observation",
			want:    models.MissionTypeEarthObservation,
		},
		{
			name:    "navigation",
			users:   "Navigation",
			purpose: "",
			want:    models.MissionTypeNavigation,
		},
		{
			name:    "science",
			users:   "",
			purpose: "Science",
			want:    models.MissionTypeScience,
		},
		{
			name:    "weather",
			users:   "",
			purpose: "Weather",
			want:    models.MissionTypeWeather,
		},
		{
			name:    "technology",
			users:   "",
			purpose: "Technology demonstration",
			want:    models.MissionTypeTechDemo,
		},
		{
			name:    "amateur",
			users:   "Amateur",
			purpose: "",
			want:    models.MissionTypeAmateur,
		},
		{
			name:    "case insensitive",
			users:   "COMMERCIAL",
			purpose: "COMMUNICATIONS",
			want:    models.MissionTypeCommunications,
		},
		{
			name:    "unknown",
			users:   "Government",
			purpose: "Military",
			want:    models.MissionTypeUnspecified,
		},
		{
			name:    "empty",
			users:   "",
			purpose: "",
			want:    models.MissionTypeUnspecified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(
				t,
				tt.want,
				mapMissionType(tt.users, tt.purpose),
			)
		})
	}
}

func TestMapOrbit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  models.OrbitRegime
	}{
		{
			name:  "LEO",
			input: "LEO",
			want:  models.OrbitRegimeLEO,
		},
		{
			name:  "MEO",
			input: "MEO",
			want:  models.OrbitRegimeMEO,
		},
		{
			name:  "GEO",
			input: "GEO",
			want:  models.OrbitRegimeGEO,
		},
		{
			name:  "HEO",
			input: "HEO",
			want:  models.OrbitRegimeHEO,
		},
		{
			name:  "lowercase",
			input: "leo",
			want:  models.OrbitRegimeLEO,
		},
		{
			name:  "spaces",
			input: "  GEO  ",
			want:  models.OrbitRegimeGEO,
		},
		{
			name:  "unknown",
			input: "UNKNOWN",
			want:  models.OrbitRegimeUnspecified,
		},
		{
			name:  "empty",
			input: "",
			want:  models.OrbitRegimeUnspecified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mapOrbit(tt.input))
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
