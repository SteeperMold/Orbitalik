//go:build integration
// +build integration

package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWorker_Init_CreatesConsumerGroup(t *testing.T) {
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})

	w := newTestWorker(rdb, nil, nil)

	err := w.Init(context.Background())

	require.NoError(t, err)

	groups, err := rdb.XInfoGroups(context.Background(), "satellite-events").Result()
	require.NoError(t, err)

	require.Len(t, groups, 1)
	assert.Equal(t, "metadata-workers", groups[0].Name)
}

func TestWorker_Init_IgnoresExistingGroup(t *testing.T) {
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})

	w := newTestWorker(rdb, nil, nil)

	require.NoError(t, w.Init(context.Background()))
	require.NoError(t, w.Init(context.Background()))
}

func TestWorker_Init_ReturnsRedisError(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:1",
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})

	w := newTestWorker(rdb, nil, nil)

	err := w.Init(context.Background())

	require.Error(t, err)
}

func TestIsBusyGroup(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "busy group",
			err:  errors.New("BUSYGROUP Consumer Group name already exists"),
			want: true,
		},
		{
			name: "other redis error",
			err:  errors.New("connection refused"),
			want: false,
		},
		{
			name: "nil",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isBusyGroup(tt.err))
		})
	}
}

func TestWorker_Run_ReturnsContextError(t *testing.T) {
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	w := newTestWorker(rdb, nil, nil)

	err := w.Run(ctx)

	assert.ErrorIs(t, err, context.Canceled)
}

func TestWorker_Run_ProcessesMessage(t *testing.T) {
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})

	aggregator := new(mockAggregator)

	aggregator.
		On("RebuildSatellite", mock.Anything, 25544).
		Return(nil).
		Once()

	w := newTestWorker(rdb, aggregator, nil)

	require.NoError(t, w.Init(context.Background()))

	_, err := rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "satellite-events",
		Values: map[string]any{
			"norad_id": 25544,
		},
	}).Result()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = w.Run(ctx)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	aggregator.AssertExpectations(t)
}

func TestWorker_Run_ContinuesAfterProcessingError(t *testing.T) {
	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	t.Cleanup(func() {
		_ = rdb.Close()
	})

	logger := new(mockLogger)

	logger.
		On("Error", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	aggregator := new(mockAggregator)

	aggregator.
		On("RebuildSatellite", mock.Anything, 1).
		Return(errors.New("aggregation failed")).
		Once()

	aggregator.
		On("RebuildSatellite", mock.Anything, 2).
		Return(nil).
		Once()

	w := newTestWorker(rdb, aggregator, logger)

	require.NoError(t, w.Init(context.Background()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "satellite-events",
		Values: map[string]any{
			"norad_id": 1,
		},
	}).Result()
	require.NoError(t, err)

	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "satellite-events",
		Values: map[string]any{
			"norad_id": 2,
		},
	}).Result()
	require.NoError(t, err)

	err = w.Run(ctx)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	aggregator.AssertExpectations(t)
}

func TestWorker_Run_RetriesAfterRedisError(t *testing.T) {
	rdb := new(mockRedis)
	logger := new(mockLogger)

	redisErr := errors.New("redis connection lost")

	cmd := redis.NewXStreamSliceCmd(context.Background())
	cmd.SetErr(redisErr)

	rdb.
		On(
			"XReadGroup",
			mock.Anything,
			mock.Anything,
		).
		Return(cmd, redisErr).
		Once()

	ctx, cancel := context.WithCancel(context.Background())

	w := newTestWorker(rdb, nil, logger)
	w.retryDelay = 10 * time.Millisecond

	logger.
		On("Error", "failed to read stream", mock.Anything).
		Once()

	done := make(chan error, 1)

	go func() {
		done <- w.Run(ctx)
	}()

	require.Eventually(t, func() bool {
		return len(rdb.Calls) >= 1
	}, time.Second, time.Millisecond)

	cancel()

	err := <-done

	assert.ErrorIs(t, err, context.Canceled)
	rdb.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestWorker_Run_ContinuesAfterRedisError(t *testing.T) {
	rdb := new(mockRedis)
	logger := new(mockLogger)

	redisErr := errors.New("redis connection lost")

	errCmd := redis.NewXStreamSliceCmd(context.Background())
	errCmd.SetErr(redisErr)

	retryObserved := make(chan struct{})
	var calls int

	rdb.
		On("XReadGroup", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			calls++

			if calls == 2 {
				close(retryObserved)
			}
		}).
		Return(errCmd)

	logger.
		On("Error", "failed to read stream", mock.Anything).
		Maybe()

	w := newTestWorker(rdb, nil, logger)
	w.retryDelay = time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- w.Run(ctx)
	}()

	select {
	case <-retryObserved:
	case <-time.After(time.Second):
		t.Fatal("worker did not retry XReadGroup")
	}

	cancel()

	err := <-done

	assert.ErrorIs(t, err, context.Canceled)
	assert.GreaterOrEqual(t, calls, 2)

	rdb.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestSleep_ReturnsNilAfterDuration(t *testing.T) {
	ctx := context.Background()

	start := time.Now()

	err := sleep(ctx, 5*time.Millisecond)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, time.Since(start), 5*time.Millisecond)
}

func TestSleep_ReturnsContextErrorWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sleep(ctx, time.Second)

	assert.ErrorIs(t, err, context.Canceled)
}
