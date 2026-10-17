package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logger(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
		}

		if id, ok := RequestIDFromContext(ctx); ok {
			fields = append(fields, zap.String("request_id", id))
		}

		if err != nil {
			stat, _ := status.FromError(err)
			fields = append(fields, zap.String("status_code", stat.Code().String()),
				zap.String("status_message", stat.Message()),
			)
			logger.Warn("request failed", fields...)
			return resp, err
		}
		logger.Debug("request done", fields...)
		return resp, nil

	}
}
