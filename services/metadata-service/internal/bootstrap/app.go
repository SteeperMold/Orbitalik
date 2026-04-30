package bootstrap

import (
	"context"
	"log"
	"net"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/gen/metadatapb"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/celestrak"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/fallback"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/filesource"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/ucs"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/repository"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/scheduler"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/service"
	transportGrpc "github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/transport/grpc"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/transport/http/route"
)

func StartSchedulers(ctx context.Context, cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	ucsSource := ucs.NewSource(
		ucs.NewClient(cfg.UCS.SourceURL, cfg.UCS.FetchTimeout, cfg.UCS.FetchRetries),
		ucs.NewParser(),
		ucs.NewMapper(),
		cfg.UCS.BatchSize,
	)
	ucsFallbackFile := filesource.NewSource(string(models.SourceUCS), cfg.UCS.FallbackFilePath, ucs.NewParser(), ucs.NewMapper(), cfg.UCS.BatchSize)
	ucsWithFallback := fallback.NewFallbackSource(ucsSource, ucsFallbackFile, logger)
	celestrakSource := celestrak.NewSource(
		celestrak.NewClient(cfg.Celestrak.SourceURL, cfg.Celestrak.FetchTimeout, cfg.Celestrak.FetchRetries),
		celestrak.NewParser(),
		celestrak.NewMapper(),
		cfg.Celestrak.BatchSize,
	)
	celecstrakFallbackFile := filesource.NewSource(string(models.SourceCelestrak), cfg.Celestrak.FallbackFilePath, celestrak.NewParser(), celestrak.NewMapper(), cfg.Celestrak.BatchSize)
	celestrakWithFallback := fallback.NewFallbackSource(celestrakSource, celecstrakFallbackFile, logger)
	sources := []ingestion.MetadataSource{ucsWithFallback, celestrakWithFallback}
	metadataRepo := repository.NewMetadataRepository(db)
	ingestionService := ingestion.NewService(sources, metadataRepo, logger)
	fetchMetadataScheduler := scheduler.NewFetchMetadataScheduler(ingestionService, logger, cfg.IngestionInterval, cfg.IngestionTimeout)
	fetchMetadataScheduler.Start(ctx)
}

func StartHTTPServer(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	route.Serve(cfg, db, logger)
}

func StartGRPCServer(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	metadataRepo := repository.NewMetadataRepository(db)
	metadataService := service.NewMetadataService(metadataRepo)

	grpcServer := transportGrpc.NewServer(logger)
	metadatapb.RegisterSatelliteMetadataServiceServer(
		grpcServer,
		transportGrpc.NewMetadataServer(
			metadataService,
			logger,
			uint32(cfg.MaxPageSize),
			uint32(cfg.DefaultPageSize),
		),
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
