package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockMetadataSource struct {
	mock.Mock
	name string
}

func (m *mockMetadataSource) Name() string {
	return m.name
}

func (m *mockMetadataSource) StreamBatch(
	ctx context.Context,
	fn func([]*SatelliteSourceRecord) error,
) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type mockRawMetadataRepository struct {
	mock.Mock
}

func (m *mockRawMetadataRepository) SaveRawBatch(
	ctx context.Context,
	data []*models.SatelliteIngestRecord,
) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type mockDirtyMarker struct {
	mock.Mock
}

func (m *mockDirtyMarker) MarkDirty(
	ctx context.Context,
	noradIDs []int,
) error {
	args := m.Called(ctx, noradIDs)
	return args.Error(0)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, fields ...applog.Field) {
	m.Called(msg, fields)
}

func (m *mockLogger) Error(msg string, fields ...applog.Field) {
	m.Called(msg, fields)
}

func (m *mockLogger) Sync() error {
	return nil
}

func TestToIngestRecord(t *testing.T) {
	fetchedAt := time.Date(
		2026, 8, 18,
		10, 0, 0, 0,
		time.UTC,
	)

	raw := json.RawMessage(`{
		"norad_id": 25544,
		"cospar_id": "1998-067A",
		"name": "ISS"
	}`)

	src := &SatelliteSourceRecord{
		Source:    models.SourceUCS,
		FetchedAt: fetchedAt,
		Raw:       raw,
	}

	before := time.Now()

	got, err := toIngestRecord(src)

	after := time.Now()

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, 25544, got.NoradID)

	require.NotNil(t, got.CosparID)
	assert.Equal(t, "1998-067A", *got.CosparID)

	assert.Equal(t, models.SourceUCS, got.Source)
	assert.Equal(t, raw, got.Payload)
	assert.Equal(t, fetchedAt, got.FetchedAt)

	assert.False(t, got.StoredAt.Before(before))
	assert.False(t, got.StoredAt.After(after))
}

func TestToIngestRecord_WithoutCosparID(t *testing.T) {
	raw := json.RawMessage(`{"norad_id": 25544}`)

	got, err := toIngestRecord(&SatelliteSourceRecord{
		Source:    models.SourceCelestrak,
		FetchedAt: time.Now(),
		Raw:       raw,
	})

	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, 25544, got.NoradID)
	assert.Nil(t, got.CosparID)
}

func TestToIngestRecord_InvalidJSON(t *testing.T) {
	got, err := toIngestRecord(&SatelliteSourceRecord{
		Raw: []byte(`not valid json`),
	})

	assert.Nil(t, got)
	require.Error(t, err)
}

func TestService_IngestSource_Success(t *testing.T) {
	ctx := context.Background()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	fetchedAt := time.Now()

	batch := []*SatelliteSourceRecord{
		{
			Source:    models.SourceUCS,
			FetchedAt: fetchedAt,
			Raw: json.RawMessage(`{
				"norad_id": 25544,
				"cospar_id": "1998-067A"
			}`),
		},
		{
			Source:    models.SourceUCS,
			FetchedAt: fetchedAt,
			Raw: json.RawMessage(`{
				"norad_id": 12345,
				"cospar_id": "2024-001A"
			}`),
		},
	}

	source.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func([]*SatelliteSourceRecord) error)
			require.NoError(t, fn(batch))
		}).
		Return(nil).
		Once()

	repo.
		On(
			"SaveRawBatch",
			ctx,
			mock.MatchedBy(func(records []*models.SatelliteIngestRecord) bool {
				return len(records) == 2 &&
					records[0].NoradID == 25544 &&
					records[1].NoradID == 12345
			}),
		).
		Return(nil).
		Once()

	dirty.
		On("MarkDirty", ctx, []int{25544, 12345}).
		Return(nil).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Info", "ingestion completed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.NoError(t, err)

	source.AssertExpectations(t)
	repo.AssertExpectations(t)
	dirty.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_IngestSource_EmptyBatch(t *testing.T) {
	ctx := context.Background()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	source.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func([]*SatelliteSourceRecord) error)
			require.NoError(t, fn(nil))
		}).
		Return(nil).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Info", "ingestion completed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.NoError(t, err)

	repo.AssertNotCalled(t, "SaveRawBatch", mock.Anything, mock.Anything)
	dirty.AssertNotCalled(t, "MarkDirty", mock.Anything, mock.Anything)
}

func TestService_IngestSource_SkipsInvalidRecords(t *testing.T) {
	ctx := context.Background()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	batch := []*SatelliteSourceRecord{
		{
			Source:    models.SourceUCS,
			FetchedAt: time.Now(),
			Raw:       []byte(`invalid json`),
		},
		{
			Source:    models.SourceUCS,
			FetchedAt: time.Now(),
			Raw:       json.RawMessage(`{"norad_id": 25544}`),
		},
	}

	source.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func([]*SatelliteSourceRecord) error)
			require.NoError(t, fn(batch))
		}).
		Return(nil).
		Once()

	repo.
		On(
			"SaveRawBatch",
			ctx,
			mock.MatchedBy(func(records []*models.SatelliteIngestRecord) bool {
				return len(records) == 1 && records[0].NoradID == 25544
			}),
		).
		Return(nil).
		Once()

	dirty.
		On("MarkDirty", ctx, []int{25544}).
		Return(nil).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Info", "ingestion completed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.NoError(t, err)

	repo.AssertExpectations(t)
	dirty.AssertExpectations(t)
}

