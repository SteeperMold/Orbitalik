package bootstrap

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/gen/tlepb"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/repository"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/scheduler"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/service"
	transportGrpc "github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/transport/grpc"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/transport/http/route"
	"go.uber.org/zap"
)

func StartSchedulers(ctx context.Context, cfg *infrastructure.Config, db domain.DBConn, logger *zap.Logger) {
	tleRepository := repository.NewTLERepository(db)

	httpClient := &http.Client{Timeout: cfg.TLEFetchTimeout}
	fetchService := service.NewFetchTLEService(tleRepository, httpClient, cfg.TLESourceUrl)

	tleScheduler := scheduler.NewTLEFetchScheduler(fetchService, logger, cfg.TLEFetchInterval, cfg.ContextTimeout)

	go tleScheduler.Start(ctx)
}

func StartHTTPServer(cfg *infrastructure.Config, db domain.DBConn, logger *zap.Logger) {
	go route.Serve(cfg, db, logger)
}

func StartGRPCServer(cfg *infrastructure.Config, db domain.DBConn, logger *zap.Logger) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	tleRepo := repository.NewTLERepository(db)
	tleService := service.NewTLEService(tleRepo)

	grpcServer := transportGrpc.NewServer(logger)
	tlepb.RegisterTleServiceServer(grpcServer, transportGrpc.NewTLEServiceServer(tleService, logger))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
