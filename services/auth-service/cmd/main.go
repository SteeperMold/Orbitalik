package main

import (
	"log"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/bootstrap"
	infra "github.com/SteeperMold/Orbitalik/auth-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	commonlog "github.com/SteeperMold/Orbitalik/common/go/log"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
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

	go bootstrap.StartHTTPServer(cfg, db, logger)
	bootstrap.StartGRPCServer(cfg, db, logger)
}
