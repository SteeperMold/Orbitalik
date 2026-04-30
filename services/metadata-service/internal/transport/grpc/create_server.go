package grpc

import (
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"google.golang.org/grpc"
)

func NewServer(logger applog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryLoggingInterceptor(logger)),
	)

	return server
}
