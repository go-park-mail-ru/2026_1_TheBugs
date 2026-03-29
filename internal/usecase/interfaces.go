package usecase

import (
	"context"
	"io"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_repo.go -package=mocks

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error)
	CreateByProvider(ctx context.Context, dto dto.CreateUserByProviderDTO) (*entity.User, error)
	GetByProvider(ctx context.Context, provider string, email string) (*entity.User, error)
	UpdatePwd(ctx context.Context, email string, pwd string, salt string) error
}

type PosterRepo interface {
	GetAll(ctx context.Context, dto dto.PostersFiltersDTO) ([]entity.Poster, error)
	CountPosters(ctx context.Context) (int, error)
	GetByAlias(ctx context.Context, posterAlias string) (*entity.PosterById, error)
	GetFlatByPropetyID(ctx context.Context, propertyID int) (*entity.Flat, error)
	GetByUserID(ctx context.Context, userID int) ([]entity.Poster, error)
	GetMetroStationByRadius(ctx context.Context, buidingGeo dto.GeographyDTO, radius entity.Metre) ([]entity.MetroStation, error)
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
	GetAllByDeveloperID(ctx context.Context, companyID int) ([]dto.UtilityCompanyCardDTO, error)
	GetAllDevelopers(ctx context.Context) ([]dto.DeveloperDTO, error)
}

type UnitOfWork interface {
	Users() UserRepo
	Posters() PosterRepo
	Autho() AuthRepo
	UtilityCompany() UtilityCompanyRepo
	Do(ctx context.Context, fn func(r UnitOfWork) error) error
}

type Сache interface {
	SetBlacklist(ctx context.Context, val string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, val string) (bool, error)
	CreateRecoverSession(ctx context.Context, sessionID string, data entity.RecoverSession, ttl time.Duration) error
	GetRecoverSession(ctx context.Context, sessionID string) (*entity.RecoverSession, error)
	DeleteRecoverSession(ctx context.Context, sessionID string) error
	IncrementRecoverAttempts(ctx context.Context, sessionID string) (int64, error)
	SetRecoverVerified(ctx context.Context, sessionID string, verified bool) error
}

type MailSender interface {
	SendCode(ctx context.Context, to string, code string) error
}

type FileRepo interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}
