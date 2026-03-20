package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_repo.go -package=mocks

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error)
	CreateUserByProvider(ctx context.Context, dto dto.CreateUserByProviderDTO) (*entity.User, error)
	GetUserByProvider(ctx context.Context, provider string, email string) (*entity.User, error)
}

type PosterRepo interface {
	GetPosters(ctx context.Context, dto dto.PostersFiltersDTO) ([]entity.Poster, error)
	CountPosters(ctx context.Context) (int, error)
	GetPosterByAlias(ctx context.Context, posterAlias string) (*entity.PosterById, error)
	GetFlatByPropetyID(ctx context.Context, propertyID int) (*entity.Flat, error)
}

type AuthRepo interface {
	GetToken(ctx context.Context, tokenID string, userID int) (*entity.RefreshToken, error)
	CreateToken(ctx context.Context, dto dto.CreateRefreshTokenDTO) error
	DeleteToken(ctx context.Context, tokenID string, userID int) error
	BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

type UtilityCompanyRepo interface {
	GetUtilityCompanyByAlias(ctx context.Context, alias string) (*dto.UtilityCompanyDTO, error)
}
