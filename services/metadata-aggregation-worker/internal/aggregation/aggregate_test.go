package aggregation

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRecord(
	t *testing.T,
	source models.Source,
	fetchedAt time.Time,
	payload any,
) models.SatelliteIngestRecord {
	t.Helper()

	data, err := json.Marshal(payload)
	require.NoError(t, err)

	return models.SatelliteIngestRecord{
		Source:    source,
		FetchedAt: fetchedAt,
		Payload:   data,
	}
}

func fieldSourceFor(
	t *testing.T,
	result *models.SatelliteMetadata,
	field string,
) models.FieldSource {
	t.Helper()

	var matches []models.FieldSource
	for _, source := range result.Sources {
		if source.Field == field {
			matches = append(matches, source)
		}
	}

	require.Len(t, matches, 1, "expected exactly one source for field %q", field)

	return matches[0]
}

func assertFieldSource(
	t *testing.T,
	result *models.SatelliteMetadata,
	field string,
	source models.Source,
	fetchedAt time.Time,
) {
	t.Helper()

	fs := fieldSourceFor(t, result, field)

	assert.Equal(t, []models.Source{source}, fs.Sources)
	assert.Equal(t, fetchedAt, fs.FetchedAt)
}

func TestAggregate_EmptyRecords(t *testing.T) {
	result := Aggregate(nil)
	assert.Nil(t, result)
}

func TestAggregate_EmptySlice(t *testing.T) {
	result := Aggregate([]models.SatelliteIngestRecord{})
	assert.Nil(t, result)
}

func TestAggregate_AllPayloadsInvalid(t *testing.T) {
	records := []models.SatelliteIngestRecord{
		{
			Source:    models.SourceUCS,
			FetchedAt: time.Now(),
			Payload:   []byte(`not valid json`),
		},
		{
			Source:    models.SourceCelestrak,
			FetchedAt: time.Now(),
			Payload:   []byte(`{invalid`),
		},
	}

	result := Aggregate(records)
	assert.Nil(t, result)
}

func TestAggregate_SkipsInvalidPayloads(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		{
			Source:    models.SourceUCS,
			FetchedAt: now,
			Payload:   []byte(`not valid json`),
		},
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Minute),
			map[string]any{
				"norad_id":  12345,
				"cospar_id": "2024-001A",
				"name":      "Test Satellite",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, 12345, result.NoradID)
	assert.Equal(t, "Test Satellite", result.Name)

	require.NotNil(t, result.CosparID)
	assert.Equal(t, "2024-001A", *result.CosparID)
}

func TestAggregate_PreservesIDs(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  25544,
				"cospar_id": "1998-067A",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, 25544, result.NoradID)

	require.NotNil(t, result.CosparID)
	assert.Equal(t, "1998-067A", *result.CosparID)
}

func TestAggregate_Name_PrefersUCS(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "UCS Name",
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "Celestrak Name",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, "UCS Name", result.Name)
}

func TestAggregate_Name_PrefersFreshestWhenSamePriority(t *testing.T) {
	now := time.Now()

	older := makeRecord(
		t,
		models.SourceCelestrak,
		now,
		map[string]any{
			"norad_id":  1,
			"cospar_id": "test",
			"name":      "Old Name",
		},
	)

	newer := makeRecord(
		t,
		models.SourceCelestrak,
		now.Add(time.Hour),
		map[string]any{
			"norad_id":  1,
			"cospar_id": "test",
			"name":      "New Name",
		},
	)

	result := Aggregate([]models.SatelliteIngestRecord{older, newer})

	require.NotNil(t, result)

	assert.Equal(t, "New Name", result.Name)
}

func TestAggregate_Name_IgnoresEmptyValue(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "",
			},
		),
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "Valid Name",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, "Valid Name", result.Name)
}

func TestAggregate_Name_MissingFromAllSources(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Empty(t, result.Name)
}

