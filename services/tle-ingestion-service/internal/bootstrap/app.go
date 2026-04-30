package bootstrap

import (
	"context"
	"log"
	"net"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/gen/tlepb"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/repository"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/scheduler"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/service"
	transportGrpc "github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/transport/grpc"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/transport/http/route"
)

func StartSchedulers(ctx context.Context, cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	tleRepository := repository.NewTLERepository(db)
	fetchService := ingestion.NewFetchTLEService(
		tleRepository,
		cfg.TLESourceUrl,
		cfg.TLEFetchTimeout,
		cfg.TLEFetchMaxRetries,
	)

	tleFetchScheduler := scheduler.NewTLEFetchScheduler(
		fetchService,
		logger,
		cfg.TLEFetchInterval,
		cfg.TLEFetchTimeout,
	)
	tleCleanupScheduler := scheduler.NewTLECleanupScheduler(
		tleRepository,
		logger,
		cfg.TLECleanupInterval,
		cfg.TLECleanupTimeout,
		cfg.TLERetentionPeriod,
	)

	go tleFetchScheduler.Start(ctx)
	tleCleanupScheduler.Start(ctx)
}

func StartHTTPServer(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	route.Serve(cfg, db, logger)
}

func StartGRPCServer(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
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
