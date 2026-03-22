package tokens

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepo struct {
	redisClient *redis.Client
}

func NewTokenRepo(rdb *redis.Client) *TokenRepo {
	return &TokenRepo{
		redisClient: rdb,
	}
}

func (r *TokenRepo) BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error {
	if r.redisClient == nil {
		return nil
	}
	key := "blacklist:access:" + tokenID
	return r.redisClient.Set(ctx, key, "1", ttl).Err()
}

func (r *TokenRepo) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	if r.redisClient == nil {
		return false, nil
	}
	key := "blacklist:access:" + tokenID
	res, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return res == "1", nil
}
