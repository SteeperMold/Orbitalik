package interceptors

import (
	"context"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor(logger applog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)
		duration := time.Since(start)

		st := status.Convert(err)

		logger.Info("gRPC request",
			applog.NewField("method", info.FullMethod),
			applog.NewField("duration", duration),
			applog.NewField("code", st.Code().String()),
			applog.NewErrorField(err),
		)

		return resp, err
	}
}
