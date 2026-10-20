package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Username string
	Password string
	DB       int
}

const contextTimeout = 5 * time.Second

func NewClient(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	return client, nil
}
