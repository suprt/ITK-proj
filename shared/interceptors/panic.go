package interceptors

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecovery(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				fields := []zap.Field{
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", debug.Stack()),
				}
				if id, ok := RequestIDFromContext(ctx); ok {
					fields = append(fields, zap.String("request_id", id))
				}
				logger.Error("panic recovery", fields...)
				err = status.Errorf(codes.Internal, "internal server error")
				resp = nil
			}
		}()
		return handler(ctx, req)
	}
}
