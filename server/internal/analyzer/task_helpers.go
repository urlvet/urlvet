package analyzer

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/urlvet/urlvet/internal/metrics"
)

// cachedTask executes a task with caching support
// Returns true if value was found in cache, false otherwise
func cachedTask[T any](
	ctx context.Context,
	cache CacheInterface,
	cacheKey string,
	ttl time.Duration,
	fetch func() (T, error),
	setOutput func(*Output, T),
	out *Output,
) (bool, error) {
	// Derive a short label from the cache key prefix (e.g. "domain_rank" from "domain_rank:example.com").
	keyLabel := cacheKey
	if i := strings.IndexByte(cacheKey, ':'); i >= 0 {
		keyLabel = cacheKey[:i]
	}

	// Try cache first
	if cache != nil {
		var cached T
		if err := cache.GetJSON(ctx, cacheKey, &cached); err == nil {
			metrics.CacheHits.WithLabelValues(keyLabel).Inc()
			out.mu.Lock()
			setOutput(out, cached)
			out.mu.Unlock()
			return true, nil
		} else if err != redis.Nil {
			// Cache error (not a miss) - continue to fetch
		} else {
			metrics.CacheMisses.WithLabelValues(keyLabel).Inc()
		}
	}

	// Cache miss - fetch from source
	value, err := fetch()
	if err != nil {
		return false, err
	}

	// Store in cache
	if cache != nil {
		_ = cache.SetJSON(ctx, cacheKey, value, ttl)
	}

	out.mu.Lock()
	setOutput(out, value)
	out.mu.Unlock()
	return false, nil
}

// Helper to update output fields with mutex lock
func updateOutput(fn func(*Output)) func(*Output) {
	return func(out *Output) {
		out.mu.Lock()
		defer out.mu.Unlock()
		fn(out)
	}
}

// Helper to update output fields with mutex lock (direct version)
func updateOutputDirect(out *Output, fn func(*Output)) {
	out.mu.Lock()
	defer out.mu.Unlock()
	fn(out)
}
