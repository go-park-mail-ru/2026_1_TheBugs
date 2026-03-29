package usecase

import (
	"context"
	"io"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_repo.go -package=mocks

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error)
	CreateByProvider(ctx context.Context, dto dto.CreateUserByProviderDTO) (*entity.User, error)
	GetByProvider(ctx context.Context, provider string, email string) (*entity.User, error)
}

type PosterRepo interface {
	GetAll(ctx context.Context, dto dto.PostersFiltersDTO) ([]entity.Poster, error)
	CountPosters(ctx context.Context) (int, error)
	GetByAlias(ctx context.Context, posterAlias string) (*entity.PosterById, error)
	GetFlatByPropetyID(ctx context.Context, propertyID int) (*entity.Flat, error)
	CreateBuilding(ctx context.Context, poster *entity.PosterInput) (int, error)
	CreateProperty(ctx context.Context, poster *entity.PosterInput, buildingID int) (int, error)
	Create(ctx context.Context, poster *entity.PosterInput, propertyID int) (int, error)
	InsertFlat(ctx context.Context, flat *entity.FlatInput) error
	InsertPhotos(ctx context.Context, posterID int, photos []entity.PhotoInput) error
	InsertMainPhoto(ctx context.Context, posterID int, avatarURL string) error
}

type AuthRepo interface {
	GetToken(ctx context.Context, tokenID string, userID int) (*entity.RefreshToken, error)
	CreateToken(ctx context.Context, dto dto.CreateRefreshTokenDTO) error
	DeleteToken(ctx context.Context, tokenID string, userID int) error
}

type UtilityCompanyRepo interface {
	GetByAlias(ctx context.Context, alias string) (*dto.UtilityCompanyDTO, error)
}

type UnitOfWork interface {
	Users() UserRepo
	Posters() PosterRepo
	Autho() AuthRepo
	UtilityCompany() UtilityCompanyRepo
	Do(ctx context.Context, fn func(r UnitOfWork) error) error
}

type TokenRepo interface {
	BlacklistToken(ctx context.Context, tokenID string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, tokenID string) (bool, error)
}

type FileRepo interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}
