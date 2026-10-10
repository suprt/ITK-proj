package bootstrap

import (
	"github.com/suprt/trading/shared/grpcServer"
	"go.uber.org/fx"
)

func NewApp(opts Options) *fx.App {
	return fx.New(
		fx.Provide(grpcServer.New),
		fx.Options(opts.Modules...))
}
