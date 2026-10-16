//минимальный конфиг, потом расширю

package config

import (
	"time"

	sharedconfig "github.com/suprt/trading/shared/config"
)

type GRPC struct {
	Addr        string        `env:"GRPC_ADDR" envDefault:":50051"`
	StopTimeout time.Duration `env:"GRPC_STOP_TIMEOUT" envDefault:"15s"`
}

type JWT struct {
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15s"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"168h"`
}

type SeedUser struct {
	Email    string
	Password string
	Roles    []string
}

type Seed struct {
	Enabled bool
	Users   []SeedUser
}
type Config struct {
	sharedconfig.Log
	sharedconfig.JWTSecret
	GRPC GRPC
	JWT  JWT
	Seed Seed `yaml:"seed"`
	Auth Auth `yaml:"auth"`
}

type Auth struct {
	Whitelist map[string]struct{}
}

func Load(yamlPath string) (*Config, error) {
	return sharedconfig.Load[Config](yamlPath)
}
