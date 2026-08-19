package worker

import (
	"context"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type mockAggregator struct {
	mock.Mock
}

func (m *mockAggregator) RebuildSatellite(ctx context.Context, noradID int) error {
	args := m.Called(ctx, noradID)
	return args.Error(0)
}

type mockRedis struct {
	mock.Mock
}

func (m *mockRedis) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd {
	args := m.Called(ctx, stream, group, start)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *mockRedis) XReadGroup(ctx context.Context, args *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	call := m.Called(ctx, args)
	return call.Get(0).(*redis.XStreamSliceCmd)
}

func (m *mockRedis) XAck(ctx context.Context, stream string, group string, ids ...string) *redis.IntCmd {
	args := m.Called(ctx, stream, group, ids)
	return args.Get(0).(*redis.IntCmd)
}

func (m *mockRedis) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, members)
	return args.Get(0).(*redis.IntCmd)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, fields ...applog.Field) {
	m.Called(msg, fields)
}

func (m *mockLogger) Sync() error {
	g := m.Called()
	return g.Error(0)
}

func (m *mockLogger) Error(msg string, fields ...applog.Field) {
	m.Called(msg, fields)
}

func newTestWorker(
	rdb RedisClient,
	aggregator Aggregator,
	logger applog.Logger,
) *Worker {
	return New(
		rdb,
		aggregator,
		logger,
		"consumer-1",
		"satellite-events",
		"metadata-workers",
		"dirty-satellites",
		10,
		50*time.Millisecond,
		time.Millisecond,
	)
}
