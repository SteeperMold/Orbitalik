package celestrak

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
		want models.SatelliteMetadataPartial
	}{
		{
			name: "maps all fields",
			row: ingestion.Row{
				"norad_id":    "25544",
				"cospar_id":   "1998-067A",
				"name":        "ISS",
				"owner":       "NASA",
				"launch_date": "1998-11-20",
				"launch_site": "Baikonur",
				"flags":       "",
				"decay_date":  "",
				"apogee":      "420",
			},
			want: models.SatelliteMetadataPartial{
				NoradID:    25544,
				CosparID:   getPtr("1998-067A"),
				Name:       getPtr("ISS"),
				Owner:      getPtr("NASA"),
				LaunchSite: getPtr("Baikonur"),
				ObjectType: func() *models.ObjectType {
					v := models.ObjectTypePayload
					return &v
				}(),
				OperationalStatus: func() *models.OperationalStatus {
					v := models.OperationalStatusUnknown
					return &v
				}(),
				OrbitRegime: func() *models.OrbitRegime {
					v := models.OrbitRegimeLEO
					return &v
				}(),
			},
		},
		{
			name: "trims string fields",
			row: ingestion.Row{
				"norad_id":    " 25544 ",
				"cospar_id":   " 1998-067A ",
				"name":        " ISS ",
				"owner":       " NASA ",
				"launch_site": " Baikonur ",
				"apogee":      "1000",
			},
			want: models.SatelliteMetadataPartial{
				NoradID:    25544,
				CosparID:   getPtr("1998-067A"),
				Name:       getPtr("ISS"),
				Owner:      getPtr("NASA"),
				LaunchSite: getPtr("Baikonur"),
				ObjectType: func() *models.ObjectType {
					v := models.ObjectTypePayload
					return &v
				}(),
				OperationalStatus: func() *models.OperationalStatus {
					v := models.OperationalStatusUnknown
					return &v
				}(),
				OrbitRegime: func() *models.OrbitRegime {
					v := models.OrbitRegimeLEO
					return &v
				}(),
			},
		},
		{
			name: "empty optional fields remain nil",
			row: ingestion.Row{
				"norad_id": "25544",
			},
			want: models.SatelliteMetadataPartial{
				NoradID:           25544,
				ObjectType:        getPtr(models.ObjectTypePayload),
				OperationalStatus: getPtr(models.OperationalStatusUnknown),
			},
		},
		{
			name: "decayed satellite from flag",
			row: ingestion.Row{
				"norad_id": "25544",
				"name":     "TEST DEB",
				"flags":    "D",
			},
			want: models.SatelliteMetadataPartial{
				NoradID: 25544,
				Name:    getPtr("TEST DEB"),
				ObjectType: func() *models.ObjectType {
					v := models.ObjectTypeDebris
					return &v
				}(),
				OperationalStatus: func() *models.OperationalStatus {
					v := models.OperationalStatusDecayed
					return &v
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := mapper.Map(tt.row)

			require.NoError(t, err)

			var got models.SatelliteMetadataPartial
			require.NoError(t, json.Unmarshal(data, &got))

			assert.Equal(t, tt.want.NoradID, got.NoradID)
			assert.Equal(t, tt.want.CosparID, got.CosparID)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Owner, got.Owner)
			assert.Equal(t, tt.want.LaunchSite, got.LaunchSite)
			assert.Equal(t, tt.want.ObjectType, got.ObjectType)
			assert.Equal(t, tt.want.OperationalStatus, got.OperationalStatus)
			assert.Equal(t, tt.want.OrbitRegime, got.OrbitRegime)
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

func TestMapper_Map_SetsFetchedAt(t *testing.T) {
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

func TestParseDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *time.Time
	}{
		{
			name:  "valid date",
			input: "1998-11-20",
			want:  getPtr(time.Date(1998, 11, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:  "leading and trailing spaces",
			input: " 1998-11-20 ",
			want:  getPtr(time.Date(1998, 11, 20, 0, 0, 0, 0, time.UTC)),
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
			name:  "invalid format",
			input: "20/11/1998",
			want:  nil,
		},
		{
			name:  "invalid date",
			input: "1998-99-99",
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

func TestDeriveStatus(t *testing.T) {
	tests := []struct {
		name   string
		flags  string
		decay  string
		expect models.OperationalStatus
	}{
		{
			name:   "normal",
			expect: models.OperationalStatusUnknown,
		},
		{
			name:   "decayed flag",
			flags:  "D",
			expect: models.OperationalStatusDecayed,
		},
		{
			name:   "decayed flag surrounded by other flags",
			flags:  "ABC D XYZ",
			expect: models.OperationalStatusDecayed,
		},
		{
			name:   "decay date",
			decay:  "2020-01-01",
			expect: models.OperationalStatusDecayed,
		},
		{
			name:   "whitespace decay date",
			decay:  "  2020-01-01  ",
			expect: models.OperationalStatusDecayed,
		},
		{
			name:   "whitespace flags",
			flags:  "   ",
			expect: models.OperationalStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveStatus(tt.flags, tt.decay)
			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestDeriveObjectType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  models.ObjectType
	}{
		{
			name:  "payload",
			input: "ISS",
			want:  models.ObjectTypePayload,
		},
		{
			name:  "rocket body",
			input: "FALCON 9 R/B",
			want:  models.ObjectTypeRocketBody,
		},
		{
			name:  "debris",
			input: "COSMOS 1234 DEB",
			want:  models.ObjectTypeDebris,
		},
		{
			name:  "lowercase rocket body",
			input: "falcon 9 r/b",
			want:  models.ObjectTypeRocketBody,
		},
		{
			name:  "lowercase debris",
			input: "cosmos 1234 deb",
			want:  models.ObjectTypeDebris,
		},
		{
			name:  "empty defaults to payload",
			input: "",
			want:  models.ObjectTypePayload,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, deriveObjectType(tt.input))
		})
	}
}

func TestDeriveOrbit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  models.OrbitRegime
	}{
		{
			name:  "empty",
			input: "",
			want:  models.OrbitRegimeUnspecified,
		},
		{
			name:  "whitespace",
			input: "   ",
			want:  models.OrbitRegimeUnspecified,
		},
		{
			name:  "invalid",
			input: "not-a-number",
			want:  models.OrbitRegimeUnspecified,
		},
		{
			name:  "LEO",
			input: "1999",
			want:  models.OrbitRegimeLEO,
		},
		{
			name:  "LEO boundary",
			input: "0",
			want:  models.OrbitRegimeLEO,
		},
		{
			name:  "MEO lower boundary",
			input: "2000",
			want:  models.OrbitRegimeMEO,
		},
		{
			name:  "MEO upper boundary",
			input: "35785",
			want:  models.OrbitRegimeMEO,
		},
		{
			name:  "GEO lower boundary",
			input: "35786",
			want:  models.OrbitRegimeGEO,
		},
		{
			name:  "GEO upper boundary",
			input: "36000",
			want:  models.OrbitRegimeGEO,
		},
		{
			name:  "HEO",
			input: "36001",
			want:  models.OrbitRegimeHEO,
		},
		{
			name:  "spaces trimmed",
			input: "  5000  ",
			want:  models.OrbitRegimeMEO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, deriveOrbit(tt.input))
		})
	}
}

func getPtr[T any](v T) *T {
	return &v
}
