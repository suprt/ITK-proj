package bootstrap

import (
	"github.com/suprt/trading/shared/grpcServer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewApp(opts Options) *fx.App {
	stopTimeout := opts.StopTimeout
	return fx.New(
		fx.Provide(func() grpcServer.Addr {
			return grpcServer.Addr(opts.Addr)
		}),
		fx.Provide(func() *zap.Logger {
			return opts.Log
		}),
		fx.StopTimeout(stopTimeout),
		fx.Provide(grpcServer.New),
		fx.Options(opts.Modules...))
}
