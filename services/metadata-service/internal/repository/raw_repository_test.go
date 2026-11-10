//go:build integration
// +build integration

package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearRawMetadata(t *testing.T) {
	t.Helper()

	_, err := testDB.Exec(
		"TRUNCATE satellite_metadata_raw RESTART IDENTITY CASCADE",
	)
	require.NoError(t, err)
}

func TestRawMetadataRepository_SaveRawBatch_Empty(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	repo := repository.NewRawMetadataRepository(testConn)

	err := repo.SaveRawBatch(
		context.Background(),
		nil,
	)

	require.NoError(t, err)

	var count int

	err = testDB.QueryRow(
		"SELECT COUNT(*) FROM satellite_metadata_raw",
	).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRawMetadataRepository_SaveRawBatch_CreatesRecords(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	cosparID := "1998-067A"
	fetchedAt := time.Date(
		2026, 8, 19,
		12, 0, 0, 0,
		time.UTC,
	)

	data := []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			CosparID:  &cosparID,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":25544,"name":"ISS"}`),
			FetchedAt: fetchedAt,
		},
	}

	err := repo.SaveRawBatch(ctx, data)

	require.NoError(t, err)

	var (
		noradID    int
		cospar     *string
		source     string
		payload    []byte
		gotFetched time.Time
	)

	err = testDB.QueryRow(`
		SELECT
			norad_id,
			cospar_id,
			source,
			payload,
			fetched_at
		FROM satellite_metadata_raw
		WHERE norad_id = $1
	`, 25544).Scan(
		&noradID,
		&cospar,
		&source,
		&payload,
		&gotFetched,
	)

	require.NoError(t, err)

	assert.Equal(t, 25544, noradID)

	require.NotNil(t, cospar)
	assert.Equal(t, cosparID, *cospar)

	assert.Equal(t, string(models.SourceUCS), source)
	assert.JSONEq(
		t,
		string(data[0].Payload),
		string(payload),
	)
	assert.True(t, fetchedAt.Equal(gotFetched))
}

func TestRawMetadataRepository_SaveRawBatch_MultipleRecords(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	fetchedAt := time.Date(
		2026, 8, 19,
		12, 0, 0, 0,
		time.UTC,
	)

	data := []*models.SatelliteIngestRecord{
		{
			NoradID:   1,
			CosparID:  stringPtr("2026-001A"),
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":1}`),
			FetchedAt: fetchedAt,
		},
		{
			NoradID:   2,
			CosparID:  stringPtr("2026-002A"),
			Source:    models.SourceCelestrak,
			Payload:   json.RawMessage(`{"norad_id":2}`),
			FetchedAt: fetchedAt.Add(time.Minute),
		},
		{
			NoradID:   3,
			CosparID:  nil,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":3}`),
			FetchedAt: fetchedAt.Add(2 * time.Minute),
		},
	}

	err := repo.SaveRawBatch(ctx, data)

	require.NoError(t, err)

	var count int

	err = testDB.QueryRow(`
		SELECT COUNT(*)
		FROM satellite_metadata_raw
	`).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 3, count)

	rows, err := testDB.Query(`
		SELECT norad_id, cospar_id, source, fetched_at
		FROM satellite_metadata_raw
		ORDER BY norad_id
	`)
	require.NoError(t, err)
	defer func() {
		_ = rows.Close()
	}()

	type row struct {
		noradID   int
		cosparID  *string
		source    string
		fetchedAt time.Time
	}

	var got []row

	for rows.Next() {
		var r row

		err := rows.Scan(
			&r.noradID,
			&r.cosparID,
			&r.source,
			&r.fetchedAt,
		)
		require.NoError(t, err)

		got = append(got, r)
	}

	require.NoError(t, rows.Err())
	require.Len(t, got, 3)

	assert.Equal(t, 1, got[0].noradID)
	require.NotNil(t, got[0].cosparID)
	assert.Equal(t, "2026-001A", *got[0].cosparID)
	assert.Equal(t, string(models.SourceUCS), got[0].source)

	assert.Equal(t, 2, got[1].noradID)
	require.NotNil(t, got[1].cosparID)
	assert.Equal(t, "2026-002A", *got[1].cosparID)
	assert.Equal(t, string(models.SourceCelestrak), got[1].source)

	assert.Equal(t, 3, got[2].noradID)
	assert.Nil(t, got[2].cosparID)
	assert.Equal(t, string(models.SourceUCS), got[2].source)
}

