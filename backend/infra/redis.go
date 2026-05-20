package infra

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/huodaoshi/harness/backend/conf"
)

var (
	redisOnce   sync.Once
	redisClient redis.UniversalClient
	redisErr    error
)

// NewRedisClient returns a shared Redis UniversalClient.
func NewRedisClient(cfg conf.RedisConfig) (redis.UniversalClient, error) {
	redisOnce.Do(func() {
		client := redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    cfg.Addrs,
			Password: cfg.Password,
		})
		if err := client.Ping(context.Background()).Err(); err != nil {
			redisErr = fmt.Errorf("infra: redis ping: %w", err)
			return
		}
		redisClient = client
	})
	return redisClient, redisErr
}
