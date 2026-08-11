package main

import (
	"log"

	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
	"github.com/SteeperMold/Orbitalik/user-service/internal/bootstrap"
	"github.com/SteeperMold/Orbitalik/user-service/internal/infrastructure"
)

func main() {
	cfg := infrastructure.NewConfig()

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
	defer logger.Sync()

	go bootstrap.StartHTTPServer(cfg, db, logger)
	bootstrap.StartGRPCServer(cfg, db, logger)
}
