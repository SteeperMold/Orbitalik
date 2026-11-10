//go:build integration
// +build integration

package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearSatelliteMetadata(t *testing.T) {
	t.Helper()

	_, err := testDB.Exec(
		"TRUNCATE satellite_metadata RESTART IDENTITY CASCADE",
	)

	require.NoError(t, err)
}

func insertMetadata(t *testing.T, meta *models.SatelliteMetadata) {
	t.Helper()

	frequencies, err := json.Marshal(meta.Frequencies)
	require.NoError(t, err)

	sources, err := json.Marshal(meta.Sources)
	require.NoError(t, err)

	_, err = testDB.Exec(`
		INSERT INTO satellite_metadata (
			norad_id,
			cospar_id,
			name,
			aliases,
			object_type,
			mission_type,
			orbit_regime,
			operator,
			owner,
			constellation,
			launch_date,
			launch_site,
			launch_vehicle,
			operational_status,
			frequencies,
			sources,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17
		)
	`,
		meta.NoradID,
		meta.CosparID,
		meta.Name,
		meta.Aliases,
		meta.ObjectType,
		meta.MissionType,
		meta.OrbitRegime,
		meta.Operator,
		meta.Owner,
		meta.Constellation,
		meta.LaunchDate,
		meta.LaunchSite,
		meta.LaunchVehicle,
		meta.OperationalStatus,
		frequencies,
		sources,
		meta.UpdatedAt,
	)
	require.NoError(t, err)
}

func makeMetadata(
	noradID int,
	name string,
	objectType models.ObjectType,
	missionType models.MissionType,
	orbit models.OrbitRegime,
	status models.OperationalStatus,
) *models.SatelliteMetadata {
	return &models.SatelliteMetadata{
		NoradID:           noradID,
		Name:              name,
		Aliases:           []string{"ALPHA"},
		ObjectType:        objectType,
		MissionType:       missionType,
		OrbitRegime:       orbit,
		OperationalStatus: status,
		Operator:          stringPtr("NASA"),
		Owner:             stringPtr("USA"),
		Constellation:     stringPtr("ISS"),
		LaunchDate:        timePtr(time.Date(1998, 11, 20, 0, 0, 0, 0, time.UTC)),
		LaunchSite:        stringPtr("Baikonur"),
		LaunchVehicle:     stringPtr("Proton-K"),
		Frequencies: []models.Frequency{
			{
				Direction:    models.FrequencyDirectionDownlink,
				FrequencyMHz: 145.8,
				Modulation:   "FM",
				Mode:         "Beacon",
			},
		},
		Sources: []models.FieldSource{
			{
				Field:   "name",
				Sources: []models.Source{models.SourceUCS},
			},
		},
		UpdatedAt: time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC),
	}
}

func stringPtr(v string) *string {
	return &v
}

func timePtr(v time.Time) *time.Time {
	return &v
}

func TestMetadataRepository_GetMetadataByNoradID(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	meta := makeMetadata(
		25544,
		"ISS",
		models.ObjectTypePayload,
		models.MissionTypeScience,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)

	insertMetadata(t, meta)

	got, err := repo.GetMetadataByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, meta.NoradID, got.NoradID)
	assert.Equal(t, meta.Name, got.Name)
	assert.Equal(t, meta.Aliases, got.Aliases)
	assert.Equal(t, meta.ObjectType, got.ObjectType)
	assert.Equal(t, meta.MissionType, got.MissionType)
	assert.Equal(t, meta.OrbitRegime, got.OrbitRegime)
	assert.Equal(t, meta.OperationalStatus, got.OperationalStatus)
	assert.Equal(t, meta.Frequencies, got.Frequencies)
	assert.Equal(t, meta.Sources, got.Sources)

	require.NotNil(t, got.Operator)
	assert.Equal(t, *meta.Operator, *got.Operator)

	require.NotNil(t, got.Owner)
	assert.Equal(t, *meta.Owner, *got.Owner)

	require.NotNil(t, got.Constellation)
	assert.Equal(t, *meta.Constellation, *got.Constellation)

	require.NotNil(t, got.LaunchDate)
	assert.True(t, meta.LaunchDate.Equal(*got.LaunchDate))

	require.NotNil(t, got.LaunchSite)
	assert.Equal(t, *meta.LaunchSite, *got.LaunchSite)

	require.NotNil(t, got.LaunchVehicle)
	assert.Equal(t, *meta.LaunchVehicle, *got.LaunchVehicle)

	assert.False(t, got.UpdatedAt.IsZero())
}

