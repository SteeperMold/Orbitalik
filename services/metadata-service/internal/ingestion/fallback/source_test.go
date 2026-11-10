package fallback

import (
	"context"
	"errors"
	"testing"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSource struct {
	mock.Mock
	name string
}

func (m *mockSource) Name() string {
	return m.name
}

func (m *mockSource) StreamBatch(
	ctx context.Context,
	fn func([]*ingestion.SatelliteSourceRecord) error,
) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, fields ...applog.Field) {
	m.Called(msg, fields)
}

func (m *mockLogger) Sync() error {
	return nil
}

func (m *mockLogger) Error(msg string, fields ...applog.Field) {
	m.Called(msg, fields)
}

func TestSource_Name(t *testing.T) {
	primary := &mockSource{name: "ucs"}
	secondary := &mockSource{name: "celestrak"}

	source := NewFallbackSource(primary, secondary, nil)

	assert.Equal(t, "ucs_fallback", source.Name())
}

func TestSource_StreamBatch_PrimarySucceeds(t *testing.T) {
	primary := &mockSource{name: "ucs"}
	secondary := &mockSource{name: "celestrak"}
	logger := new(mockLogger)

	ctx := context.Background()
	callback := func([]*ingestion.SatelliteSourceRecord) error {
		return nil
	}

	primary.
		On("StreamBatch", ctx, mock.Anything).
		Return(nil).
		Once()

	source := NewFallbackSource(primary, secondary, logger)

	err := source.StreamBatch(ctx, callback)

	require.NoError(t, err)

	primary.AssertExpectations(t)
	secondary.AssertNotCalled(t, "StreamBatch", mock.Anything, mock.Anything)
	logger.AssertNotCalled(t, "Error", mock.Anything, mock.Anything)
}

func TestSource_StreamBatch_FallsBackToSecondary(t *testing.T) {
	primary := &mockSource{name: "ucs"}
	secondary := &mockSource{name: "celestrak"}
	logger := new(mockLogger)

	ctx := context.Background()
	primaryErr := errors.New("ucs unavailable")

	callback := func([]*ingestion.SatelliteSourceRecord) error {
		return nil
	}

	primary.
		On("StreamBatch", ctx, mock.Anything).
		Return(primaryErr).
		Once()

	secondary.
		On("StreamBatch", ctx, mock.Anything).
		Return(nil).
		Once()

	logger.
		On(
			"Error",
			"primary source failed, switching to fallback",
			mock.Anything,
		).
		Return().
		Once()

	source := NewFallbackSource(primary, secondary, logger)

	err := source.StreamBatch(ctx, callback)

	require.NoError(t, err)

	primary.AssertExpectations(t)
	secondary.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestSource_StreamBatch_ReturnsSecondaryError(t *testing.T) {
	primary := &mockSource{name: "ucs"}
	secondary := &mockSource{name: "celestrak"}
	logger := new(mockLogger)

	ctx := context.Background()

	primaryErr := errors.New("ucs unavailable")
	secondaryErr := errors.New("celestrak unavailable")

	callback := func([]*ingestion.SatelliteSourceRecord) error {
		return nil
	}

	primary.
		On("StreamBatch", ctx, mock.Anything).
		Return(primaryErr).
		Once()

	secondary.
		On("StreamBatch", ctx, mock.Anything).
		Return(secondaryErr).
		Once()

	logger.
		On(
			"Error",
			"primary source failed, switching to fallback",
			mock.Anything,
		).
		Return().
		Once()

	source := NewFallbackSource(primary, secondary, logger)

	err := source.StreamBatch(ctx, callback)

	require.ErrorIs(t, err, secondaryErr)

	primary.AssertExpectations(t)
	secondary.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestSource_StreamBatch_PassesSameCallbackToFallback(t *testing.T) {
	primary := &mockSource{name: "ucs"}
	secondary := &mockSource{name: "celestrak"}
	logger := new(mockLogger)

	ctx := context.Background()
	primaryErr := errors.New("primary failed")

	var receivedCallback func([]*ingestion.SatelliteSourceRecord) error

	primary.
		On("StreamBatch", ctx, mock.Anything).
		Return(primaryErr).
		Once()

	secondary.
		On("StreamBatch", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			receivedCallback = args.Get(1).(func([]*ingestion.SatelliteSourceRecord) error)
		}).
		Return(nil).
		Once()

	logger.
		On(
			"Error",
			"primary source failed, switching to fallback",
			mock.Anything,
		).
		Return().
		Once()

	source := NewFallbackSource(primary, secondary, logger)

	callback := func([]*ingestion.SatelliteSourceRecord) error {
		return nil
	}

	err := source.StreamBatch(ctx, callback)

	require.NoError(t, err)
	require.NotNil(t, receivedCallback)

	_ = receivedCallback(nil)

	primary.AssertExpectations(t)
	secondary.AssertExpectations(t)
}