func TestRawMetadataRepository_SaveRawBatch_AllowsNilCosparID(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	data := []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			CosparID:  nil,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":25544}`),
			FetchedAt: time.Now(),
		},
	}

	require.NoError(t, repo.SaveRawBatch(ctx, data))

	var cosparID *string

	err := testDB.QueryRow(`
		SELECT cospar_id
		FROM satellite_metadata_raw
		WHERE norad_id = $1
	`, 25544).Scan(&cosparID)

	require.NoError(t, err)
	assert.Nil(t, cosparID)
}

func TestRawMetadataRepository_SaveRawBatch_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	repo := repository.NewRawMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.SaveRawBatch(
		ctx,
		[]*models.SatelliteIngestRecord{
			{
				NoradID:   25544,
				Source:    models.SourceUCS,
				Payload:   json.RawMessage(`{"norad_id":25544}`),
				FetchedAt: time.Now(),
			},
		},
	)

	require.Error(t, err)

	var count int

	err = testDB.QueryRow(`
		SELECT COUNT(*)
		FROM satellite_metadata_raw
	`).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRawMetadataRepository_SaveRawBatch_SuccessDoesNotLeaveTransactionOpen(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	data := []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":25544}`),
			FetchedAt: time.Now(),
		},
	}

	require.NoError(t, repo.SaveRawBatch(ctx, data))

	data2 := []*models.SatelliteIngestRecord{
		{
			NoradID:   12345,
			Source:    models.SourceCelestrak,
			Payload:   json.RawMessage(`{"norad_id":12345}`),
			FetchedAt: time.Now(),
		},
	}

	require.NoError(t, repo.SaveRawBatch(ctx, data2))

	var count int

	err := testDB.QueryRow(`
		SELECT COUNT(*)
		FROM satellite_metadata_raw
	`).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestRawMetadataRepository_SaveRawBatch_ContextCancellationRollsBack(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	repo := repository.NewRawMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.SaveRawBatch(ctx, []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":25544}`),
			FetchedAt: time.Now(),
		},
	})

	require.Error(t, err)

	var count int
	err = testDB.QueryRow(`
		SELECT COUNT(*)
		FROM satellite_metadata_raw
	`).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRawMetadataRepository_SaveRawBatch_RollsBackOnCopyError(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	data := []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":25544}`),
			FetchedAt: time.Now(),
		},
		{
			NoradID:   12345,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{invalid json`),
			FetchedAt: time.Now(),
		},
	}

	err := repo.SaveRawBatch(ctx, data)

	require.Error(t, err)

	var count int

	err = testDB.QueryRow(`
		SELECT COUNT(*)
		FROM satellite_metadata_raw
	`).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRawMetadataRepository_SaveRawBatch_PreservesPayload(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	payload := json.RawMessage(`{
		"norad_id": 25544,
		"name": "ISS",
		"aliases": ["ZARYA", "ISS"],
		"nested": {
			"value": true
		}
	}`)

	err := repo.SaveRawBatch(ctx, []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			Source:    models.SourceUCS,
			Payload:   payload,
			FetchedAt: time.Now(),
		},
	})

	require.NoError(t, err)

	var got []byte

	err = testDB.QueryRow(`
		SELECT payload
		FROM satellite_metadata_raw
		WHERE norad_id = $1
	`, 25544).Scan(&got)

	require.NoError(t, err)

	assert.JSONEq(t, string(payload), string(got))
}

func TestRawMetadataRepository_SaveRawBatch_ReturnsErrorForCancelledBegin(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	repo := repository.NewRawMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.SaveRawBatch(ctx, []*models.SatelliteIngestRecord{
		{
			NoradID:   25544,
			Source:    models.SourceUCS,
			Payload:   json.RawMessage(`{"norad_id":25544}`),
			FetchedAt: time.Now(),
		},
	})

	require.Error(t, err)
	assert.NotEqual(t, errors.Is(err, context.Canceled), false)
}