func TestMetadataRepository_GetMetadataByNoradID_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	got, err := repo.GetMetadataByNoradID(
		context.Background(),
		999999,
	)

	assert.Nil(t, got)
	require.Error(t, err)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestMetadataRepository_GetMetadataByNoradID_NullableFields(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	meta := &models.SatelliteMetadata{
		NoradID:           25544,
		Name:              "ISS",
		Aliases:           []string{},
		ObjectType:        models.ObjectTypePayload,
		MissionType:       models.MissionTypeUnspecified,
		OrbitRegime:       models.OrbitRegimeUnspecified,
		OperationalStatus: models.OperationalStatusUnknown,
		Frequencies:       []models.Frequency{},
		Sources:           []models.FieldSource{},
		UpdatedAt:         time.Now(),
	}

	insertMetadata(t, meta)

	got, err := repo.GetMetadataByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Nil(t, got.CosparID)
	assert.Nil(t, got.Operator)
	assert.Nil(t, got.Owner)
	assert.Nil(t, got.Constellation)
	assert.Nil(t, got.LaunchDate)
	assert.Nil(t, got.LaunchSite)
	assert.Nil(t, got.LaunchVehicle)

	assert.Empty(t, got.Aliases)
	assert.Empty(t, got.Frequencies)
	assert.Empty(t, got.Sources)
}

func TestMetadataRepository_GetMetadataByName(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			25544,
			"International Space Station",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	got, err := repo.GetMetadataByName(
		ctx,
		"International Space Station",
	)

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, 25544, got.NoradID)
	assert.Equal(t, "International Space Station", got.Name)
}

func TestMetadataRepository_GetMetadataByName_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	got, err := repo.GetMetadataByName(
		context.Background(),
		"does-not-exist",
	)

	assert.Nil(t, got)
	require.Error(t, err)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestMetadataRepository_ListSatellites_NilFilter(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	for i := 1; i <= 3; i++ {
		insertMetadata(
			t,
			makeMetadata(
				i,
				"Satellite",
				models.ObjectTypePayload,
				models.MissionTypeScience,
				models.OrbitRegimeLEO,
				models.OperationalStatusActive,
			),
		)
	}

	items, nextToken, total, err := repo.ListSatellites(ctx, nil)

	require.NoError(t, err)
	require.Len(t, items, 3)

	assert.Equal(t, uint32(3), total)
	assert.Empty(t, nextToken)

	assert.Equal(t, 1, items[0].NoradID)
	assert.Equal(t, 2, items[1].NoradID)
	assert.Equal(t, 3, items[2].NoradID)
}

func TestMetadataRepository_ListSatellites_PageSize(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	for i := 1; i <= 5; i++ {
		insertMetadata(
			t,
			makeMetadata(
				i,
				"Satellite",
				models.ObjectTypePayload,
				models.MissionTypeScience,
				models.OrbitRegimeLEO,
				models.OperationalStatusActive,
			),
		)
	}

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize: 2,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 2)

	assert.Equal(t, uint32(5), total)
	assert.Equal(t, "2", nextToken)

	assert.Equal(t, 1, items[0].NoradID)
	assert.Equal(t, 2, items[1].NoradID)
}

func TestMetadataRepository_ListSatellites_PageToken(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	for i := 1; i <= 5; i++ {
		insertMetadata(
			t,
			makeMetadata(
				i,
				"Satellite",
				models.ObjectTypePayload,
				models.MissionTypeScience,
				models.OrbitRegimeLEO,
				models.OperationalStatusActive,
			),
		)
	}

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize:  2,
			PageToken: "2",
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 2)

	assert.Equal(t, uint32(5), total)
	assert.Equal(t, "4", nextToken)

	assert.Equal(t, 3, items[0].NoradID)
	assert.Equal(t, 4, items[1].NoradID)
}

