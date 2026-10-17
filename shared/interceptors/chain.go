package interceptors

import (
	"context"
	"slices"

	"google.golang.org/grpc"
)

func Chain(inter ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		h := handler
		for _, interceptor := range slices.Backward(inter) {

			next := h
			h = func(ctx context.Context, req any) (any, error) {
				return interceptor(ctx, req, info, next)
			}
		}
		return h(ctx, req)
	}
}
