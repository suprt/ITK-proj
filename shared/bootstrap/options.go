package bootstrap

import (
	"time"

	"go.uber.org/fx"
)

type Options struct {
	Addr        string
	Modules     []fx.Option
	StopTimeout time.Duration
} //TODO:  log
