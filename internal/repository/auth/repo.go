package auth

import (
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

func (r *AuthRepo) CreateToken(token dto.CreateRefreshTokenDTO) error {
	r.tokensSlice = append(r.tokensSlice, entity.RefreshToken{
		TokenID:   token.TokenID,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
	})
	return nil
}

func (r *AuthRepo) GetToken(tokenID string, userID int) (*entity.RefreshToken, error) {
	for _, token := range r.tokensSlice {
		if token.TokenID == tokenID && token.UserID == userID {
			return &token, nil
		}
	}
	return nil, nil
}
func (r *AuthRepo) DeleteToken(tokenID string, userID int) error {
	for i, token := range r.tokensSlice {
		if token.TokenID == tokenID && token.UserID == userID {
			r.tokensSlice = append(r.tokensSlice[:i], r.tokensSlice[i+1:]...)
			return nil
		}
	}
	return nil
}
