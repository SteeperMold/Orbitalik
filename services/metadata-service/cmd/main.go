package main

import (
	"log"

	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/bootstrap"
	infra "github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/infrastructure"
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

	rdb := infra.NewRedisClient(cfg.Redis)

	logger, err := applog.NewLogger(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func() {
		err := logger.Sync()
		log.Fatal(err)
	}()

	go bootstrap.StartSchedulers(ctx, cfg, db, rdb, logger)
	go bootstrap.StartHTTPServer(cfg, db, logger)
	bootstrap.StartGRPCServer(cfg, db, logger)
}
