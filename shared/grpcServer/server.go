package grpcServer

import (
	"context"
	"errors"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Params struct {
	fx.In
	LC           fx.Lifecycle
	Shutdowner   fx.Shutdowner
	Addr         Addr
	Interceptors []grpc.UnaryServerInterceptor `group:"grpc_server_interceptor"`
}
type Addr string
type Result struct {
	fx.Out
	Server *grpc.Server
}

func New(p Params) Result {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(p.Interceptors...))

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", string(p.Addr))
			if err != nil {
				return err
			}
			go func() {
				err := srv.Serve(lis)

				if err != nil &&
					!errors.Is(err, grpc.ErrServerStopped) &&
					!errors.Is(err, net.ErrClosed) {
					log.Printf("grpc serve failed: %v", err)
					if shutdownErr := p.Shutdowner.Shutdown(); shutdownErr != nil {
						log.Printf("shutdown failed: %v", shutdownErr)
					}
				}

			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			done := make(chan struct{})
			go func() {
				srv.GracefulStop()
				close(done)
			}()
			select {
			case <-done:
				return nil
			case <-ctx.Done():
				srv.Stop()
				return ctx.Err()
			}

		},
	})

	return Result{Server: srv}
}
