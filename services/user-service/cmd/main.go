package main

import (
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
	"github.com/SteeperMold/Orbitalik/user-service/internal/bootstrap"
	"github.com/SteeperMold/Orbitalik/user-service/internal/infrastructure"
)

func main() {
	cfg := infrastructure.NewConfig()

	pgxPool := infrastructure.NewSQLDatabase(cfg.DB)
	defer infrastructure.CloseDBConnection(pgxPool)

	db := postgres.NewDB(pgxPool)

	zapLog := infrastructure.NewLogger(cfg.AppEnv)
	defer infrastructure.LoggerSync(zapLog)

	logger := applog.New(zapLog)

	_, cancel := bootstrap.SignalContext()
	defer cancel()

	go bootstrap.StartHTTPServer(cfg, db, logger)
	bootstrap.StartGRPCServer(cfg, db, logger)
}
