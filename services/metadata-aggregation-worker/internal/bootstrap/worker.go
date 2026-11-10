package bootstrap

import (
	"context"
	"fmt"

	commonlog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/aggregation"
	infra "github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/infrastructure"
	worker "github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/queue"
	"github.com/redis/go-redis/v9"
)

func StartWorker(
	ctx context.Context,
	cfg *infra.Config,
	rdb *redis.Client,
	aggSvc *aggregation.Service,
	logger commonlog.Logger,
) error {
	w := worker.New(
		rdb,
		aggSvc,
		logger,
		cfg.Redis.ConsumerName,
		cfg.Redis.DirtySatellitesQueueKey,
		cfg.Redis.GroupName,
		cfg.Redis.DirtySatellitesSetKey,
		int64(cfg.Redis.StreamsCount),
		cfg.Redis.StreamsBlock,
		cfg.Redis.RetryDelay,
	)

	if err := w.Init(ctx); err != nil {
		return fmt.Errorf("initialize worker: %w", err)
	}

	if err := w.Run(ctx); err != nil {
		return fmt.Errorf("run worker: %w", err)
	}

	return nil
}
