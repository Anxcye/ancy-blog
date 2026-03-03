// File: cache.go
// Purpose: Define cache abstraction used by services for cache-aside behavior.
// Module: backend/internal/cache, infrastructure abstraction layer.
// Related: redis cache implementation and content service.
package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}
