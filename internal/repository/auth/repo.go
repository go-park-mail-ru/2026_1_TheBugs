package auth

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

type AuthRepo struct {
	tokensSlice []entity.RefreshToken
}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		tokensSlice: []entity.RefreshToken{},
	}
}

func (r *AuthRepo) CreateToken(ctx context.Context, token dto.CreateRefreshTokenDTO) error {
	r.tokensSlice = append(r.tokensSlice, entity.RefreshToken{
		TokenID:   token.TokenID,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
	})
	return nil
}

func (r *AuthRepo) GetToken(ctx context.Context, tokenID string, userID int) (*entity.RefreshToken, error) {
	for _, token := range r.tokensSlice {
		if token.TokenID == tokenID && token.UserID == userID {
			return &token, nil
		}
	}
	return nil, entity.NotFoundError
}
func (r *AuthRepo) DeleteToken(ctx context.Context, tokenID string, userID int) error {
	for i, token := range r.tokensSlice {
		if token.TokenID == tokenID && token.UserID == userID {
			if i < len(r.tokensSlice)-1 {
				r.tokensSlice = append(r.tokensSlice[:i], r.tokensSlice[i+1:]...)
			} else {
				r.tokensSlice = r.tokensSlice[:i]
			}
			return nil
		}
	}
	return entity.NotFoundError
}
