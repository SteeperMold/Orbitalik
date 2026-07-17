package bootstrap

import (
	"log"
	"net"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/user-service/gen/userpb"
	"github.com/SteeperMold/Orbitalik/user-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/user-service/internal/repository"
	"github.com/SteeperMold/Orbitalik/user-service/internal/service"
	transportGrpc "github.com/SteeperMold/Orbitalik/user-service/internal/transport/grpc"
	"github.com/SteeperMold/Orbitalik/user-service/internal/transport/http/route"
)

func StartHTTPServer(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	route.Serve(cfg, db, logger)
}

func StartGRPCServer(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db, logger)
	userService := service.NewUserService(userRepo, logger)

	grpcServer := transportGrpc.NewServer(logger)
	userpb.RegisterUserServiceServer(grpcServer, transportGrpc.NewUserServiceServer(userService, logger))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
