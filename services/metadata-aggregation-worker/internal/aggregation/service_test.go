package aggregation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRawRepository struct {
	mock.Mock
}

func (m *mockRawRepository) GetByNoradID(
	ctx context.Context,
	noradID int,
) ([]models.SatelliteIngestRecord, error) {
	args := m.Called(ctx, noradID)

	var records []models.SatelliteIngestRecord
	if value := args.Get(0); value != nil {
		records = value.([]models.SatelliteIngestRecord)
	}

	return records, args.Error(1)
}

type mockMetadataRepository struct {
	mock.Mock
}

func (m *mockMetadataRepository) Upsert(
	ctx context.Context,
	meta *models.SatelliteMetadata,
) error {
	args := m.Called(ctx, meta)
	return args.Error(0)
}

func TestService_RebuildSatellite_RawRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("raw repository failed")

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return(nil, expectedErr).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	require.ErrorIs(t, err, expectedErr)

	rawRepo.AssertExpectations(t)
	metaRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestService_RebuildSatellite_NoRecords(t *testing.T) {
	ctx := context.Background()

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return([]models.SatelliteIngestRecord{}, nil).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	require.NoError(t, err)

	rawRepo.AssertExpectations(t)
	metaRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestService_RebuildSatellite_NilRecords(t *testing.T) {
	ctx := context.Background()

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return(nil, nil).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	require.NoError(t, err)

	rawRepo.AssertExpectations(t)
	metaRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}

func TestService_RebuildSatellite_UpsertsAggregatedMetadata(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id":  12345,
			"cospar_id": "2024-001A",
			"name":      "Satellite Alpha",
			"operator":  "Operator Alpha",
		}),
	}

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return(records, nil).
		Once()

	metaRepo.
		On(
			"Upsert",
			ctx,
			mock.MatchedBy(func(meta *models.SatelliteMetadata) bool {
				return meta != nil &&
					meta.NoradID == 12345 &&
					meta.CosparID != nil &&
					*meta.CosparID == "2024-001A" &&
					meta.Name == "Satellite Alpha" &&
					meta.Operator != nil &&
					*meta.Operator == "Operator Alpha"
			}),
		).
		Return(nil).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	require.NoError(t, err)

	rawRepo.AssertExpectations(t)
	metaRepo.AssertExpectations(t)
}

func TestService_RebuildSatellite_UpsertError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("upsert failed")
	now := time.Now()

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, now, map[string]any{
			"norad_id": 12345,
			"name":     "Satellite Alpha",
		}),
	}

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return(records, nil).
		Once()

	metaRepo.
		On(
			"Upsert",
			ctx,
			mock.AnythingOfType("*models.SatelliteMetadata"),
		).
		Return(expectedErr).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	require.ErrorIs(t, err, expectedErr)

	rawRepo.AssertExpectations(t)
	metaRepo.AssertExpectations(t)
}

func TestService_RebuildSatellite_PassesNoradIDToRepository(t *testing.T) {
	ctx := context.Background()

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	rawRepo.
		On("GetByNoradID", ctx, 99999).
		Return([]models.SatelliteIngestRecord{}, nil).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 99999)

	require.NoError(t, err)

	rawRepo.AssertExpectations(t)
}

func TestService_RebuildSatellite_PassesContextToRepositories(t *testing.T) {
	//nolint:staticcheck // string context key is fine in tests
	ctx := context.WithValue(context.Background(), "test-key", "test-value")

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	records := []models.SatelliteIngestRecord{
		makeRecord(t, models.SourceUCS, time.Now(), map[string]any{
			"norad_id": 12345,
			"name":     "Satellite Alpha",
		}),
	}

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return(records, nil).
		Once()

	metaRepo.
		On(
			"Upsert",
			ctx,
			mock.AnythingOfType("*models.SatelliteMetadata"),
		).
		Return(nil).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	require.NoError(t, err)

	rawRepo.AssertExpectations(t)
	metaRepo.AssertExpectations(t)
}

func TestService_RebuildSatellite_DoesNotUpsertWhenAggregationProducesNil(t *testing.T) {
	ctx := context.Background()

	rawRepo := new(mockRawRepository)
	metaRepo := new(mockMetadataRepository)

	rawRepo.
		On("GetByNoradID", ctx, 12345).
		Return(nil, nil).
		Once()

	service := NewService(rawRepo, metaRepo)

	err := service.RebuildSatellite(ctx, 12345)

	assert.NoError(t, err)
	metaRepo.AssertNotCalled(t, "Upsert", mock.Anything, mock.Anything)
}
