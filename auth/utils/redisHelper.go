package utils

import (
	"time"

	"github.com/khaingminhtun/rssagg/config"
)

// redis helper
// SetKey stores a value in Redis with TTL
func SetKey(cfg *config.Config, key string, value interface{}, ttl time.Duration) error {
	return cfg.Redis.Set(cfg.Ctx, key, value, ttl).Err()
}

// GetKey retrieves a value from Redis
func GetKey(cfg *config.Config, key string) (string, error) {
	return cfg.Redis.Get(cfg.Ctx, key).Result()
}

// DeleteKey deletes a key from Redis
func DeleteKey(cfg *config.Config, key string) error {
	return cfg.Redis.Del(cfg.Ctx, key).Err()
}