func TestAggregate_OperationalStatus_PrefersUCS(t *testing.T) {
	now := time.Now()

	ucsValue := models.OperationalStatusActive
	celestrakValue := models.OperationalStatusInactive

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":           1,
				"cospar_id":          "test",
				"operational_status": celestrakValue,
			},
		),
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":           1,
				"cospar_id":          "test",
				"operational_status": ucsValue,
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, ucsValue, result.OperationalStatus)
}

func TestAggregate_Fields_PreferConfiguredSourcePriority(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		field      string
		ucsValue   any
		otherValue any
		assert     func(t *testing.T, result *models.SatelliteMetadata)
	}{
		{
			name:       "name",
			field:      "name",
			ucsValue:   "UCS Name",
			otherValue: "Celestrak Name",
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				assert.Equal(t, "UCS Name", result.Name)

				assertFieldSource(
					t,
					result,
					"name",
					models.SourceUCS,
					now,
				)
			},
		},
		{
			name:       "operator",
			field:      "operator",
			ucsValue:   "UCS Operator",
			otherValue: "Celestrak Operator",
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				require.NotNil(t, result.Operator)
				assert.Equal(t, "UCS Operator", *result.Operator)

				assertFieldSource(
					t,
					result,
					"operator",
					models.SourceUCS,
					now,
				)
			},
		},
		{
			name:       "orbit_regime",
			field:      "orbit_regime",
			ucsValue:   models.OrbitRegimeGEO,
			otherValue: models.OrbitRegimeLEO,
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				assert.Equal(t, models.OrbitRegimeLEO, result.OrbitRegime)

				assertFieldSource(
					t,
					result,
					"orbit_regime",
					models.SourceCelestrak,
					now.Add(time.Hour),
				)
			},
		},
		{
			name:       "object_type",
			field:      "object_type",
			ucsValue:   models.ObjectTypeDebris,
			otherValue: models.ObjectTypePayload,
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				assert.Equal(t, models.ObjectTypePayload, result.ObjectType)

				assertFieldSource(
					t,
					result,
					"object_type",
					models.SourceCelestrak,
					now.Add(time.Hour),
				)
			},
		},
		{
			name:       "operational_status",
			field:      "operational_status",
			ucsValue:   models.OperationalStatusActive,
			otherValue: models.OperationalStatusInactive,
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				assert.Equal(
					t,
					models.OperationalStatusActive,
					result.OperationalStatus,
				)

				assertFieldSource(
					t,
					result,
					"operational_status",
					models.SourceUCS,
					now,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records := []models.SatelliteIngestRecord{
				makeRecord(
					t,
					models.SourceUCS,
					now,
					map[string]any{
						"norad_id": 1,
						tt.field:   tt.ucsValue,
					},
				),
				makeRecord(
					t,
					models.SourceCelestrak,
					now.Add(time.Hour),
					map[string]any{
						"norad_id": 1,
						tt.field:   tt.otherValue,
					},
				),
			}

			result := Aggregate(records)

			require.NotNil(t, result)

			tt.assert(t, result)
		})
	}
}

func TestAggregate_Fields_FallBackToOlderRecordWhenNewestMissing(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		field  string
		value  any
		assert func(t *testing.T, result *models.SatelliteMetadata)
	}{
		{
			name:  "launch_site",
			field: "launch_site",
			value: "Launch Site",
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				require.NotNil(t, result.LaunchSite)
				assert.Equal(t, "Launch Site", *result.LaunchSite)
			},
		},
		{
			name:  "launch_vehicle",
			field: "launch_vehicle",
			value: "Falcon 9",
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				require.NotNil(t, result.LaunchVehicle)
				assert.Equal(t, "Falcon 9", *result.LaunchVehicle)
			},
		},
		{
			name:  "owner",
			field: "owner",
			value: "Owner",
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				require.NotNil(t, result.Owner)
				assert.Equal(t, "Owner", *result.Owner)
			},
		},
		{
			name:  "constellation",
			field: "constellation",
			value: "Starlink",
			assert: func(t *testing.T, result *models.SatelliteMetadata) {
				require.NotNil(t, result.Constellation)
				assert.Equal(t, "Starlink", *result.Constellation)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records := []models.SatelliteIngestRecord{
				makeRecord(
					t,
					models.SourceUCS,
					now,
					map[string]any{
						"norad_id": 1,
						tt.field:   tt.value,
					},
				),
				makeRecord(
					t,
					models.SourceCelestrak,
					now.Add(time.Hour),
					map[string]any{
						"norad_id": 1,
					},
				),
			}

			result := Aggregate(records)

			require.NotNil(t, result)

			tt.assert(t, result)

			assertFieldSource(
				t,
				result,
				tt.field,
				models.SourceUCS,
				now,
			)
		})
	}
}

