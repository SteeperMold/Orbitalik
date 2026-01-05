package main

import (
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/bootstrap"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/infrastructure"
)

func main() {
	cfg := infrastructure.NewConfig()

	db := infrastructure.NewSQLDatabase(cfg.DB)
	defer infrastructure.CloseDBConnection(db)

	logger := infrastructure.NewLogger(cfg.AppEnv)
	defer infrastructure.LoggerSync(logger)

	ctx, cancel := bootstrap.SignalContext()
	defer cancel()

	go bootstrap.StartSchedulers(ctx, cfg, db, logger)
	go bootstrap.StartHTTPServer(cfg, db, logger)
	bootstrap.StartGRPCServer(cfg, db, logger)
}

//func main() {
//	cfg := infrastructure.NewConfig()
//
//	db := infrastructure.NewSQLDatabase(cfg.DB)
//	defer infrastructure.CloseDBConnection(db)
//
//	logger := infrastructure.NewLogger(cfg.General.AppEnv)
//	defer infrastructure.LoggerSync(logger)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	sigs := make(chan os.Signal, 1)
//	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//	go func() {
//		<-sigs
//		log.Println("shutting down rebalancer service")
//		cancel()
//	}()
//
//	tleRepository := repository.NewTLERepository(db)
//	fetchTLEService := service.NewFetchTLEService(tleRepository, http.DefaultClient, cfg.General.TLESourceUrl)
//	fetchTLEScheduler := scheduler.NewTLEFetchScheduler(fetchTLEService, logger, cfg.General.TLEFetchInterval, cfg.General.ContextTimeout)
//
//	go func() {
//		fetchTLEScheduler.Start(ctx)
//	}()
//
//	go func() {
//		route.Serve(cfg, db, logger)
//	}()
//
//	lis, err := net.Listen("tcp", ":"+cfg.General.GRPCPort)
//	if err != nil {
//		log.Fatal(err)
//	}
//	grpcServer := transportGrpc.NewServer(logger)
//	tleService := service.NewTLEService(tleRepository)
//	tlepb.RegisterTleServiceServer(grpcServer, transportGrpc.NewTLEServiceServer(tleService, logger))
//
//	err = grpcServer.Serve(lis)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
