package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const CacheKeyPrefix = "cache:"

type RedisCache struct {
	client redis.Cmdable
}

func NewRedisCache(rdb redis.Cmdable) *RedisCache {
	return &RedisCache{
		client: rdb,
	}
}
func buildKey(key string) string {
	return CacheKeyPrefix + key
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, buildKey(key)).Bytes()
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, buildKey(key), value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, buildKey(key)).Err()
}