func TestPickBest_PriorityBeatsFreshness(t *testing.T) {
	now := time.Now()

	items := []item{
		{
			meta: &models.SatelliteMetadataPartial{
				Name: stringPtr("UCS Old"),
			},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceUCS,
				FetchedAt: now,
			},
		},
		{
			meta: &models.SatelliteMetadataPartial{
				Name: stringPtr("Celestrak New"),
			},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceCelestrak,
				FetchedAt: now.Add(24 * time.Hour),
			},
		},
	}

	got, bestItem, ok := pickBest(
		items,
		"name",
		func(m *models.SatelliteMetadataPartial) (string, bool) {
			if m.Name == nil {
				return "", false
			}

			return *m.Name, true
		},
	)

	require.True(t, ok)
	assert.Equal(t, "UCS Old", got)

	require.NotNil(t, bestItem)
	assert.Equal(t, models.SourceUCS, bestItem.rec.Source)
	assert.Equal(t, now, bestItem.rec.FetchedAt)
}

func TestAggregate_OptionalFieldsUseFirstNonNilByFreshness(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":    1,
				"cospar_id":   "test",
				"launch_site": "older-site",
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	require.NotNil(t, result.LaunchSite)
	assert.Equal(t, "older-site", *result.LaunchSite)
}

func TestAggregate_LaunchDate(t *testing.T) {
	now := time.Now()

	oldDate := now.Add(-24 * time.Hour)
	newDate := now.Add(-12 * time.Hour)

	tests := []struct {
		name          string
		records       []models.SatelliteIngestRecord
		expectedDate  *time.Time
		expectedSrc   models.Source
		expectedFetch time.Time
	}{
		{
			name: "uses value when present",
			records: []models.SatelliteIngestRecord{
				makeRecord(t, models.SourceUCS, now, map[string]any{
					"norad_id":    1,
					"launch_date": oldDate,
				}),
			},
			expectedDate:  &oldDate,
			expectedSrc:   models.SourceUCS,
			expectedFetch: now,
		},
		{
			name: "uses value from second record when first is missing",
			records: []models.SatelliteIngestRecord{
				makeRecord(t, models.SourceUCS, now, map[string]any{
					"norad_id": 1,
				}),
				makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{
					"norad_id":    1,
					"launch_date": oldDate,
				}),
			},
			expectedDate:  &oldDate,
			expectedSrc:   models.SourceCelestrak,
			expectedFetch: now.Add(time.Hour),
		},
		{
			name: "source priority beats freshness",
			records: []models.SatelliteIngestRecord{
				makeRecord(t, models.SourceUCS, now, map[string]any{
					"norad_id":    1,
					"launch_date": oldDate,
				}),
				makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{
					"norad_id":    1,
					"launch_date": newDate,
				}),
			},
			expectedDate:  &oldDate,
			expectedSrc:   models.SourceUCS,
			expectedFetch: now,
		},
		{
			name: "freshness wins for same source",
			records: []models.SatelliteIngestRecord{
				makeRecord(t, models.SourceUCS, now, map[string]any{
					"norad_id":    1,
					"launch_date": oldDate,
				}),
				makeRecord(t, models.SourceUCS, now.Add(time.Hour), map[string]any{
					"norad_id":    1,
					"launch_date": newDate,
				}),
			},
			expectedDate:  &newDate,
			expectedSrc:   models.SourceUCS,
			expectedFetch: now.Add(time.Hour),
		},
		{
			name: "nil when missing from all sources",
			records: []models.SatelliteIngestRecord{
				makeRecord(t, models.SourceUCS, now, map[string]any{"norad_id": 1}),
				makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{"norad_id": 1}),
			},
			expectedDate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Aggregate(tt.records)

			require.NotNil(t, result)

			if tt.expectedDate == nil {
				assert.Nil(t, result.LaunchDate)
				return
			}

			require.NotNil(t, result.LaunchDate)
			require.NotNil(t, result.LaunchDate)
			assert.True(
				t,
				tt.expectedDate.Equal(*result.LaunchDate),
				"expected %v, got %v",
				*tt.expectedDate,
				*result.LaunchDate,
			)

			assertFieldSource(
				t,
				result,
				"launch_date",
				tt.expectedSrc,
				tt.expectedFetch,
			)
		})
	}
}

