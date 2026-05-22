package limits

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitRepository struct {
	redisClient redis.Cmdable
}

func NewRateLimitRepositoryRepo(rdb redis.Cmdable) *RateLimitRepository {
	return &RateLimitRepository{
		redisClient: rdb,
	}
}

func (r *RateLimitRepository) IncIPAttempts(ctx context.Context, ip string, ttl time.Duration) (int64, error) {
	key := "ratelimit:" + ip

	lua := `
        local key = KEYS[1]
        local ttl = tonumber(ARGV[1])
        local count = redis.call('INCR', key)
        if count == 1 then
            redis.call('EXPIRE', key, ttl)
        end
        return count
    `

	res, err := r.redisClient.Eval(ctx, lua, []string{key}, ttl.Seconds()).Int64()
	if err != nil {
		return 0, err
	}
	return res, nil
}
