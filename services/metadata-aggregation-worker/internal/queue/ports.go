package worker

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Aggregator interface {
	RebuildSatellite(ctx context.Context, noradID int) error
}

type RedisClient interface {
	XGroupCreateMkStream(ctx context.Context, stream, group, start string) *redis.StatusCmd
	XReadGroup(ctx context.Context, args *redis.XReadGroupArgs) *redis.XStreamSliceCmd
	XAck(ctx context.Context, stream string, group string, ids ...string) *redis.IntCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
}
