package celestrak

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeCelestrakLine(noradID, name string) string {
	line := make([]byte, 117)

	for i := range line {
		line[i] = ' '
	}

	put := func(start, end int, value string) {
		value = strings.TrimSpace(value)

		if len(value) > end-start {
			value = value[:end-start]
		}

		copy(line[start:end], value)
	}

	put(13, 18, noradID)
	put(23, 47, name)
	put(49, 54, "US")
	put(56, 66, "1998-11-20")
	put(68, 73, "KSC")
	put(75, 85, "")
	put(87, 94, "92.6")
	put(96, 101, "51.6")
	put(103, 109, "420")
	put(111, 117, "410")

	return string(line)
}

func newTestSource(t *testing.T, body string, batchSize int) *Source {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.URL, time.Second, 0)

	return NewSource(
		client,
		NewParser(),
		NewMapper(),
		batchSize,
	)
}

func TestSource_Name(t *testing.T) {
	source := NewSource(
		nil,
		nil,
		nil,
		1,
	)

	assert.Equal(t, string(models.SourceCelestrak), source.Name())
}

func TestSource_StreamBatch_Success(t *testing.T) {
	body := strings.Join([]string{
		makeCelestrakLine("25544", "ISS"),
		makeCelestrakLine("12345", "SAT-A"),
	}, "\n")

	source := newTestSource(t, body, 10)

	var batches [][]*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			batches = append(batches, batch)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, batches, 1)
	require.Len(t, batches[0], 2)

	assert.Equal(t, models.SourceCelestrak, batches[0][0].Source)
	assert.Equal(t, models.SourceCelestrak, batches[0][1].Source)

	assert.NotEmpty(t, batches[0][0].Raw)
	assert.NotEmpty(t, batches[0][1].Raw)

	var first models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(batches[0][0].Raw, &first))

	var second models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(batches[0][1].Raw, &second))

	assert.Equal(t, 25544, first.NoradID)
	assert.Equal(t, 12345, second.NoradID)
}

func TestSource_StreamBatch_FlushesFullBatches(t *testing.T) {
	body := strings.Join([]string{
		makeCelestrakLine("1", "SAT-1"),
		makeCelestrakLine("2", "SAT-2"),
		makeCelestrakLine("3", "SAT-3"),
		makeCelestrakLine("4", "SAT-4"),
		makeCelestrakLine("5", "SAT-5"),
	}, "\n")

	source := newTestSource(t, body, 2)

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
	body := strings.Join([]string{
		makeCelestrakLine("1", "SAT-1"),
		makeCelestrakLine("2", "SAT-2"),
		makeCelestrakLine("3", "SAT-3"),
	}, "\n")

	source := newTestSource(t, body, 3)

	var batches [][]*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			batches = append(batches, batch)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, batches, 1)
	assert.Len(t, batches[0], 3)
}

func TestSource_StreamBatch_EmptyInput(t *testing.T) {
	source := newTestSource(t, "", 10)

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

func TestSource_StreamBatch_SkipsInvalidLines(t *testing.T) {
	body := strings.Join([]string{
		"too short",
		makeCelestrakLine("25544", "ISS"),
		"still too short",
	}, "\n")

	source := newTestSource(t, body, 10)

	var batches [][]*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			batches = append(batches, batch)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, batches, 1)
	require.Len(t, batches[0], 1)
}

func TestSource_StreamBatch_CallbackErrorOnFullBatch(t *testing.T) {
	body := strings.Join([]string{
		makeCelestrakLine("1", "SAT-1"),
		makeCelestrakLine("2", "SAT-2"),
		makeCelestrakLine("3", "SAT-3"),
	}, "\n")

	source := newTestSource(t, body, 2)

	expectedErr := errors.New("save batch failed")
	callCount := 0

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			callCount++
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 1, callCount)
}

func TestSource_StreamBatch_CallbackErrorOnFinalBatch(t *testing.T) {
	body := strings.Join([]string{
		makeCelestrakLine("1", "SAT-1"),
		makeCelestrakLine("2", "SAT-2"),
		makeCelestrakLine("3", "SAT-3"),
	}, "\n")

	source := newTestSource(t, body, 2)

	expectedErr := errors.New("save final batch failed")
	callCount := 0

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			callCount++

			if callCount == 1 {
				return nil
			}

			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
	assert.Equal(t, 2, callCount)
}

func TestSource_StreamBatch_ContextCancellation(t *testing.T) {
	body := makeCelestrakLine("25544", "ISS")

	source := newTestSource(t, body, 10)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := source.StreamBatch(
		ctx,
		func(batch []*ingestion.SatelliteSourceRecord) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	require.ErrorIs(t, err, context.Canceled)
}

func TestSource_StreamBatch_RecordFetchTimesAreSame(t *testing.T) {
	body := strings.Join([]string{
		makeCelestrakLine("1", "SAT-1"),
		makeCelestrakLine("2", "SAT-2"),
		makeCelestrakLine("3", "SAT-3"),
	}, "\n")

	source := newTestSource(t, body, 10)

	var batch []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(records []*ingestion.SatelliteSourceRecord) error {
			batch = records
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, batch, 3)

	assert.False(t, batch[0].FetchedAt.IsZero())
	assert.Equal(t, batch[0].FetchedAt, batch[1].FetchedAt)
	assert.Equal(t, batch[0].FetchedAt, batch[2].FetchedAt)
}

func TestSource_StreamBatch_RecordsHaveExpectedSource(t *testing.T) {
	body := makeCelestrakLine("25544", "ISS")

	source := newTestSource(t, body, 10)

	var batch []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(records []*ingestion.SatelliteSourceRecord) error {
			batch = records
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, batch, 1)

	assert.Equal(t, models.SourceCelestrak, batch[0].Source)
}
