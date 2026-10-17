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
	GRPC     GRPC
	Postgres Postgres
	JWT      JWT
	Crypto   Crypto
	Seed     Seed `yaml:"seed"`
	Auth     Auth `yaml:"auth"`
}

type Auth struct {
	Whitelist map[string]struct{}
}

type Crypto struct {
	HashCost int `env:"CRYPTO_HASH_COST" envDefault:"10"`
}

type Postgres struct {
	DSN             string        `env:"POSTGRES_DSN,required"`
	MaxConns        int32         `env:"POSTGRES_MAX_CONNS" envDefault:"25"`
	MinConns        int32         `env:"POSTGRES_MIN_CONNS" envDefault:"5"`
	MaxConnLifetime time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME" envDefault:"5m"`
}

func Load(yamlPath string) (*Config, error) {
	return sharedconfig.Load[Config](yamlPath)
}
