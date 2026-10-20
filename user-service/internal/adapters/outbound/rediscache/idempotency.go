package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/suprt/trading/user-service/internal/core/domain"
)

type IdempotencyStore struct {
	client *redis.Client
}

func NewIdempotencyStore(client *redis.Client) *IdempotencyStore {
	return &IdempotencyStore{client: client}
}

func (s *IdempotencyStore) Get(ctx context.Context, key string) (*domain.IdempotencyRecord, error) {
	data, err := s.client.Get(ctx, idemKey(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var rec domain.IdempotencyRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (s *IdempotencyStore) Set(ctx context.Context, key string, rec *domain.IdempotencyRecord,
	ttl time.Duration) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, idemKey(key), data, ttl).Err()
}

func idemKey(key string) string {
	return fmt.Sprintf("idem:%s", key)
}
