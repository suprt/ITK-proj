package interceptors

import (
	"context"
	"strings"

	userv1 "github.com/suprt/trading/proto/user/v1"
	"github.com/suprt/trading/user-service/internal/core/auth"
	"github.com/suprt/trading/user-service/internal/core/domain"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func UnaryAuth(tokens outbound.TokenIssuer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if !requiresAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := meta.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		raw := strings.TrimPrefix(values[0], "Bearer ")
		if raw == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
		}

		userID, roles, _, err := tokens.ParseAccess(domain.AccessToken(raw))
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = auth.WithUserID(ctx, userID)
		ctx = auth.WithRoles(ctx, roles)
		return handler(ctx, req)
	}
}

func requiresAuth(fullMethod string) bool {
	name := strings.Replace(strings.TrimPrefix(fullMethod, "/"), "/", ".", 1)

	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil {
		return true
	}

	method, ok := desc.(protoreflect.MethodDescriptor)
	if !ok {
		return true
	}

	opts := method.Options()
	if !proto.HasExtension(opts, userv1.E_RequiresAuth) {
		return true
	}

	val, _ := proto.GetExtension(opts, userv1.E_RequiresAuth).(bool)
	return val
}