func TestMetadataRepository_ListSatellites_LastPage(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	for i := 1; i <= 5; i++ {
		insertMetadata(
			t,
			makeMetadata(
				i,
				"Satellite",
				models.ObjectTypePayload,
				models.MissionTypeScience,
				models.OrbitRegimeLEO,
				models.OperationalStatusActive,
			),
		)
	}

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize:  2,
			PageToken: "4",
		},
	)

	require.NoError(t, err)

	require.Len(t, items, 1)
	assert.Equal(t, 5, items[0].NoradID)
	assert.Equal(t, uint32(5), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_OffsetBeyondEnd(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	for i := 1; i <= 3; i++ {
		insertMetadata(
			t,
			makeMetadata(
				i,
				"Satellite",
				models.ObjectTypePayload,
				models.MissionTypeScience,
				models.OrbitRegimeLEO,
				models.OperationalStatusActive,
			),
		)
	}

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize:  10,
			PageToken: "100",
		},
	)

	require.NoError(t, err)

	assert.Empty(t, items)
	assert.Empty(t, nextToken)
	assert.Equal(t, uint32(3), total)
}

func TestMetadataRepository_ListSatellites_InvalidPageToken(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	items, nextToken, total, err := repo.ListSatellites(
		context.Background(),
		&models.ListFilter{
			PageToken: "not-a-number",
		},
	)

	assert.Nil(t, items)
	assert.Empty(t, nextToken)
	assert.Zero(t, total)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid page token")
}

func TestMetadataRepository_ListSatellites_EmptyDatabase(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	items, nextToken, total, err := repo.ListSatellites(
		context.Background(),
		&models.ListFilter{
			PageSize: 10,
		},
	)

	require.NoError(t, err)

	assert.Empty(t, items)
	assert.Empty(t, nextToken)
	assert.Equal(t, uint32(0), total)
}

func TestMetadataRepository_ListSatellites_ObjectTypeFilter(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			1,
			"Payload",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	insertMetadata(
		t,
		makeMetadata(
			2,
			"Debris",
			models.ObjectTypeDebris,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	filter := models.ObjectTypeDebris

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			ObjectType: &filter,
			PageSize:   10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, 2, items[0].NoradID)
	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_MissionTypeFilter(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			1,
			"Science",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	insertMetadata(
		t,
		makeMetadata(
			2,
			"Weather",
			models.ObjectTypePayload,
			models.MissionTypeWeather,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	filter := models.MissionTypeWeather

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			MissionType: &filter,
			PageSize:    10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, 2, items[0].NoradID)
	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_OperationalStatusFilter(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			1,
			"Active",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	insertMetadata(
		t,
		makeMetadata(
			2,
			"Inactive",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusInactive,
		),
	)

	filter := models.OperationalStatusInactive

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			OperationalStatus: &filter,
			PageSize:          10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, 2, items[0].NoradID)
	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_OrbitRegimeFilter(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			1,
			"LEO",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	insertMetadata(
		t,
		makeMetadata(
			2,
			"GEO",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeGEO,
			models.OperationalStatusActive,
		),
	)

	filter := models.OrbitRegimeGEO

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			OrbitRegime: &filter,
			PageSize:    10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, 2, items[0].NoradID)
	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_ConstellationFilter(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	iss := makeMetadata(
		1,
		"ISS",
		models.ObjectTypePayload,
		models.MissionTypeScience,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)
	iss.Constellation = stringPtr("ISS")

	starlink := makeMetadata(
		2,
		"Starlink",
		models.ObjectTypePayload,
		models.MissionTypeScience,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)
	starlink.Constellation = stringPtr("STARLINK")

	insertMetadata(t, iss)
	insertMetadata(t, starlink)

	filter := "ISS"

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			Constellation: &filter,
			PageSize:      10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, 1, items[0].NoradID)
	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_EmptyConstellationIsIgnored(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			1,
			"Satellite",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	empty := ""

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			Constellation: &empty,
			PageSize:      10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_MultipleFilters(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	match := makeMetadata(
		1,
		"Match",
		models.ObjectTypePayload,
		models.MissionTypeScience,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)
	match.Constellation = stringPtr("ISS")

	wrongObject := makeMetadata(
		2,
		"Wrong Object",
		models.ObjectTypeDebris,
		models.MissionTypeScience,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)
	wrongObject.Constellation = stringPtr("ISS")

	wrongMission := makeMetadata(
		3,
		"Wrong Mission",
		models.ObjectTypePayload,
		models.MissionTypeWeather,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)
	wrongMission.Constellation = stringPtr("ISS")

	insertMetadata(t, match)
	insertMetadata(t, wrongObject)
	insertMetadata(t, wrongMission)

	objectType := models.ObjectTypePayload
	missionType := models.MissionTypeScience
	status := models.OperationalStatusActive
	orbit := models.OrbitRegimeLEO
	constellation := "ISS"

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			ObjectType:        &objectType,
			MissionType:       &missionType,
			OperationalStatus: &status,
			OrbitRegime:       &orbit,
			Constellation:     &constellation,
			PageSize:          10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, 1, items[0].NoradID)
	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)
}

func TestMetadataRepository_ListSatellites_NoMatches(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	insertMetadata(
		t,
		makeMetadata(
			1,
			"Satellite",
			models.ObjectTypePayload,
			models.MissionTypeScience,
			models.OrbitRegimeLEO,
			models.OperationalStatusActive,
		),
	)

	filter := models.ObjectTypeDebris

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			ObjectType: &filter,
			PageSize:   10,
		},
	)

	require.NoError(t, err)

	assert.Empty(t, items)
	assert.Empty(t, nextToken)
	assert.Equal(t, uint32(0), total)
}

func TestMetadataRepository_ListSatellites_OrdersByNoradID(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	for _, id := range []int{300, 100, 200} {
		insertMetadata(
			t,
			makeMetadata(
				id,
				"Satellite",
				models.ObjectTypePayload,
				models.MissionTypeScience,
				models.OrbitRegimeLEO,
				models.OperationalStatusActive,
			),
		)
	}

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize: 10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 3)

	assert.Equal(t, uint32(3), total)
	assert.Empty(t, nextToken)

	assert.Equal(t, 100, items[0].NoradID)
	assert.Equal(t, 200, items[1].NoradID)
	assert.Equal(t, 300, items[2].NoradID)
}

