//минимальный конфиг, потом расширю

package config

import (
	"time"
)

type GRPC struct {
	Addr string
}

type JWT struct {
	Secret     []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration
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
	GRPC GRPC
	JWT  JWT
	Seed Seed
}

func Default() Config {
	return Config{
		GRPC: GRPC{Addr: ":50051"},
		JWT: JWT{
			Secret:     []byte("dev-secret"),
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 7 * 24 * time.Hour,
		},
		Seed: Seed{
			Enabled: true,
			Users: []SeedUser{
				{
					Email:    "admin@local",
					Password: "admin12345",
					Roles:    []string{"USER", "ADMIN"},
				},
				{
					Email:    "trader@local",
					Password: "trader12345",
					Roles:    []string{"USER", "TRADER"},
				},
				{
					Email:    "user@local",
					Password: "user12345",
					Roles:    []string{"USER"},
				},
			},
		},
	}
}
