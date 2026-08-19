package ucs

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

func TestSource_Name(t *testing.T) {
	source := NewSource(nil, nil, nil, 10)

	assert.Equal(t, string(models.SourceUCS), source.Name())
}

func TestSource_StreamBatch_Success(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite\tOperator/Owner\tCountry of Operator/Owner\tUsers\tPurpose\tClass of Orbit\tDate of Launch\tLaunch Site\tLaunch Vehicle\tCOSPAR Number",
		"25544\tISS\tNASA\tUSA\tGovernment\tCommunications\tLEO\t11/20/1998\tBaikonur\tProton-K\t1998-067A",
	}, "\n")

	source := newTestUCSSource(t, body, 10)

	var records []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			records = append(records, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, records, 1)

	record := records[0]

	assert.Equal(t, models.SourceUCS, record.Source)
	assert.False(t, record.FetchedAt.IsZero())
	assert.NotEmpty(t, record.Raw)

	var meta models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(record.Raw, &meta))

	assert.Equal(t, 25544, meta.NoradID)

	require.NotNil(t, meta.CosparID)
	assert.Equal(t, "1998-067A", *meta.CosparID)

	require.NotNil(t, meta.Name)
	assert.Equal(t, "ISS", *meta.Name)
}

func TestSource_StreamBatch_BatchesRecords(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite\tOperator/Owner",
		"1\tSAT-1\tNASA",
		"2\tSAT-2\tNASA",
		"3\tSAT-3\tESA",
		"4\tSAT-4\tESA",
		"5\tSAT-5\tNASA",
	}, "\n")

	source := newTestUCSSource(t, body, 2)

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
		"Norad Number\tCurrent Official Name of Satellite",
		"1\tSAT-1",
		"2\tSAT-2",
	}, "\n")

	source := newTestUCSSource(t, body, 2)

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

func TestSource_StreamBatch_FinalPartialBatch(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"1\tSAT-1",
		"2\tSAT-2",
		"3\tSAT-3",
	}, "\n")

	source := newTestUCSSource(t, body, 2)

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

func TestSource_StreamBatch_EmptyInput(t *testing.T) {
	source := newTestUCSSource(t, "", 10)

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

func TestSource_StreamBatch_SkipsInvalidRows(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"\tInvalid",
		"25544\tISS",
	}, "\n")

	source := newTestUCSSource(t, body, 10)

	var records []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			records = append(records, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, records, 1)

	var meta models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(records[0].Raw, &meta))

	assert.Equal(t, 25544, meta.NoradID)
}

func TestSource_StreamBatch_SkipsMapperErrors(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"not-a-number\tBroken",
		"25544\tISS",
	}, "\n")

	source := newTestUCSSource(t, body, 10)

	var records []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			records = append(records, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, records, 1)

	var meta models.SatelliteMetadataPartial
	require.NoError(t, json.Unmarshal(records[0].Raw, &meta))

	assert.Equal(t, 25544, meta.NoradID)
}

func TestSource_StreamBatch_CallbackErrorOnFullBatch(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"1\tSAT-1",
		"2\tSAT-2",
		"3\tSAT-3",
	}, "\n")

	source := newTestUCSSource(t, body, 2)

	expectedErr := errors.New("save batch failed")
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
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"1\tSAT-1",
	}, "\n")

	source := newTestUCSSource(t, body, 10)

	expectedErr := errors.New("save final batch failed")

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
}

func TestSource_StreamBatch_FetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream failure", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	client := NewClient(server.URL, time.Second, 0)

	source := NewSource(
		client,
		NewParser(),
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
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"25544\tISS",
	}, "\n")

	source := newTestUCSSource(t, body, 10)

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
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"1\tSAT-1",
		"2\tSAT-2",
		"3\tSAT-3",
	}, "\n")

	source := newTestUCSSource(t, body, 2)

	var records []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			records = append(records, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, records, 3)

	assert.False(t, records[0].FetchedAt.IsZero())
	assert.Equal(t, records[0].FetchedAt, records[1].FetchedAt)
	assert.Equal(t, records[1].FetchedAt, records[2].FetchedAt)
}

func TestSource_StreamBatch_RecordsHaveExpectedSource(t *testing.T) {
	body := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"25544\tISS",
	}, "\n")

	source := newTestUCSSource(t, body, 10)

	var records []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			records = append(records, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, records, 1)

	assert.Equal(t, models.SourceUCS, records[0].Source)
}

func newTestUCSSource(
	t *testing.T,
	body string,
	batchSize int,
) *Source {
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