func TestMetadataRepository_ListSatellites_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize: 10,
		},
	)

	assert.Nil(t, items)
	assert.Empty(t, nextToken)
	assert.Zero(t, total)
	require.Error(t, err)
}

func TestMetadataRepository_GetMetadataByNoradID_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := repo.GetMetadataByNoradID(ctx, 25544)

	assert.Nil(t, got)
	require.Error(t, err)
}

func TestMetadataRepository_GetMetadataByName_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := repo.GetMetadataByName(ctx, "ISS")

	assert.Nil(t, got)
	require.Error(t, err)
}

func TestMetadataRepository_ListSatellites_CollectionsRoundTrip(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	bandwidth := 12.5

	meta := makeMetadata(
		25544,
		"ISS",
		models.ObjectTypePayload,
		models.MissionTypeScience,
		models.OrbitRegimeLEO,
		models.OperationalStatusActive,
	)

	meta.Aliases = []string{"ZARYA", "ALPHA", "BETA"}
	meta.Frequencies = []models.Frequency{
		{
			Direction:    models.FrequencyDirectionDownlink,
			FrequencyMHz: 145.8,
			BandwidthKHz: &bandwidth,
			Modulation:   "FM",
			Mode:         "Beacon",
		},
		{
			Direction:    models.FrequencyDirectionUplink,
			FrequencyMHz: 435.25,
			Modulation:   "FM",
			Mode:         "Voice",
		},
	}
	meta.Sources = []models.FieldSource{
		{
			Field:   "name",
			Sources: []models.Source{models.SourceUCS},
		},
		{
			Field:   "aliases",
			Sources: []models.Source{models.SourceUCS, models.SourceCelestrak},
		},
	}

	insertMetadata(t, meta)

	items, nextToken, total, err := repo.ListSatellites(
		ctx,
		&models.ListFilter{
			PageSize: 10,
		},
	)

	require.NoError(t, err)
	require.Len(t, items, 1)

	assert.Equal(t, uint32(1), total)
	assert.Empty(t, nextToken)

	assert.Equal(t, meta.Aliases, items[0].Aliases)
	assert.Equal(t, meta.Frequencies, items[0].Frequencies)
	assert.Equal(t, meta.Sources, items[0].Sources)
}
