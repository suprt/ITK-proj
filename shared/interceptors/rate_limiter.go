package interceptors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type RateLimiterRule struct {
	Limit  int
	Window time.Duration
}

type RateLimiterConfig struct {
	Default RateLimiterRule
	Methods map[string]RateLimiterRule
}

type ErrorReporter interface {
	Report(ctx context.Context, err error, attrs map[string]string)
}

func UnaryRateLimiter(client *redis.Client, cfg RateLimiterConfig, reporter ErrorReporter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		rule, ok := cfg.Methods[info.FullMethod]
		if !ok {
			rule = cfg.Default
		}
		if rule.Limit <= 0 {
			return handler(ctx, req)
		}

		ip, err := clientIP(ctx)
		if err != nil {
			return handler(ctx, req) // не блокируем если не получилось определить IP
		}

		key := fmt.Sprintf("ratelimit:%s:%s", info.FullMethod, ip)
		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			if reporter != nil {
				reporter.Report(ctx, err, map[string]string{
					"op":     "rate_limit_incr",
					"method": info.FullMethod,
					"key":    key,
				})
			}
			return handler(ctx, req)
		}
		if count == 1 {
			err = client.Expire(ctx, key, rule.Window).Err()
			if err != nil {
				if reporter != nil {
					reporter.Report(ctx, err, map[string]string{
						"op":     "rate_limit_expire",
						"method": info.FullMethod,
						"key":    key,
					})
				}
			}
		}

		if count > int64(rule.Limit) {
			return nil, status.Error(codes.ResourceExhausted, "too many requests")
		}
		return handler(ctx, req)
	}
}

// TODO: relay ip -> last x-forwarder-for -> p.Addr
func clientIP(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-forwarded-for"); len(vals) > 0 {
			return vals[0], nil
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			return addr.IP.String(), nil
		}
	}
	return "", errors.New("no client ip found")
}
