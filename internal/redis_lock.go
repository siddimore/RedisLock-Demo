package internal

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Wrapper
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient Create instance
func NewRedisClient() *RedisClient {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

// AcquireLock tries to acquire a distributed lock with a TTL.
func (r *RedisClient) AcquireLock(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, ttl).Result()
}

// ReleaseLock deletes the lock key. Should be called after work is done.
// TTL acts as a safety net if the process crashes before release.
func (r *RedisClient) ReleaseLock(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