func TestAggregate_AliasesAreUnique(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"aliases":   []string{"Alpha", "Beta"},
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"aliases":   []string{"Beta", "Gamma"},
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.ElementsMatch(
		t,
		[]string{"Alpha", "Beta", "Gamma"},
		result.Aliases,
	)

	assert.Len(t, result.Aliases, 3)
}

func TestAggregate_AliasesFromMultipleRecords(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"aliases":   []string{"Alpha"},
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"aliases":   []string{"Beta"},
			},
		),
		makeRecord(
			t,
			models.SourceUCS,
			now.Add(2*time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"aliases":   []string{"Alpha", "Gamma"},
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.ElementsMatch(
		t,
		[]string{"Alpha", "Beta", "Gamma"},
		result.Aliases,
	)
}

func TestAggregate_FieldSources(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "UCS Name",
				"operator":  "UCS Operator",
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":     1,
				"cospar_id":    "test",
				"name":         "Celestrak Name",
				"orbit_regime": models.OrbitRegimeLEO,
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assertFieldSource(
		t,
		result,
		"name",
		models.SourceUCS,
		now,
	)

	assertFieldSource(
		t,
		result,
		"operator",
		models.SourceUCS,
		now,
	)

	assertFieldSource(
		t,
		result,
		"orbit_regime",
		models.SourceCelestrak,
		now.Add(time.Hour),
	)
}

func TestAggregate_FieldSourcesAreUniquePerField(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id": 1,
				"name":     "Old UCS Name",
			},
		),
		makeRecord(
			t,
			models.SourceUCS,
			now.Add(time.Hour),
			map[string]any{
				"norad_id": 1,
				"name":     "New UCS Name",
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(2*time.Hour),
			map[string]any{
				"norad_id":     1,
				"name":         "Celestrak Name",
				"orbit_regime": models.OrbitRegimeLEO,
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, "New UCS Name", result.Name)

	nameSources := 0
	for _, source := range result.Sources {
		if source.Field == "name" {
			nameSources++
		}
	}

	assert.Equal(t, 1, nameSources)

	assertFieldSource(
		t,
		result,
		"name",
		models.SourceUCS,
		now.Add(time.Hour),
	)

	assertFieldSource(
		t,
		result,
		"orbit_regime",
		models.SourceCelestrak,
		now.Add(2*time.Hour),
	)
}

func TestAggregate_NoDuplicateFieldSources(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id": 1,
			"name":     "Name 1",
			"operator": "Operator 1",
		}),
		makeRecord(t, models.SourceUCS, now.Add(time.Hour), map[string]any{
			"norad_id": 1,
			"name":     "Name 2",
			"operator": "Operator 2",
		}),
		makeRecord(t, models.SourceCelestrak, now.Add(2*time.Hour), map[string]any{
			"norad_id":     1,
			"name":         "Name 3",
			"orbit_regime": models.OrbitRegimeLEO,
		}),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	counts := make(map[string]int)

	for _, source := range result.Sources {
		counts[source.Field]++
	}

	for field, count := range counts {
		assert.Equal(t, 1, count, "field %q has duplicate sources", field)
	}
}

