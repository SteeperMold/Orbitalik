package filesource

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockParser struct {
	mock.Mock
}

func (m *mockParser) Stream(
	r io.Reader,
	fn func(ingestion.Row) error,
) error {
	args := m.Called(r, mock.Anything)

	if err := args.Error(0); err != nil {
		return err
	}

	// Let the mock decide which rows to send to the callback.
	rows := args.Get(1)
	if rows == nil {
		return nil
	}

	for _, row := range rows.([]ingestion.Row) {
		if err := fn(row); err != nil {
			return err
		}
	}

	return nil
}

type mockMapper struct {
	mock.Mock
}

func (m *mockMapper) Map(row ingestion.Row) (json.RawMessage, error) {
	args := m.Called(row)

	var raw json.RawMessage
	if args.Get(0) != nil {
		raw = args.Get(0).(json.RawMessage)
	}

	return raw, args.Error(1)
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "metadata.txt")

	require.NoError(
		t,
		os.WriteFile(path, []byte(content), 0600),
	)

	return path
}

func TestSource_Name(t *testing.T) {
	source := NewSource(
		"ucs",
		"/tmp/test",
		nil,
		nil,
		10,
	)

	assert.Equal(t, "ucs", source.Name())
}

func TestSource_StreamBatch_Success(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
	}

	raw1 := json.RawMessage(`{"norad_id":1}`)
	raw2 := json.RawMessage(`{"norad_id":2}`)

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	mapper.
		On("Map", rows[0]).
		Return(raw1, nil).
		Once()

	mapper.
		On("Map", rows[1]).
		Return(raw2, nil).
		Once()

	source := NewSource("ucs", path, parser, mapper, 10)

	var received []*ingestion.SatelliteSourceRecord

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			received = append(received, batch...)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, received, 2)

	assert.Equal(t, models.Source("ucs"), received[0].Source)
	assert.Equal(t, models.Source("ucs"), received[1].Source)

	assert.Equal(t, raw1, received[0].Raw)
	assert.Equal(t, raw2, received[1].Raw)

	assert.False(t, received[0].FetchedAt.IsZero())
	assert.Equal(t, received[0].FetchedAt, received[1].FetchedAt)

	parser.AssertExpectations(t)
	mapper.AssertExpectations(t)
}

func TestSource_StreamBatch_FlushesFullBatches(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
		{"norad_id": "3"},
		{"norad_id": "4"},
		{"norad_id": "5"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	for i, row := range rows {
		mapper.
			On("Map", row).
			Return(json.RawMessage(`{"norad_id":1}`), nil).
			Once()

		_ = i
	}

	source := NewSource("ucs", path, parser, mapper, 2)

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
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	for _, row := range rows {
		mapper.
			On("Map", row).
			Return(json.RawMessage(`{}`), nil).
			Once()
	}

	source := NewSource("ucs", path, parser, mapper, 2)

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
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	for _, row := range rows {
		mapper.
			On("Map", row).
			Return(json.RawMessage(`{}`), nil).
			Once()
	}

	source := NewSource("ucs", path, parser, mapper, 10)

	var batchSizes []int

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			batchSizes = append(batchSizes, len(batch))
			return nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, []int{2}, batchSizes)
}

func TestSource_StreamBatch_MapperErrorIsSkipped(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	mapper.
		On("Map", rows[0]).
		Return(nil, errors.New("invalid row")).
		Once()

	mapper.
		On("Map", rows[1]).
		Return(json.RawMessage(`{"norad_id":2}`), nil).
		Once()

	source := NewSource("ucs", path, parser, mapper, 10)

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
	assert.Equal(t, json.RawMessage(`{"norad_id":2}`), received[0].Raw)

	mapper.AssertExpectations(t)
}

func TestSource_StreamBatch_ParserError(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	expectedErr := errors.New("parser failed")

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(expectedErr, nil).
		Once()

	source := NewSource("ucs", path, parser, mapper, 10)

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	require.ErrorIs(t, err, expectedErr)
	mapper.AssertNotCalled(t, "Map", mock.Anything)
}

func TestSource_StreamBatch_CallbackErrorOnFullBatch(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	for _, row := range rows {
		mapper.
			On("Map", row).
			Return(json.RawMessage(`{}`), nil).
			Once()
	}

	expectedErr := errors.New("save failed")
	source := NewSource("ucs", path, parser, mapper, 2)

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
}

func TestSource_StreamBatch_CallbackErrorOnFinalBatch(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	mapper.
		On("Map", rows[0]).
		Return(json.RawMessage(`{}`), nil).
		Once()

	expectedErr := errors.New("save failed")

	source := NewSource("ucs", path, parser, mapper, 10)

	err := source.StreamBatch(
		context.Background(),
		func(batch []*ingestion.SatelliteSourceRecord) error {
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
}

func TestSource_StreamBatch_ContextCancellation(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	source := NewSource("ucs", path, parser, mapper, 10)

	err := source.StreamBatch(
		ctx,
		func(batch []*ingestion.SatelliteSourceRecord) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	require.ErrorIs(t, err, context.Canceled)

	mapper.AssertNotCalled(t, "Map", mock.Anything)
}

func TestSource_StreamBatch_FileOpenError(t *testing.T) {
	parser := new(mockParser)
	mapper := new(mockMapper)

	source := NewSource(
		"ucs",
		"/definitely/does/not/exist",
		parser,
		mapper,
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

	parser.AssertNotCalled(t, "Stream", mock.Anything, mock.Anything)
	mapper.AssertNotCalled(t, "Map", mock.Anything)
}

func TestSource_StreamBatch_EmptyParserResult(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, []ingestion.Row{}).
		Once()

	source := NewSource("ucs", path, parser, mapper, 10)

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

func TestSource_StreamBatch_FetchedAtIsSharedAcrossBatches(t *testing.T) {
	path := writeTempFile(t, "test data")

	parser := new(mockParser)
	mapper := new(mockMapper)

	rows := []ingestion.Row{
		{"norad_id": "1"},
		{"norad_id": "2"},
		{"norad_id": "3"},
	}

	parser.
		On("Stream", mock.Anything, mock.Anything).
		Return(nil, rows).
		Once()

	for _, row := range rows {
		mapper.
			On("Map", row).
			Return(json.RawMessage(`{}`), nil).
			Once()
	}

	source := NewSource("ucs", path, parser, mapper, 2)

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
