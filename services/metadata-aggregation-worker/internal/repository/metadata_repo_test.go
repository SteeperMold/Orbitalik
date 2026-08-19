//go:build integration
// +build integration

package repository_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMeta = &models.SatelliteMetadata{
	NoradID:           25544,
	CosparID:          getPtr("1998-067A"),
	Name:              "ISS",
	Aliases:           []string{"ZARYA"},
	ObjectType:        "PAYLOAD",
	MissionType:       "CREWED",
	OrbitRegime:       "LEO",
	Operator:          getPtr("NASA"),
	Owner:             getPtr("NASA"),
	Constellation:     nil,
	LaunchDate:        getPtr(time.Now()),
	LaunchSite:        getPtr("Baikonur"),
	LaunchVehicle:     getPtr("Proton-K"),
	OperationalStatus: "OPERATIONAL",
	Frequencies: []models.Frequency{
		{
			FrequencyMHz: 145.8,
		},
	},
	Sources: []models.FieldSource{
		{
			Field:   "name",
			Sources: []models.Source{"source-1"},
		},
	},
}

func getPtr[T any](v T) *T {
	return &v
}

func clearSatelliteMetadata(t *testing.T) {
	t.Helper()

	_, err := testDB.Exec(
		"TRUNCATE satellite_metadata RESTART IDENTITY CASCADE",
	)
	require.NoError(t, err)
}

func TestMetadataRepository_Upsert_CreatesMetadata(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	err := repo.Upsert(ctx, testMeta)
	require.NoError(t, err)

	var (
		noradID           int
		cosparID          sql.NullString
		name              sql.NullString
		aliasesJSON       string
		objectType        sql.NullString
		missionType       sql.NullString
		orbitRegime       sql.NullString
		operator          sql.NullString
		owner             sql.NullString
		constellation     sql.NullString
		launchDate        sql.NullTime
		launchSite        sql.NullString
		launchVehicle     sql.NullString
		operationalStatus sql.NullString
		frequencies       []byte
		sources           []byte
		updatedAt         time.Time
	)

	const query = `
		SELECT
			norad_id,
			cospar_id,
			name,
			array_to_json(aliases)::text,
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
		FROM satellite_metadata
		WHERE norad_id = $1
	`

	err = testDB.QueryRow(query, testMeta.NoradID).Scan(
		&noradID,
		&cosparID,
		&name,
		&aliasesJSON,
		&objectType,
		&missionType,
		&orbitRegime,
		&operator,
		&owner,
		&constellation,
		&launchDate,
		&launchSite,
		&launchVehicle,
		&operationalStatus,
		&frequencies,
		&sources,
		&updatedAt,
	)

	require.NoError(t, err)

	assert.Equal(t, testMeta.NoradID, noradID)

	require.True(t, cosparID.Valid)
	assert.Equal(t, *testMeta.CosparID, cosparID.String)

	require.True(t, name.Valid)
	assert.Equal(t, testMeta.Name, name.String)

	var gotAliases []string
	require.NoError(t, json.Unmarshal([]byte(aliasesJSON), &gotAliases))
	assert.Equal(t, testMeta.Aliases, gotAliases)

	require.True(t, objectType.Valid)
	assert.Equal(t, string(testMeta.ObjectType), objectType.String)

	require.True(t, missionType.Valid)
	assert.Equal(t, string(testMeta.MissionType), missionType.String)

	require.True(t, orbitRegime.Valid)
	assert.Equal(t, string(testMeta.OrbitRegime), orbitRegime.String)

	require.True(t, operator.Valid)
	assert.Equal(t, *testMeta.Operator, operator.String)

	require.True(t, owner.Valid)
	assert.Equal(t, *testMeta.Owner, owner.String)

	assert.False(t, constellation.Valid)

	require.True(t, launchDate.Valid)
	assert.WithinDuration(
		t,
		*testMeta.LaunchDate,
		launchDate.Time,
		time.Microsecond,
	)

	require.True(t, launchSite.Valid)
	assert.Equal(t, *testMeta.LaunchSite, launchSite.String)

	require.True(t, launchVehicle.Valid)
	assert.Equal(t, *testMeta.LaunchVehicle, launchVehicle.String)

	require.True(t, operationalStatus.Valid)
	assert.Equal(
		t,
		string(testMeta.OperationalStatus),
		operationalStatus.String,
	)

	var gotFrequencies []models.Frequency
	require.NoError(t, json.Unmarshal(frequencies, &gotFrequencies))
	assert.Equal(t, testMeta.Frequencies, gotFrequencies)

	var gotSources []models.FieldSource
	require.NoError(t, json.Unmarshal(sources, &gotSources))
	assert.Equal(t, testMeta.Sources, gotSources)

	assert.False(t, updatedAt.IsZero())
}