func TestAggregate_SourcesPreserveRecordOrder(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "test",
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
				"name":      "test",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	require.Len(t, result.Sources, 2)

	assert.Equal(t, []models.Source{models.SourceUCS}, result.Sources[0].Sources)
	assert.Equal(t, []models.Source{models.SourceCelestrak}, result.Sources[1].Sources)
}

func TestAggregate_AliasesHaveMixedSources(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id": 1,
			"aliases":  []string{"Alpha", "Shared"},
		}),
		makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{
			"norad_id": 1,
			"aliases":  []string{"Beta", "Shared"},
		}),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.ElementsMatch(
		t,
		[]string{"Alpha", "Beta", "Shared"},
		result.Aliases,
	)

	fs := fieldSourceFor(t, result, "aliases")

	assert.ElementsMatch(
		t,
		[]models.Source{
			models.SourceUCS,
			models.SourceCelestrak,
		},
		fs.Sources,
	)
}

func TestAggregate_UpdatedAtIsSet(t *testing.T) {
	before := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			before,
			map[string]any{
				"norad_id":  1,
				"cospar_id": "test",
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	after := time.Now()

	assert.False(t, result.UpdatedAt.IsZero())
	assert.False(t, result.UpdatedAt.Before(before))
	assert.False(t, result.UpdatedAt.After(after))
}

func TestDecodePayload_CosparID(t *testing.T) {
	record := makeRecord(
		t,
		models.SourceUCS,
		time.Now(),
		map[string]any{
			"norad_id":  12345,
			"cospar_id": "2024-001A",
		},
	)

	m, err := decodePayload(record)

	require.NoError(t, err)
	require.NotNil(t, m)
	require.NotNil(t, m.CosparID)
	assert.Equal(t, "2024-001A", *m.CosparID)
	assert.Equal(t, 12345, m.NoradID)
}

func TestAggregate_CosparIDFromSecondRecord(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id": 12345,
			"name":     "Satellite Alpha",
		}),
		makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{
			"norad_id":  12345,
			"cospar_id": "2024-001A",
		}),
	}

	result := Aggregate(records)

	require.NotNil(t, result)
	require.NotNil(t, result.CosparID)

	assert.Equal(t, "2024-001A", *result.CosparID)
}

func TestAggregate_CosparID_PrefersFreshestWhenSamePriority(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id":  12345,
			"cospar_id": "old-cospar",
		}),
		makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{
			"norad_id":  12345,
			"cospar_id": "new-cospar",
		}),
	}

	result := Aggregate(records)

	require.NotNil(t, result)
	require.NotNil(t, result.CosparID)

	assert.Equal(t, "new-cospar", *result.CosparID)

	assertFieldSource(
		t,
		result,
		"cospar_id",
		models.SourceCelestrak,
		now.Add(time.Hour),
	)
}

func TestAggregate_CosparID_IgnoresEmptyValue(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id":  12345,
			"cospar_id": "",
		}),
		makeRecord(t, models.SourceCelestrak, now.Add(time.Hour), map[string]any{
			"norad_id":  12345,
			"cospar_id": "2024-001A",
		}),
	}

	result := Aggregate(records)

	require.NotNil(t, result)
	require.NotNil(t, result.CosparID)

	assert.Equal(t, "2024-001A", *result.CosparID)

	assertFieldSource(
		t,
		result,
		"cospar_id",
		models.SourceCelestrak,
		now.Add(time.Hour),
	)
}

