package ratelimit

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/redis/limits"
)

type RateLimitUC struct {
	repo *limits.RateLimitRepository
}

func NewRateLimitUC(repo *limits.RateLimitRepository) *RateLimitUC {
	return &RateLimitUC{repo: repo}
}

func (uc *RateLimitUC) CheckIPLimit(ctx context.Context, ip string, limit int, ttl time.Duration) (bool, error) {
	count, err := uc.repo.IncIPAttempts(ctx, ip, ttl)
	if err != nil {
		return false, err
	}
	return count <= int64(limit), nil
}
