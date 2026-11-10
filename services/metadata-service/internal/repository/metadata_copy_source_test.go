package repository

import (
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawCopySource_Next(t *testing.T) {
	data := []*models.SatelliteIngestRecord{
		{
			NoradID: 25544,
			Payload: []byte(`{"name":"ISS"}`),
		},
		{
			NoradID: 12345,
			Payload: []byte(`{"name":"Satellite"}`),
		},
	}

	source := newRawMetadataCopySource(data)
	require.NotNil(t, source)

	assert.True(t, source.Next())
	_, err := source.Values()
	require.NoError(t, err)

	assert.True(t, source.Next())
	_, err = source.Values()
	require.NoError(t, err)

	assert.False(t, source.Next())
}

func TestRawCopySource_Values(t *testing.T) {
	fetchedAt := time.Date(
		2026, 8, 18,
		10, 0, 0, 0,
		time.UTC,
	)

	cosparID := "1998-067A"

	record := &models.SatelliteIngestRecord{
		NoradID:   25544,
		CosparID:  &cosparID,
		Source:    models.SourceUCS,
		Payload:   []byte(`{"name":"ISS"}`),
		FetchedAt: fetchedAt,
	}

	source := newRawMetadataCopySource([]*models.SatelliteIngestRecord{
		record,
	})

	require.True(t, source.Next())

	values, err := source.Values()

	require.NoError(t, err)
	require.Len(t, values, 5)

	assert.Equal(t, 25544, values[0])
	assert.Equal(t, &cosparID, values[1])
	assert.Equal(t, models.SourceUCS, values[2])
	assert.Equal(t, []byte(`{"name":"ISS"}`), values[3])
	assert.Equal(t, fetchedAt, values[4])

	assert.False(t, source.Next())
}

func TestRawCopySource_Values_NilCosparID(t *testing.T) {
	record := &models.SatelliteIngestRecord{
		NoradID:  25544,
		CosparID: nil,
		Source:   models.SourceCelestrak,
		Payload:  []byte(`{"name":"ISS"}`),
	}

	source := newRawMetadataCopySource([]*models.SatelliteIngestRecord{
		record,
	})

	require.True(t, source.Next())

	values, err := source.Values()

	require.NoError(t, err)
	require.Len(t, values, 5)

	assert.Nil(t, values[1])
	assert.Equal(t, models.SourceCelestrak, values[2])
}

func TestRawCopySource_Values_AdvancesIndex(t *testing.T) {
	first := &models.SatelliteIngestRecord{
		NoradID: 1,
		Payload: []byte(`{"name":"one"}`),
	}

	second := &models.SatelliteIngestRecord{
		NoradID: 2,
		Payload: []byte(`{"name":"two"}`),
	}

	source := newRawMetadataCopySource([]*models.SatelliteIngestRecord{
		first,
		second,
	})

	require.True(t, source.Next())

	values, err := source.Values()
	require.NoError(t, err)
	assert.Equal(t, 1, values[0])

	require.True(t, source.Next())

	values, err = source.Values()
	require.NoError(t, err)
	assert.Equal(t, 2, values[0])

	assert.False(t, source.Next())
}

func TestRawCopySource_Empty(t *testing.T) {
	source := newRawMetadataCopySource(nil)

	assert.False(t, source.Next())
	assert.NoError(t, source.Err())
}

func TestRawCopySource_Err(t *testing.T) {
	source := newRawMetadataCopySource(nil)

	assert.NoError(t, source.Err())
}