func TestService_IngestSource_SaveRawBatchError(t *testing.T) {
	ctx := context.Background()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	saveErr := errors.New("database unavailable")

	batch := []*SatelliteSourceRecord{
		{
			Source:    models.SourceUCS,
			FetchedAt: time.Now(),
			Raw:       json.RawMessage(`{"norad_id": 25544}`),
		},
	}

	source.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func([]*SatelliteSourceRecord) error)
			require.ErrorIs(t, fn(batch), saveErr)
		}).
		Return(saveErr).
		Once()

	repo.
		On("SaveRawBatch", ctx, mock.Anything).
		Return(saveErr).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Error", "ingestion failed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.ErrorIs(t, err, saveErr)

	dirty.AssertNotCalled(t, "MarkDirty", mock.Anything, mock.Anything)
	logger.AssertExpectations(t)
}

func TestService_IngestSource_DirtyMarkerErrorIsIgnored(t *testing.T) {
	ctx := context.Background()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	dirtyErr := errors.New("redis unavailable")

	batch := []*SatelliteSourceRecord{
		{
			Source:    models.SourceUCS,
			FetchedAt: time.Now(),
			Raw:       json.RawMessage(`{"norad_id": 25544}`),
		},
	}

	source.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func([]*SatelliteSourceRecord) error)
			require.NoError(t, fn(batch))
		}).
		Return(nil).
		Once()

	repo.
		On("SaveRawBatch", ctx, mock.Anything).
		Return(nil).
		Once()

	dirty.
		On("MarkDirty", ctx, []int{25544}).
		Return(dirtyErr).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Error", "failed to mark satellites dirty", mock.Anything).
		Once()

	logger.
		On("Info", "ingestion completed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.NoError(t, err)

	repo.AssertExpectations(t)
	dirty.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_IngestSource_SourceError(t *testing.T) {
	ctx := context.Background()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	sourceErr := errors.New("source unavailable")

	source.
		On("StreamBatch", ctx, mock.Anything).
		Return(sourceErr).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Error", "ingestion failed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.ErrorIs(t, err, sourceErr)

	repo.AssertNotCalled(t, "SaveRawBatch", mock.Anything, mock.Anything)
	dirty.AssertNotCalled(t, "MarkDirty", mock.Anything, mock.Anything)
}

func TestService_IngestSource_ContextCancelledBeforeBatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	source := &mockMetadataSource{name: "ucs"}
	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	source.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func([]*SatelliteSourceRecord) error)

			err := fn([]*SatelliteSourceRecord{
				{
					Source:    models.SourceUCS,
					FetchedAt: time.Now(),
					Raw:       json.RawMessage(`{"norad_id": 25544}`),
				},
			})

			require.ErrorIs(t, err, context.Canceled)
		}).
		Return(context.Canceled).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Once()

	logger.
		On("Error", "ingestion failed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{source},
		repo,
		dirty,
		logger,
	)

	err := service.ingestSource(ctx, source)

	require.ErrorIs(t, err, context.Canceled)

	repo.AssertNotCalled(t, "SaveRawBatch", mock.Anything, mock.Anything)
	dirty.AssertNotCalled(t, "MarkDirty", mock.Anything, mock.Anything)
}

func TestService_IngestMetadata_ContinuesAfterSourceFailure(t *testing.T) {
	ctx := context.Background()

	first := &mockMetadataSource{name: "ucs"}
	second := &mockMetadataSource{name: "celestrak"}

	repo := new(mockRawMetadataRepository)
	dirty := new(mockDirtyMarker)
	logger := new(mockLogger)

	sourceErr := errors.New("ucs unavailable")

	first.
		On("StreamBatch", ctx, mock.Anything).
		Return(sourceErr).
		Once()

	second.
		On("StreamBatch", ctx, mock.Anything).
		Return(nil).
		Once()

	logger.
		On("Info", "starting ingestion", mock.Anything).
		Twice()

	logger.
		On("Error", "ingestion failed", mock.Anything).
		Once()

	logger.
		On("Error", "failed to ingest source", mock.Anything).
		Once()

	logger.
		On("Info", "ingestion completed", mock.Anything).
		Once()

	service := NewService(
		[]MetadataSource{first, second},
		repo,
		dirty,
		logger,
	)

	err := service.IngestMetadata(ctx)

	// Public API intentionally treats individual source failures
	// as non-fatal.
	require.NoError(t, err)

	first.AssertExpectations(t)
	second.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_IngestMetadata_NilSources(t *testing.T) {
	service := NewService(nil, nil, nil, nil)

	assert.NoError(t, service.IngestMetadata(context.Background()))
}