func TestMetadataRepository_Upsert_UpdatesExistingMetadata(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	initial := &models.SatelliteMetadata{
		NoradID:           25544,
		CosparID:          getPtr("1998-067A"),
		Name:              "ISS",
		Aliases:           []string{"ZARYA"},
		ObjectType:        "PAYLOAD",
		MissionType:       "CREWED",
		OrbitRegime:       "LEO",
		Operator:          getPtr("NASA"),
		Owner:             getPtr("NASA"),
		LaunchDate:        getPtr(time.Now()),
		LaunchSite:        getPtr("Baikonur"),
		LaunchVehicle:     getPtr("Proton-K"),
		OperationalStatus: "OPERATIONAL",
		Frequencies: []models.Frequency{
			{
				FrequencyMHz: 145.8,
			},
		},
		Sources: []models.FieldSource{
			{
				Field:   "name",
				Sources: []models.Source{"source-1"},
			},
		},
	}

	require.NoError(t, repo.Upsert(ctx, initial))

	var firstUpdatedAt time.Time

	err := testDB.QueryRow(`
		SELECT updated_at
		FROM satellite_metadata
		WHERE norad_id = $1
	`, initial.NoradID).Scan(&firstUpdatedAt)

	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	updated := &models.SatelliteMetadata{
		NoradID:           25544,
		CosparID:          getPtr("1998-067A"),
		Name:              "International Space Station",
		Aliases:           []string{"ISS", "ZARYA"},
		ObjectType:        "PAYLOAD",
		MissionType:       "SCIENCE",
		OrbitRegime:       "LEO",
		Operator:          getPtr("NASA"),
		Owner:             getPtr("NASA"),
		Constellation:     getPtr("ISS"),
		LaunchDate:        getPtr(time.Now()),
		LaunchSite:        getPtr("Baikonur"),
		LaunchVehicle:     getPtr("Proton-K"),
		OperationalStatus: "ACTIVE",
		Frequencies: []models.Frequency{
			{
				FrequencyMHz: 145.8,
			},
			{
				FrequencyMHz: 437.8,
			},
		},
		Sources: []models.FieldSource{
			{
				Field:   "name",
				Sources: []models.Source{"source-2"},
			},
		},
	}

	require.NoError(t, repo.Upsert(ctx, updated))

	var (
		name              sql.NullString
		aliasesJSON       string
		missionType       sql.NullString
		constellation     sql.NullString
		operationalStatus sql.NullString
		frequencies       []byte
		sources           []byte
		afterUpdatedAt    time.Time
	)

	err = testDB.QueryRow(`
		SELECT
			name,
			array_to_json(aliases)::text,
			mission_type,
			constellation,
			operational_status,
			frequencies,
			sources,
			updated_at
		FROM satellite_metadata
		WHERE norad_id = $1
	`, updated.NoradID).Scan(
		&name,
		&aliasesJSON,
		&missionType,
		&constellation,
		&operationalStatus,
		&frequencies,
		&sources,
		&afterUpdatedAt,
	)

	require.NoError(t, err)

	require.True(t, name.Valid)
	assert.Equal(t, updated.Name, name.String)

	var gotAliases []string
	require.NoError(t, json.Unmarshal([]byte(aliasesJSON), &gotAliases))
	assert.Equal(t, updated.Aliases, gotAliases)

	require.True(t, missionType.Valid)
	assert.Equal(t, string(updated.MissionType), missionType.String)

	require.True(t, constellation.Valid)
	assert.Equal(t, *updated.Constellation, constellation.String)

	require.True(t, operationalStatus.Valid)
	assert.Equal(
		t,
		string(updated.OperationalStatus),
		operationalStatus.String,
	)

	var gotFrequencies []models.Frequency
	require.NoError(t, json.Unmarshal(frequencies, &gotFrequencies))
	assert.Equal(t, updated.Frequencies, gotFrequencies)

	var gotSources []models.FieldSource
	require.NoError(t, json.Unmarshal(sources, &gotSources))
	assert.Equal(t, updated.Sources, gotSources)

	assert.True(t, afterUpdatedAt.After(firstUpdatedAt))
}

