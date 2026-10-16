package bootstrap

import (
	"time"

	"github.com/suprt/trading/shared/grpcServer"
	"go.uber.org/fx"
)

func NewApp(opts Options) *fx.App {
	stopTimeout := opts.StopTimeout
	if stopTimeout == 0 {
		stopTimeout = 15 * time.Second
	}

	return fx.New(
		fx.Provide(func() grpcServer.Addr {
			return grpcServer.Addr(opts.Addr)
		}),
		fx.StopTimeout(stopTimeout),
		fx.Provide(grpcServer.New),
		fx.Options(opts.Modules...))
}
