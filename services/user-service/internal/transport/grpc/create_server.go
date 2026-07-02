package grpc

import (
	"github.com/SteeperMold/Orbitalik/common/go/grpc/interceptors"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"google.golang.org/grpc"
)

func NewServer(logger applog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryLoggingInterceptor(logger)),
	)

	return server
}
