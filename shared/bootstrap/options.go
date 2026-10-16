package bootstrap

import (
	"time"

	"go.uber.org/fx"
)

type Options struct {
	Modules     []fx.Option
	StopTimeout time.Duration
} //TODO: addr, log
