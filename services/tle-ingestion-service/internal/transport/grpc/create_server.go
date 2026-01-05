package grpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(logger *zap.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryLoggingInterceptor(logger)),
	)

	return server
}
