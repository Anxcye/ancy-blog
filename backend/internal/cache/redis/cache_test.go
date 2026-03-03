// File: cache_test.go
// Purpose: Verify Redis cache provider initialization behavior on invalid endpoint.
// Module: backend/internal/cache/redis, cache provider test layer.
// Related: cache.go.
package redis

import (
	"context"
	"testing"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/config"
)

func TestNewRedisCacheFailsOnInvalidAddr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	_, err := New(ctx, config.RedisConfig{Addr: "127.0.0.1:1", DB: 0, PoolSize: 1, MinIdleConns: 1})
	if err == nil {
		t.Fatalf("expected redis init to fail for unreachable address")
	}
}
