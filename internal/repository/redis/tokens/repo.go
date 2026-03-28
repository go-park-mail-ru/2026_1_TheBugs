package tokens

import (
	"context"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
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

func (r *TokenRepo) SetBlacklist(ctx context.Context, val string, ttl time.Duration) error {
	key := "blacklist:" + val
	return r.redisClient.Set(ctx, key, "1", ttl).Err()
}

func (r *TokenRepo) IsBlacklisted(ctx context.Context, val string) (bool, error) {
	key := "blacklist:" + val
	res, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return res == "1", nil
}

func (r *TokenRepo) CreateRecoverSession(ctx context.Context, sessionID string, data entity.RecoverSession, ttl time.Duration) error {
	key := "auth:recover:session:" + sessionID

	err := r.redisClient.HSet(ctx, key, map[string]interface{}{
		"email":    data.Email,
		"code":     data.Code,
		"attempts": data.Attempts,
		"verified": data.Verified,
	}).Err()
	if err != nil {
		return err
	}

	return r.redisClient.Expire(ctx, key, ttl).Err()
}

func (r *TokenRepo) GetRecoverSession(ctx context.Context, sessionID string) (*entity.RecoverSession, error) {
	key := "auth:recover:session:" + sessionID

	res, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	attempts, _ := strconv.Atoi(res["attempts"])
	verified := res["verified"] == "1"

	return &entity.RecoverSession{
		Email:    res["email"],
		Code:     res["code"],
		Attempts: attempts,
		Verified: verified,
	}, nil
}

func (r *TokenRepo) DeleteRecoverSession(ctx context.Context, sessionID string) error {
	key := "auth:recover:session:" + sessionID
	return r.redisClient.Del(ctx, key).Err()
}

func (r *TokenRepo) IncrementRecoverAttempts(ctx context.Context, sessionID string) (int64, error) {
	key := "auth:recover:session:" + sessionID

	attempts, err := r.redisClient.HIncrBy(ctx, key, "attempts", 1).Result()
	if err != nil {
		return 0, err
	}

	return attempts, nil
}

func (r *TokenRepo) SetRecoverVerified(ctx context.Context, sessionID string, verified bool) error {
	key := "auth:recover:session:" + sessionID

	return r.redisClient.HSet(ctx, key, "verified", verified).Err()
}