func TestAggregate_MergesMultipleFields(t *testing.T) {
	now := time.Now()

	records := []models.SatelliteIngestRecord{
		makeRecord(
			t,
			models.SourceUCS,
			now,
			map[string]any{
				"norad_id":  12345,
				"cospar_id": "2024-001A",
				"name":      "Satellite Alpha",
				"operator":  "Operator Alpha",
				"aliases":   []string{"Alpha", "A-1"},
			},
		),
		makeRecord(
			t,
			models.SourceCelestrak,
			now.Add(time.Hour),
			map[string]any{
				"norad_id":  12345,
				"cospar_id": "2024-001A",
				"name":      "Satellite Beta",
				"aliases":   []string{"Beta", "A-1"},
			},
		),
	}

	result := Aggregate(records)

	require.NotNil(t, result)

	assert.Equal(t, 12345, result.NoradID)

	require.NotNil(t, result.CosparID)
	assert.Equal(t, "2024-001A", *result.CosparID)

	assert.Equal(t, "Satellite Alpha", result.Name)

	require.NotNil(t, result.Operator)
	assert.Equal(t, "Operator Alpha", *result.Operator)

	assert.ElementsMatch(
		t,
		[]string{"Alpha", "A-1", "Beta"},
		result.Aliases,
	)

	assertFieldSource(t, result, "name", models.SourceUCS, now)
	assertFieldSource(t, result, "cospar_id", models.SourceCelestrak, now.Add(time.Hour))
	assertFieldSource(t, result, "operator", models.SourceUCS, now)
}

func TestPickBest_Priority(t *testing.T) {
	now := time.Now()

	items := []item{
		{
			meta: &models.SatelliteMetadataPartial{
				Name: stringPtr("UCS"),
			},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceUCS,
				FetchedAt: now.Add(time.Hour),
			},
		},
		{
			meta: &models.SatelliteMetadataPartial{
				Name: stringPtr("Celestrak"),
			},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceCelestrak,
				FetchedAt: now,
			},
		},
	}

	got, bestItem, ok := pickBest(
		items,
		"name",
		func(m *models.SatelliteMetadataPartial) (string, bool) {
			if m.Name == nil {
				return "", false
			}

			return *m.Name, true
		},
	)

	require.True(t, ok)
	assert.Equal(t, "UCS", got)
	require.NotNil(t, bestItem)
	assert.Equal(t, models.SourceUCS, bestItem.rec.Source)
	assert.Equal(t, now.Add(time.Hour), bestItem.rec.FetchedAt)
}

func TestPickBest_SamePriorityUsesFreshness(t *testing.T) {
	now := time.Now()

	items := []item{
		{
			meta: &models.SatelliteMetadataPartial{
				Name: stringPtr("Older"),
			},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceCelestrak,
				FetchedAt: now,
			},
		},
		{
			meta: &models.SatelliteMetadataPartial{
				Name: stringPtr("Newer"),
			},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceCelestrak,
				FetchedAt: now.Add(time.Hour),
			},
		},
	}

	got, _, ok := pickBest(
		items,
		"name",
		func(m *models.SatelliteMetadataPartial) (string, bool) {
			if m.Name == nil {
				return "", false
			}

			return *m.Name, true
		},
	)

	require.True(t, ok)
	assert.Equal(t, "Newer", got)
}

func TestPickBest_NoValue(t *testing.T) {
	items := []item{
		{
			meta: &models.SatelliteMetadataPartial{},
			rec: models.SatelliteIngestRecord{
				Source:    models.SourceUCS,
				FetchedAt: time.Now(),
			},
		},
	}

	got, _, ok := pickBest(
		items,
		"name",
		func(m *models.SatelliteMetadataPartial) (string, bool) {
			if m.Name == nil {
				return "", false
			}

			return *m.Name, true
		},
	)

	assert.False(t, ok)
	assert.Empty(t, got)
}

func TestUniqueAppend(t *testing.T) {
	dst := []string{"Alpha", "Beta"}
	src := []string{"Beta", "Gamma", "Alpha", "Delta"}

	result := uniqueAppend(dst, src)

	assert.Equal(
		t,
		[]string{"Alpha", "Beta", "Gamma", "Delta"},
		result,
	)
}

func TestUniqueAppend_EmptySource(t *testing.T) {
	dst := []string{"Alpha"}

	result := uniqueAppend(dst, nil)

	assert.Equal(t, []string{"Alpha"}, result)
}

func TestUniqueAppend_EmptyDestination(t *testing.T) {
	result := uniqueAppend(nil, []string{"Alpha", "Beta", "Alpha"})

	assert.Equal(t, []string{"Alpha", "Beta"}, result)
}

func stringPtr(v string) *string {
	return &v
}
