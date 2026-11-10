package main

import (
	"log"

	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	commonlog "github.com/SteeperMold/Orbitalik/common/go/log"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/aggregation"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/bootstrap"
	infra "github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/repository"
)

func main() {
	cfg := infra.NewConfig()

	ctx, cancel := bootstrap.SignalContext()
	defer cancel()

	db, err := postgres.OpenConn(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to open postgres conn: %v", err)
	}
	defer db.Close()

	logger, err := applog.NewLogger(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func(logger commonlog.Logger) {
		_ = logger.Sync()
	}(logger)

	rdb := infra.NewRedisClient(cfg.Redis)

	rawRepo := repository.NewRawMetadataRepository(db)
	metaRepo := repository.NewMetadataRepository(db)
	aggSvc := aggregation.NewService(rawRepo, metaRepo)

	if err := bootstrap.StartWorker(ctx, cfg, rdb, aggSvc, logger); err != nil {
		log.Fatalf("worker stopped: %v", err)
	}
}
