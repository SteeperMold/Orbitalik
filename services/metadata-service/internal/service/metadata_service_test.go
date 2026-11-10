package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadataService_GetMetadataByNoradID(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()

	meta := &models.SatelliteMetadata{
		NoradID: 25544,
		Name:    "ISS",
	}

	repo.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(meta, nil).
		Once()

	got, err := service.GetMetadataByNoradID(ctx, 25544)

	require.NoError(t, err)
	assert.Same(t, meta, got)

	repo.AssertExpectations(t)
}

func TestMetadataService_GetMetadataByNoradID_ReturnsError(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()
	expectedErr := errors.New("repository error")

	repo.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(nil, expectedErr).
		Once()

	got, err := service.GetMetadataByNoradID(ctx, 25544)

	assert.Nil(t, got)
	assert.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
}

func TestMetadataService_GetMetadataByNoradID_ReturnsNilResult(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()

	repo.
		On("GetMetadataByNoradID", ctx, 25544).
		Return(nil, nil).
		Once()

	got, err := service.GetMetadataByNoradID(ctx, 25544)

	require.NoError(t, err)
	assert.Nil(t, got)

	repo.AssertExpectations(t)
}

func TestMetadataService_GetMetadataByName(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()

	meta := &models.SatelliteMetadata{
		NoradID: 25544,
		Name:    "ISS",
	}

	repo.
		On("GetMetadataByName", ctx, "ISS").
		Return(meta, nil).
		Once()

	got, err := service.GetMetadataByName(ctx, "ISS")

	require.NoError(t, err)
	assert.Same(t, meta, got)

	repo.AssertExpectations(t)
}

func TestMetadataService_GetMetadataByName_ReturnsError(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()
	expectedErr := errors.New("repository error")

	repo.
		On("GetMetadataByName", ctx, "ISS").
		Return(nil, expectedErr).
		Once()

	got, err := service.GetMetadataByName(ctx, "ISS")

	assert.Nil(t, got)
	assert.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
}

func TestMetadataService_GetMetadataByName_ReturnsNilResult(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()

	repo.
		On("GetMetadataByName", ctx, "ISS").
		Return(nil, nil).
		Once()

	got, err := service.GetMetadataByName(ctx, "ISS")

	require.NoError(t, err)
	assert.Nil(t, got)

	repo.AssertExpectations(t)
}

func TestMetadataService_ListSatellites(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()

	filter := &models.ListFilter{
		PageSize: 20,
	}

	items := []*models.SatelliteMetadata{
		{
			NoradID: 1,
			Name:    "Satellite One",
		},
		{
			NoradID: 2,
			Name:    "Satellite Two",
		},
	}

	repo.
		On("ListSatellites", ctx, filter).
		Return(items, "2", uint32(5), nil).
		Once()

	gotItems, nextToken, total, err := service.ListSatellites(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, items, gotItems)
	assert.Equal(t, "2", nextToken)
	assert.Equal(t, uint32(5), total)

	repo.AssertExpectations(t)
}

func TestMetadataService_ListSatellites_ReturnsError(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()
	expectedErr := errors.New("repository error")

	filter := &models.ListFilter{
		PageSize: 20,
	}

	repo.
		On("ListSatellites", ctx, filter).
		Return(nil, "", uint32(0), expectedErr).
		Once()

	items, nextToken, total, err := service.ListSatellites(ctx, filter)

	assert.Nil(t, items)
	assert.Empty(t, nextToken)
	assert.Zero(t, total)
	assert.ErrorIs(t, err, expectedErr)

	repo.AssertExpectations(t)
}

func TestMetadataService_ListSatellites_NilFilter(t *testing.T) {
	repo := new(mockMetadataRepository)
	service := NewMetadataService(repo)

	ctx := context.Background()

	repo.
		On("ListSatellites", ctx, (*models.ListFilter)(nil)).
		Return([]*models.SatelliteMetadata{}, "", uint32(0), nil).
		Once()

	items, nextToken, total, err := service.ListSatellites(ctx, nil)

	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Empty(t, nextToken)
	assert.Zero(t, total)

	repo.AssertExpectations(t)
}

func TestNewMetadataService(t *testing.T) {
	repo := new(mockMetadataRepository)

	service := NewMetadataService(repo)

	require.NotNil(t, service)
	assert.Same(t, repo, service.repo)
}
