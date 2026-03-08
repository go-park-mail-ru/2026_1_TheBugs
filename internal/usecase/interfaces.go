package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_repo.go -package=mocks

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error)
}

type PosterRepo interface {
	GetPosters(ctx context.Context, dto dto.PostersFiltersDTO) ([]entity.Poster, error)
	CountPosters(ctx context.Context) (int, error)
}
type AuthRepo interface {
	GetToken(ctx context.Context, tokenID string, userID int) (*entity.RefreshToken, error)
	CreateToken(ctx context.Context, dto dto.CreateRefreshTokenDTO) error
	DeleteToken(ctx context.Context, tokenID string, userID int) error
}
