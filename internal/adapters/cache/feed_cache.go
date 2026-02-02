package cache

import (
	"context"
	"time"
)

type RSSURLCache struct {
	store *RedisStore
	ttl   time.Duration
}

func NewRSSURLCache(store *RedisStore, ttl time.Duration) *RSSURLCache {
	return &RSSURLCache{store: store, ttl: ttl}
}

func (c *RSSURLCache) Get(ctx context.Context, websiteURL string) (string, bool, error) {
	key := RSSURLKey(websiteURL)
	return c.store.GetString(ctx, key)
}

func (c *RSSURLCache) Set(ctx context.Context, websiteURL, rssURL string) error {
	key := RSSURLKey(websiteURL)
	return c.store.SetString(ctx, key, rssURL, c.ttl)
}
