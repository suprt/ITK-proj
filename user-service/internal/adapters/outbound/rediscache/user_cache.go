package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/suprt/trading/user-service/internal/core/domain"
	er "github.com/suprt/trading/user-service/internal/core/errors"
)

type UserCacheConfig struct {
	UserTTL         time.Duration
	UserNotFoundTTL time.Duration
}

type UserCache struct {
	client          *redis.Client
	userTTL         time.Duration
	userNotFoundTTL time.Duration
}

func NewUserCache(client *redis.Client, cfg UserCacheConfig) *UserCache {
	return &UserCache{
		client:          client,
		userTTL:         cfg.UserTTL,
		userNotFoundTTL: cfg.UserNotFoundTTL,
	}
}

func (c *UserCache) Get(ctx context.Context, id domain.UserID) (*domain.User, error) {
	exists, err := c.client.Exists(ctx, notFoundKey(id)).Result()
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, er.ErrUserNotFound
	}

	data, err := c.client.Get(ctx, userKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, er.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *UserCache) Set(ctx context.Context, u *domain.User) error {
	data, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, userKey(u.ID), data, c.userTTL).Err()
}

func (c *UserCache) SetNotFound(ctx context.Context, id domain.UserID) error {
	return c.client.Set(ctx, notFoundKey(id), "1", c.userNotFoundTTL).Err()
}

func (c *UserCache) Delete(ctx context.Context, id domain.UserID) error {
	return c.client.Del(ctx, userKey(id)).Err()
}
func (c *UserCache) DeleteNotFound(ctx context.Context, id domain.UserID) error {
	return c.client.Del(ctx, notFoundKey(id)).Err()
}

func userKey(id domain.UserID) string     { return fmt.Sprintf("user:%s", string(id)) }
func notFoundKey(id domain.UserID) string { return fmt.Sprintf("user:%s:missing", string(id)) }
