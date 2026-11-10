package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(
	ctx context.Context,
	query string,
	args ...any,
) (db.Rows, error) {
	ret := m.Called(ctx, query, args)
	return ret.Get(0).(db.Rows), ret.Error(1)
}

func (m *MockDB) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) db.Row {
	ret := m.Called(ctx, query, args)
	return ret.Get(0).(db.Row)
}

func (m *MockDB) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (db.Result, error) {
	ret := m.Called(ctx, query, args)
	return ret.Get(0).(db.Result), ret.Error(1)
}

func (m *MockDB) Begin(ctx context.Context) (db.Tx, error) {
	ret := m.Called(ctx)
	return ret.Get(0).(db.Tx), ret.Error(1)
}

func (m *MockDB) Ping(ctx context.Context) error {
	ret := m.Called(ctx)
	return ret.Error(0)
}

func (m *MockDB) Close() {
	m.Called()
}

type mockMetadataRepository struct {
	mock.Mock
}

func (m *mockMetadataRepository) GetMetadataByNoradID(
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

func (m *mockMetadataRepository) GetMetadataByName(
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

func (m *mockMetadataRepository) ListSatellites(
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
