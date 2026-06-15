package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/aggregation"
	infra "github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/infrastructure"
	worker "github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/queue"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/repository"
)

func main() {
	cfg := infra.NewConfig()

	pgxPool := infra.NewSQLDatabase(cfg.DB)
	defer infra.CloseDBConnection(pgxPool)
	db := postgres.NewDB(pgxPool)

	rdb := infra.NewRedisClient(cfg.Redis)

	zapLog := infra.NewLogger(cfg.AppEnv)
	defer infra.LoggerSync(zapLog)
	logger := applog.New(zapLog)

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	rawRepo := repository.NewRawMetadataRepository(db)
	metaRepo := repository.NewMetadataRepository(db)
	aggSvc := aggregation.NewService(rawRepo, metaRepo)

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
	)

	if err := w.Init(ctx); err != nil {
		log.Fatal(err)
	}

	if err := w.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
