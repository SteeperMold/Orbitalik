//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/repository"
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

func insertRawMetadata(
	t *testing.T,
	noradID int,
	cosparID *string,
	source models.Source,
	payload []byte,
	fetchedAt time.Time,
	storedAt time.Time,
) {
	t.Helper()

	_, err := testDB.Exec(`
		INSERT INTO satellite_metadata_raw (
			norad_id,
			cospar_id,
			source,
			payload,
			fetched_at,
			stored_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		noradID,
		cosparID,
		source,
		payload,
		fetchedAt,
		storedAt,
	)

	require.NoError(t, err)
}

func TestRawMetadataRepository_GetByNoradID(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	fetchedAt := time.Now().Add(-time.Hour)
	storedAt := time.Now()

	cosparID := "1998-067A"
	payload := []byte(`{"name": "ISS"}`)

	insertRawMetadata(
		t,
		25544,
		&cosparID,
		models.SourceUCS,
		payload,
		fetchedAt,
		storedAt,
	)

	result, err := repo.GetByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.Len(t, result, 1)

	record := result[0]

	assert.Equal(t, 25544, record.NoradID)
	assert.NotZero(t, record.ID)
	assert.NotNil(t, record.CosparID)
	assert.Equal(t, cosparID, *record.CosparID)
	assert.Equal(t, models.SourceUCS, record.Source)
	assert.Equal(t, payload, []byte(record.Payload))
	assert.WithinDuration(t, fetchedAt, record.FetchedAt, time.Microsecond)
	assert.WithinDuration(t, storedAt, record.StoredAt, time.Microsecond)
}

func TestRawMetadataRepository_GetByNoradID_ReturnsRecordsInFreshnessOrder(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	now := time.Now()

	oldFetchedAt := now.Add(-2 * time.Hour)
	middleFetchedAt := now.Add(-time.Hour)
	newFetchedAt := now

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceUCS,
		[]byte(`{"name":"old"}`),
		oldFetchedAt,
		now,
	)

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceCelestrak,
		[]byte(`{"name":"middle"}`),
		middleFetchedAt,
		now,
	)

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceUCS,
		[]byte(`{"name":"new"}`),
		newFetchedAt,
		now,
	)

	result, err := repo.GetByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.Len(t, result, 3)

	assert.WithinDuration(t, newFetchedAt, result[0].FetchedAt, time.Microsecond)
	assert.WithinDuration(t, middleFetchedAt, result[1].FetchedAt, time.Microsecond)
	assert.WithinDuration(t, oldFetchedAt, result[2].FetchedAt, time.Microsecond)

	assert.Equal(t, []byte(`{"name": "new"}`), []byte(result[0].Payload))
	assert.Equal(t, []byte(`{"name": "middle"}`), []byte(result[1].Payload))
	assert.Equal(t, []byte(`{"name": "old"}`), []byte(result[2].Payload))
}

func TestRawMetadataRepository_GetByNoradID_ReturnsMultipleRecords(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	now := time.Now()

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceUCS,
		[]byte(`{"name":"ISS"}`),
		now,
		now,
	)

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceCelestrak,
		[]byte(`{"name":"International Space Station"}`),
		now.Add(-time.Hour),
		now,
	)

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceUCS,
		[]byte(`{"operator":"NASA"}`),
		now.Add(-2*time.Hour),
		now,
	)

	result, err := repo.GetByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.Len(t, result, 3)

	assert.Equal(t, models.SourceUCS, result[0].Source)
	assert.Equal(t, models.SourceCelestrak, result[1].Source)
	assert.Equal(t, models.SourceUCS, result[2].Source)
}

func TestRawMetadataRepository_GetByNoradID_IgnoresOtherNoradIDs(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	now := time.Now()

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceUCS,
		[]byte(`{"name":"ISS"}`),
		now,
		now,
	)

	insertRawMetadata(
		t,
		12345,
		nil,
		models.SourceCelestrak,
		[]byte(`{"name":"Other Satellite"}`),
		now.Add(time.Hour),
		now,
	)

	result, err := repo.GetByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.Len(t, result, 1)

	assert.Equal(t, 25544, result[0].NoradID)
	assert.Equal(t, []byte(`{"name": "ISS"}`), []byte(result[0].Payload))
}

func TestRawMetadataRepository_GetByNoradID_NullCosparID(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	now := time.Now()

	insertRawMetadata(
		t,
		25544,
		nil,
		models.SourceUCS,
		[]byte(`{"name":"ISS"}`),
		now,
		now,
	)

	result, err := repo.GetByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.Len(t, result, 1)

	assert.Nil(t, result[0].CosparID)
}

func TestRawMetadataRepository_GetByNoradID_NonNullCosparID(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	now := time.Now()
	cosparID := "1998-067A"

	insertRawMetadata(
		t,
		25544,
		&cosparID,
		models.SourceUCS,
		[]byte(`{"name":"ISS"}`),
		now,
		now,
	)

	result, err := repo.GetByNoradID(ctx, 25544)

	require.NoError(t, err)
	require.Len(t, result, 1)

	require.NotNil(t, result[0].CosparID)
	assert.Equal(t, cosparID, *result[0].CosparID)
}

func TestRawMetadataRepository_GetByNoradID_NoRecords(
	t *testing.T,
) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	ctx := context.Background()
	repo := repository.NewRawMetadataRepository(testConn)

	result, err := repo.GetByNoradID(ctx, 999999)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestRawMetadataRepository_GetByNoradID_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearRawMetadata(t)
	})

	repo := repository.NewRawMetadataRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := repo.GetByNoradID(ctx, 25544)

	assert.Nil(t, result)
	require.Error(t, err)
}
