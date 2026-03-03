// File: cache.go
// Purpose: Implement cache abstraction using Redis.
// Module: backend/internal/cache/redis, concrete cache provider layer.
// Related: internal/cache interface and app dependency wiring.
package redis

import (
	"context"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/cache"
	"github.com/anxcye/ancy-blog/backend/internal/config"
	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *goredis.Client
}

var _ cache.Cache = (*Cache)(nil)

func New(ctx context.Context, cfg config.RedisConfig) (*Cache, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Cache{client: client}, nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool, error) {
	v, err := c.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return v, true, nil
}

func (c *Cache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}