func TestMetadataRepository_Upsert_UsesEmptyCollectionsForNilValues(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	meta := &models.SatelliteMetadata{
		NoradID:     25544,
		Name:        "ISS",
		Aliases:     nil,
		Frequencies: nil,
		Sources:     nil,
	}

	err := repo.Upsert(ctx, meta)
	require.NoError(t, err)

	var (
		aliasesJSON string
		frequencies []byte
		sources     []byte
	)

	err = testDB.QueryRow(`
		SELECT
			array_to_json(aliases)::text,
			frequencies,
			sources
		FROM satellite_metadata
		WHERE norad_id = $1
	`, meta.NoradID).Scan(
		&aliasesJSON,
		&frequencies,
		&sources,
	)

	require.NoError(t, err)

	var gotAliases []string
	require.NoError(t, json.Unmarshal([]byte(aliasesJSON), &gotAliases))

	assert.NotNil(t, gotAliases)
	assert.Empty(t, gotAliases)

	var gotFrequencies []models.Frequency
	require.NoError(t, json.Unmarshal(frequencies, &gotFrequencies))

	assert.NotNil(t, gotFrequencies)
	assert.Empty(t, gotFrequencies)

	var gotSources []models.FieldSource
	require.NoError(t, json.Unmarshal(sources, &gotSources))

	assert.NotNil(t, gotSources)
	assert.Empty(t, gotSources)

	// Upsert normalizes nil collections in the input object too.
	assert.NotNil(t, meta.Aliases)
	assert.NotNil(t, meta.Frequencies)
	assert.NotNil(t, meta.Sources)
}

func TestMetadataRepository_Upsert_CanUpdateMultipleTimes(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewMetadataRepository(testConn)

	meta := &models.SatelliteMetadata{
		NoradID: 25544,
		Name:    "first",
	}

	require.NoError(t, repo.Upsert(ctx, meta))

	meta.Name = "second"
	require.NoError(t, repo.Upsert(ctx, meta))

	meta.Name = "third"
	require.NoError(t, repo.Upsert(ctx, meta))

	var (
		count int
		name  string
	)

	err := testDB.QueryRow(`
		SELECT COUNT(*), MAX(name)
		FROM satellite_metadata
		WHERE norad_id = $1
	`, meta.NoradID).Scan(
		&count,
		&name,
	)

	require.NoError(t, err)

	assert.Equal(t, 1, count)
	assert.Equal(t, "third", name)
}

func TestMetadataRepository_Upsert_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearSatelliteMetadata(t)
	})

	repo := repository.NewMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Upsert(ctx, &models.SatelliteMetadata{
		NoradID: 25544,
		Name:    "ISS",
	})

	require.Error(t, err)
}
