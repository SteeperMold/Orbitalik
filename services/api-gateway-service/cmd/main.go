package main

import (
	"log"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/bootstrap"
	infra "github.com/SteeperMold/Orbitalik/api-gateway-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/adapter"
	commonlog "github.com/SteeperMold/Orbitalik/common/go/log"
	applog "github.com/SteeperMold/Orbitalik/common/go/log/zap"
)

func main() {
	cfg := infra.NewConfig()

	_, cancel := bootstrap.SignalContext()
	defer cancel()

	logger, err := applog.NewLogger(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer func(logger commonlog.Logger) {
		_ = logger.Sync()
	}(logger)

	connections, err := adapter.NewConnections(cfg)
	if err != nil {
		log.Fatalf("failed to connect to other microservice: %v", err)
	}
	defer func(connections *adapter.Connections) {
		_ = connections.Close()
	}(connections)

	bootstrap.StartHTTPServer(cfg, connections, logger)
}
