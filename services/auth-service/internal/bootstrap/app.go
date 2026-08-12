package bootstrap

import (
	"log"
	"net"

	"github.com/SteeperMold/Orbitalik/auth-service/gen/authpb"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/crypto/password"
	infra "github.com/SteeperMold/Orbitalik/auth-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/jwt"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/repository"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/service"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/transport/client"
	transportGrpc "github.com/SteeperMold/Orbitalik/auth-service/internal/transport/grpc"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/transport/http/route"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
)

func StartHTTPServer(cfg *infra.Config, db db.Conn, logger applog.Logger) {
	route.Serve(cfg, db, logger)
}

func StartGRPCServer(cfg *infra.Config, db db.Conn, logger applog.Logger) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := infra.LoadPrivateKey(cfg.JWT.PrivateKeyPath)
	if err != nil {
		log.Fatalf("failed to load private key: %s", err.Error())
	}

	publicKey, err := infra.LoadPublicKey(cfg.JWT.PublicKeyPath)
	if err != nil {
		log.Fatalf("failed to load private key: %s", err.Error())
	}

	tokenManager := jwt.NewTokenManager(
		privateKey,
		publicKey,
		"orbitalik-auth",
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	passwordHasher := password.NewHasher()

	credRepo := repository.NewCredentialsRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)

	userClient, err := client.NewUserClient(cfg.UserServiceAddress)
	if err != nil {
		log.Fatalf("failed to create user client: %s", err.Error())
	}
	defer func(userClient *client.Client) {
		_ = userClient.CloseConnection()
	}(userClient)

	authService := service.NewAuthService(credRepo, tokenRepo, tokenManager, passwordHasher, userClient)

	grpcServer := transportGrpc.NewServer(logger)
	authpb.RegisterAuthServiceServer(grpcServer, transportGrpc.NewAuthServiceServer(authService, logger))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
