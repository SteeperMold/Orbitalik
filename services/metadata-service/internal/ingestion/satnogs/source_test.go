package satnogs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSource(t *testing.T, transmitters []Transmitter, batchSize int) *Source {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/transmitters/", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(transmitters))
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.URL, time.Second, 0)

	return NewSource(client, NewMapper(), batchSize)
}

func TestSource_Name(t *testing.T) {
	source := NewSource(nil, nil, 10)

	assert.Equal(t, string(models.SourceSatNOGS), source.Name())
}

func TestSource_StreamBatch_Success(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{
			SatID:       "AO-91",
			NoradCatID:  43017,
			Description: "Amateur satellite",
			Mode:        "FM",
			Alive:       true,
			DownlinkLow: 145800000,
			UplinkLow:   435250000,
		},
	}, 10)

	var received []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			received = append(received, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, received, 1)

	assert.Equal(t, models.SourceSatNOGS, received[0].Source)
	assert.False(t, received[0].FetchedAt.IsZero())
	assert.NotEmpty(t, received[0].Raw)

	var meta models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(received[0].Raw, &meta))

	assert.Equal(t, 43017, meta.NoradID)

	require.Len(t, meta.Frequencies, 2)
	assert.Equal(t, models.FrequencyDirectionDownlink, meta.Frequencies[0].Direction)
	assert.Equal(t, 145.8, meta.Frequencies[0].FrequencyMHz)
	assert.Equal(t, models.FrequencyDirectionUplink, meta.Frequencies[1].Direction)
	assert.Equal(t, 435.25, meta.Frequencies[1].FrequencyMHz)
}

func TestSource_StreamBatch_SkipsDeadTransmitters(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{
			NoradCatID:  11111,
			Mode:        "FM",
			Alive:       false,
			DownlinkLow: 145800000,
		},
		{
			NoradCatID:  22222,
			Mode:        "FM",
			Alive:       true,
			DownlinkLow: 145900000,
		},
	}, 10)

	var received []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			received = append(received, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, received, 1)

	var meta models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(received[0].Raw, &meta))

	assert.Equal(t, 22222, meta.NoradID)
}

func TestSource_StreamBatch_BatchesRecords(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: true, Mode: "FM", DownlinkLow: 145000000},
		{NoradCatID: 2, Alive: true, Mode: "FM", DownlinkLow: 146000000},
		{NoradCatID: 3, Alive: true, Mode: "FM", DownlinkLow: 147000000},
		{NoradCatID: 4, Alive: true, Mode: "FM", DownlinkLow: 148000000},
		{NoradCatID: 5, Alive: true, Mode: "FM", DownlinkLow: 149000000},
	}, 2)

	var batchSizes []int

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			batchSizes = append(batchSizes, len(batch))
			return nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, []int{2, 2, 1}, batchSizes)
}

func TestSource_StreamBatch_ExactBatchSize(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: true, Mode: "FM"},
		{NoradCatID: 2, Alive: true, Mode: "FM"},
	}, 2)

	var calls int

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			calls++
			assert.Len(t, batch, 2)
			return nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestSource_StreamBatch_EmptyTransmitters(t *testing.T) {
	source := newTestSource(t, nil, 10)

	called := false

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			called = true
			return nil
		},
	)

	require.NoError(t, err)
	assert.False(t, called)
}

func TestSource_StreamBatch_AllTransmittersDead(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: false},
		{NoradCatID: 2, Alive: false},
	}, 10)

	called := false

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			called = true
			return nil
		},
	)

	require.NoError(t, err)
	assert.False(t, called)
}

func TestSource_StreamBatch_FlushesFinalPartialBatch(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: true, Mode: "FM"},
		{NoradCatID: 2, Alive: true, Mode: "FM"},
		{NoradCatID: 3, Alive: true, Mode: "FM"},
	}, 2)

	var batches [][]*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			batches = append(batches, batch)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, batches, 2)
	assert.Len(t, batches[0], 2)
	assert.Len(t, batches[1], 1)
}

func TestSource_StreamBatch_CallbackErrorOnFullBatch(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: true, Mode: "FM"},
		{NoradCatID: 2, Alive: true, Mode: "FM"},
		{NoradCatID: 3, Alive: true, Mode: "FM"},
	}, 2)

	expectedErr := assert.AnError
	calls := 0

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			calls++
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, calls)
}

func TestSource_StreamBatch_CallbackErrorOnFinalBatch(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: true, Mode: "FM"},
		{NoradCatID: 2, Alive: true, Mode: "FM"},
		{NoradCatID: 3, Alive: true, Mode: "FM"},
	}, 2)

	expectedErr := assert.AnError
	calls := 0

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			calls++
			if calls == 1 {
				return nil
			}
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 2, calls)
}

func TestSource_StreamBatch_FetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream failure", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	source := NewSource(
		NewClient(server.URL, time.Second, 0),
		NewMapper(),
		10,
	)

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status: 503")
}

func TestSource_StreamBatch_ContextCancellation(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{
			NoradCatID:  25544,
			Alive:       true,
			Mode:        "FM",
			DownlinkLow: 145800000,
		},
	}, 10)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := source.StreamBatch(
		ctx,
		func(batch []*ingestion.SatelliteSourceRecord) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSource_StreamBatch_FetchedAtIsShared(t *testing.T) {
	source := newTestSource(t, []Transmitter{
		{NoradCatID: 1, Alive: true, Mode: "FM"},
		{NoradCatID: 2, Alive: true, Mode: "FM"},
	}, 10)

	var records []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			records = append(records, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, records, 2)

	assert.False(t, records[0].FetchedAt.IsZero())
	assert.Equal(t, records[0].FetchedAt, records[1].FetchedAt)
}
