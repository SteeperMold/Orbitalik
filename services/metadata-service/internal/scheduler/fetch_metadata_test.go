package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockIngestionService struct {
	mock.Mock
}

func (m *mockIngestionService) IngestMetadata(ctx context.Context) error {
	args := m.Called(ctx)
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

func newTestScheduler(
	service domain.IngestionService,
	logger applog.Logger,
	interval time.Duration,
	timeout time.Duration,
) *FetchMetadataScheduler {
	return NewFetchMetadataScheduler(
		service,
		logger,
		interval,
		timeout,
	)
}

func TestFetchMetadataScheduler_Start_RunsInitialFetch(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	ctx := context.Background()

	service.
		On("IngestMetadata", mock.Anything).
		Return(nil).
		Once()

	logger.
		On("Info", "starting metadata fetch scheduler", mock.Anything).
		Once()

	logger.
		On("Info", "ingested metadata", mock.Anything).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		time.Hour,
		time.Second,
	)

	done := make(chan struct{})

	go func() {
		scheduler.Start(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return len(service.Calls) == 1
	}, time.Second, time.Millisecond*10)

	cancelCtx, cancel := context.WithCancel(ctx)
	_ = cancelCtx
	cancel()

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestFetchMetadataScheduler_Start_RunsPeriodically(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service.
		On("IngestMetadata", mock.Anything).
		Return(nil).
		Times(3)

	logger.
		On("Info", "starting metadata fetch scheduler", mock.Anything).
		Once()

	logger.
		On("Info", "ingested metadata", mock.Anything).
		Times(3)

	logger.
		On("Info", "metadata fetch scheduler stopped", mock.Anything).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		10*time.Millisecond,
		time.Second,
	)

	done := make(chan struct{})

	go func() {
		scheduler.Start(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return len(service.Calls) >= 3
	}, time.Second, 5*time.Millisecond)

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("scheduler did not stop")
	}

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestFetchMetadataScheduler_Start_LogsServiceError(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expectedErr := errors.New("ingestion failed")
	called := make(chan struct{})

	service.
		On("IngestMetadata", mock.Anything).
		Run(func(args mock.Arguments) {
			close(called)
		}).
		Return(expectedErr).
		Once()

	logger.
		On("Info", "starting metadata fetch scheduler", mock.Anything).
		Once()

	logger.
		On("Error", "error ingesting metadata", mock.Anything).
		Once()

	logger.
		On("Info", "metadata fetch scheduler stopped", mock.Anything).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		time.Hour,
		time.Second,
	)

	done := make(chan struct{})

	go func() {
		scheduler.Start(ctx)
		close(done)
	}()

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("scheduler did not call IngestMetadata")
	}

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("scheduler did not stop")
	}

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestFetchMetadataScheduler_Start_StopsOnContextCancellation(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	ctx, cancel := context.WithCancel(context.Background())

	service.
		On("IngestMetadata", mock.Anything).
		Return(nil).
		Once()

	logger.
		On("Info", "starting metadata fetch scheduler", mock.Anything).
		Once()

	logger.
		On("Info", "ingested metadata", mock.Anything).
		Once()

	logger.
		On(
			"Info",
			"metadata fetch scheduler stopped",
			mock.Anything,
		).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		time.Hour,
		time.Second,
	)

	done := make(chan struct{})

	go func() {
		scheduler.Start(ctx)
		close(done)
	}()

	require.Eventually(t, func() bool {
		return len(service.Calls) == 1
	}, time.Second, 5*time.Millisecond)

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("scheduler did not stop after context cancellation")
	}

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestFetchMetadataScheduler_RunFetch_UsesContextTimeout(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	expectedTimeout := 20 * time.Millisecond

	service.
		On("IngestMetadata", mock.Anything).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)

			<-ctx.Done()

			require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
		}).
		Return(context.DeadlineExceeded).
		Once()

	logger.
		On("Error", "error ingesting metadata", mock.Anything).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		time.Hour,
		expectedTimeout,
	)

	scheduler.runFetch(context.Background())

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestFetchMetadataScheduler_RunFetch_LogsSuccess(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	service.
		On("IngestMetadata", mock.Anything).
		Return(nil).
		Once()

	logger.
		On("Info", "ingested metadata", mock.Anything).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		time.Hour,
		time.Second,
	)

	scheduler.runFetch(context.Background())

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestFetchMetadataScheduler_RunFetch_LogsError(t *testing.T) {
	service := new(mockIngestionService)
	logger := new(mockLogger)

	expectedErr := errors.New("database unavailable")

	service.
		On("IngestMetadata", mock.Anything).
		Return(expectedErr).
		Once()

	logger.
		On("Error", "error ingesting metadata", mock.Anything).
		Once()

	scheduler := newTestScheduler(
		service,
		logger,
		time.Hour,
		time.Second,
	)

	scheduler.runFetch(context.Background())

	service.AssertExpectations(t)
	logger.AssertExpectations(t)
}
