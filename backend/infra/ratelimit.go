package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter wraps Redis for INCR-based rate limiting and simple KV ops.
type RedisRateLimiter struct {
	rdb redis.UniversalClient
}

// NewRedisRateLimiter creates a RedisRateLimiter.
func NewRedisRateLimiter(rdb redis.UniversalClient) *RedisRateLimiter {
	return &RedisRateLimiter{rdb: rdb}
}

// Allow reports whether the request is within limit for the sliding window starting at first hit.
// On Redis error it returns (true, err) so callers may fail open.
func (r *RedisRateLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	count, err := r.Incr(ctx, key)
	if err != nil {
		return true, err
	}
	if count == 1 {
		if expErr := r.rdb.Expire(ctx, key, window).Err(); expErr != nil {
			return true, expErr
		}
	}
	return count <= limit, nil
}

// Incr increments the counter at key and returns the new value.
func (r *RedisRateLimiter) Incr(ctx context.Context, key string) (int64, error) {
	val, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("infra: rate limiter: incr: %w", err)
	}
	return val, nil
}

// Get returns the string value stored at key.
func (r *RedisRateLimiter) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("infra: redis: get %s: %w", key, err)
	}
	return val, nil
}

// Set stores value at key with the given TTL (0 = no expiry).
func (r *RedisRateLimiter) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if err := r.rdb.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("infra: redis: set %s: %w", key, err)
	}
	return nil
}

// Del deletes one or more keys.
func (r *RedisRateLimiter) Del(ctx context.Context, keys ...string) error {
	if err := r.rdb.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("infra: redis: del: %w", err)
	}
	return nil
}
