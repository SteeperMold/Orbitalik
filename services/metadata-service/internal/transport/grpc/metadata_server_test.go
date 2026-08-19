package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/gen/metadatapb"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockMetadataService struct {
	mock.Mock
}

func (m *mockMetadataService) GetMetadataByNoradID(
	ctx context.Context,
	noradID int,
) (*models.SatelliteMetadata, error) {
	args := m.Called(ctx, noradID)

	var meta *models.SatelliteMetadata
	if args.Get(0) != nil {
		meta = args.Get(0).(*models.SatelliteMetadata)
	}

	return meta, args.Error(1)
}

func (m *mockMetadataService) GetMetadataByName(
	ctx context.Context,
	name string,
) (*models.SatelliteMetadata, error) {
	args := m.Called(ctx, name)

	var meta *models.SatelliteMetadata
	if args.Get(0) != nil {
		meta = args.Get(0).(*models.SatelliteMetadata)
	}

	return meta, args.Error(1)
}

func (m *mockMetadataService) ListSatellites(
	ctx context.Context,
	filter *models.ListFilter,
) ([]*models.SatelliteMetadata, string, uint32, error) {
	args := m.Called(ctx, filter)

	var items []*models.SatelliteMetadata
	if args.Get(0) != nil {
		items = args.Get(0).([]*models.SatelliteMetadata)
	}

	return items, args.String(1), args.Get(2).(uint32), args.Error(3)
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

func newTestMetadataServer(
	service MetadataService,
	logger applog.Logger,
) *MetadataServer {
	return NewMetadataServer(
		service,
		logger,
		100,
		20,
	)
}

func TestMetadataServer_GetSatelliteMetadata_ByNoradID(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	meta := &models.SatelliteMetadata{
		NoradID: 25544,
		Name:    "ISS",
	}

	service.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(meta, nil).
		Once()

	resp, err := server.GetSatelliteMetadata(ctx, &metadatapb.GetMetadataRequest{
		Identifier: &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_NoradId{
				NoradId: 25544,
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Metadata)

	assert.Equal(t, uint32(25544), resp.Metadata.NoradId)
	assert.Equal(t, "ISS", resp.Metadata.Name)

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestMetadataServer_GetSatelliteMetadata_ByName(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	meta := &models.SatelliteMetadata{
		NoradID: 25544,
		Name:    "ISS",
	}

	service.
		On("GetMetadataByName", ctx, "ISS").
		Return(meta, nil).
		Once()

	resp, err := server.GetSatelliteMetadata(ctx, &metadatapb.GetMetadataRequest{
		Identifier: &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_SatelliteName{
				SatelliteName: "ISS",
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Metadata)

	assert.Equal(t, uint32(25544), resp.Metadata.NoradId)
	assert.Equal(t, "ISS", resp.Metadata.Name)

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestMetadataServer_GetSatelliteMetadata_MissingIdentifier(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	resp, err := server.GetSatelliteMetadata(
		context.Background(),
		&metadatapb.GetMetadataRequest{},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(
		t,
		"either norad_id or name must be set",
		status.Convert(err).Message(),
	)

	service.AssertNotCalled(t, "GetMetadataByNoradID", mock.Anything, mock.Anything)
	service.AssertNotCalled(t, "GetMetadataByName", mock.Anything, mock.Anything)
}

func TestMetadataServer_GetSatelliteMetadata_ServiceErrorByNoradID(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()
	expectedErr := errors.New("database unavailable")

	service.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(nil, expectedErr).
		Once()

	logger.
		On("Error", "failed to get metadata by identifier", mock.Anything).
		Once()

	resp, err := server.GetSatelliteMetadata(ctx, &metadatapb.GetMetadataRequest{
		Identifier: &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_NoradId{
				NoradId: 25544,
			},
		},
	})

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(
		t,
		"failed to get metadata",
		status.Convert(err).Message(),
	)

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestMetadataServer_GetSatelliteMetadata_ServiceErrorByName(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()
	expectedErr := errors.New("repository failure")

	service.
		On("GetMetadataByName", ctx, "ISS").
		Return(nil, expectedErr).
		Once()

	logger.
		On("Error", "failed to get metadata by identifier", mock.Anything).
		Once()

	resp, err := server.GetSatelliteMetadata(ctx, &metadatapb.GetMetadataRequest{
		Identifier: &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_SatelliteName{
				SatelliteName: "ISS",
			},
		},
	})

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Internal, status.Code(err))

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestMetadataServer_GetSatelliteMetadata_NotFound(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	service.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(nil, nil).
		Once()

	resp, err := server.GetSatelliteMetadata(ctx, &metadatapb.GetMetadataRequest{
		Identifier: &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_NoradId{
				NoradId: 25544,
			},
		},
	})

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.NotFound, status.Code(err))
	assert.Equal(t, "metadata not found", status.Convert(err).Message())

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestMetadataServer_GetSatelliteMetadata_MapsMetadata(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	launchDate := time.Date(
		2020, 1, 2,
		3, 4, 5, 0,
		time.UTC,
	)

	meta := &models.SatelliteMetadata{
		NoradID:    25544,
		Name:       "ISS",
		LaunchDate: &launchDate,
		UpdatedAt:  launchDate,
	}

	service.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(meta, nil).
		Once()

	resp, err := server.GetSatelliteMetadata(ctx, &metadatapb.GetMetadataRequest{
		Identifier: &metadatapb.SatelliteIdentifier{
			Kind: &metadatapb.SatelliteIdentifier_NoradId{
				NoradId: 25544,
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Metadata)

	assert.Equal(t, uint32(25544), resp.Metadata.NoradId)
	assert.Equal(t, "ISS", resp.Metadata.Name)

	require.NotNil(t, resp.Metadata.LaunchDate)
	require.NotNil(t, resp.Metadata.UpdatedAt)

	assert.True(
		t,
		launchDate.Equal(resp.Metadata.LaunchDate.AsTime()),
	)

	assert.True(
		t,
		launchDate.Equal(resp.Metadata.UpdatedAt.AsTime()),
	)

	service.AssertExpectations(t)
}

func TestMetadataServer_ListSatelliteMetadata_DefaultPageSize(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	items := []*models.SatelliteMetadata{
		{
			NoradID: 1,
			Name:    "Satellite One",
		},
	}

	service.
		On(
			"ListSatellites",
			ctx,
			mock.MatchedBy(func(filter *models.ListFilter) bool {
				return filter.PageSize == 20 &&
					filter.PageToken == ""
			}),
		).
		Return(items, "", uint32(1), nil).
		Once()

	resp, err := server.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Len(t, resp.Items, 1)
	assert.Equal(t, uint32(1), resp.Total)
	assert.Empty(t, resp.NextPageToken)

	service.AssertExpectations(t)
}

func TestMetadataServer_ListSatelliteMetadata_ExplicitPageSize(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	service.
		On(
			"ListSatellites",
			ctx,
			mock.MatchedBy(func(filter *models.ListFilter) bool {
				return filter.PageSize == 50
			}),
		).
		Return([]*models.SatelliteMetadata{}, "", uint32(0), nil).
		Once()

	resp, err := server.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{
			PageSize: 50,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Empty(t, resp.Items)
	assert.Zero(t, resp.Total)

	service.AssertExpectations(t)
}

func TestMetadataServer_ListSatelliteMetadata_PageSizeTooLarge(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	resp, err := server.ListSatelliteMetadata(
		context.Background(),
		&metadatapb.ListSatelliteMetadataRequest{
			PageSize: 101,
		},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "page_size too large", status.Convert(err).Message())

	service.AssertNotCalled(t, "ListSatellites", mock.Anything, mock.Anything)
}

func TestMetadataServer_ListSatelliteMetadata_ConvertsFilters(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	objectType := metadatapb.ObjectType_OBJECT_TYPE_DEBRIS
	missionType := metadatapb.MissionType_MISSION_TYPE_SCIENCE
	statusValue := metadatapb.OperationalStatus_OPERATIONAL_STATUS_ACTIVE
	orbitRegime := metadatapb.OrbitRegime_ORBIT_REGIME_GEO
	constellation := "ISS"

	service.
		On("ListSatellites", ctx, mock.Anything).
		Return([]*models.SatelliteMetadata{}, "next", uint32(10), nil).
		Once()

	resp, err := server.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{
			PageSize:          20,
			PageToken:         "abc",
			ObjectType:        &objectType,
			MissionType:       &missionType,
			OperationalStatus: &statusValue,
			OrbitRegime:       &orbitRegime,
			Constellation:     &constellation,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "next", resp.NextPageToken)
	assert.Equal(t, uint32(10), resp.Total)

	service.AssertExpectations(t)
}

func TestMetadataServer_ListSatelliteMetadata_EmptyConstellationIsIgnored(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()
	empty := ""

	service.
		On(
			"ListSatellites",
			ctx,
			mock.MatchedBy(func(filter *models.ListFilter) bool {
				return filter.Constellation == nil
			}),
		).
		Return([]*models.SatelliteMetadata{}, "", uint32(0), nil).
		Once()

	resp, err := server.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{
			Constellation: &empty,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	service.AssertExpectations(t)
}

func TestMetadataServer_ListSatelliteMetadata_ServiceError(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()
	expectedErr := errors.New("database unavailable")

	service.
		On("ListSatellites", ctx, mock.Anything).
		Return(nil, "", uint32(0), expectedErr).
		Once()

	logger.
		On("Error", "failed to list satellite metadata", mock.Anything).
		Once()

	resp, err := server.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{
			PageSize: 20,
		},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(
		t,
		"failed to list satellite metadata",
		status.Convert(err).Message(),
	)

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestMetadataServer_ListSatelliteMetadata_MapsItems(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := newTestMetadataServer(service, logger)

	ctx := context.Background()

	items := []*models.SatelliteMetadata{
		{
			NoradID: 1,
			Name:    "One",
		},
		{
			NoradID: 2,
			Name:    "Two",
		},
	}

	service.
		On("ListSatellites", ctx, mock.Anything).
		Return(items, "2", uint32(5), nil).
		Once()

	resp, err := server.ListSatelliteMetadata(
		ctx,
		&metadatapb.ListSatelliteMetadataRequest{
			PageSize: 10,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, resp.Items, 2)

	assert.Equal(t, uint32(1), resp.Items[0].NoradId)
	assert.Equal(t, "One", resp.Items[0].Name)

	assert.Equal(t, uint32(2), resp.Items[1].NoradId)
	assert.Equal(t, "Two", resp.Items[1].Name)

	assert.Equal(t, "2", resp.NextPageToken)
	assert.Equal(t, uint32(5), resp.Total)

	service.AssertExpectations(t)
}

func TestNewMetadataServer(t *testing.T) {
	service := new(mockMetadataService)
	logger := new(mockLogger)

	server := NewMetadataServer(
		service,
		logger,
		100,
		20,
	)

	require.NotNil(t, server)
	assert.Same(t, service, server.service)
	assert.Same(t, logger, server.logger)
	assert.Equal(t, uint32(100), server.maxPageSize)
	assert.Equal(t, uint32(20), server.defaultPageSize)
}
