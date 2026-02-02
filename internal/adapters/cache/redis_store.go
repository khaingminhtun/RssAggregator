package cache

import (
	"context"
	"time"

	"github.com/khaingminhtun/rssagg/internal/pkg/errorHandle"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) GetString(ctx context.Context, key string) (string, bool, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil // cache miss
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (r *RedisStore) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisStore) Delete(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return errorHandle.CacheUnavailable("redis delete failed")
	}
	return nil
}
