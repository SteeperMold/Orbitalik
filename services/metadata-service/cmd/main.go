package main

import (
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/bootstrap"
	infra "github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/infrastructure"
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

	ctx, cancel := bootstrap.SignalContext()
	defer cancel()

	go bootstrap.StartSchedulers(ctx, cfg, db, rdb, logger)
	go bootstrap.StartHTTPServer(cfg, db, logger)
	bootstrap.StartGRPCServer(cfg, db, logger)
}
