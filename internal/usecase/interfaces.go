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
	GetByEmailSecurity(ctx context.Context, email string) (*entity.UserSecurity, error)
	Create(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error)
	CreateByProvider(ctx context.Context, dto dto.CreateUserByProviderDTO) (*entity.User, error)
	GetByProvider(ctx context.Context, provider string, email string) (*entity.User, error)
	UpdatePwd(ctx context.Context, email string, pwd string, salt string) error
	GetByID(ctx context.Context, id int) (*dto.UserDTO, error)
	UpdateProfile(ctx context.Context, data dto.UpdateProfileDTO) (*dto.UserDTO, error)
}

type PosterRepo interface {
	GetFlatsAll(ctx context.Context, dto dto.PostersFiltersDTO) ([]entity.PosterFlat, error)
	GetFlatsByIDs(ctx context.Context, ids []int) ([]entity.PosterFlat, error)
	CountPosters(ctx context.Context) (int, error)

	GetByAlias(ctx context.Context, posterAlias string, userID *int) (*entity.PosterById, error)
	GetFlatByPropetyID(ctx context.Context, propertyID int) (*entity.Flat, error)

	GetByUserID(ctx context.Context, userID int) ([]entity.Poster, error)
	GetMetroStationByRadius(ctx context.Context, buidingGeo dto.GeographyDTO, radius entity.Metre) ([]entity.MetroStation, error)

	CreateBuilding(ctx context.Context, poster *dto.PosterInput) (int, error)
	CreateProperty(ctx context.Context, poster *dto.PosterInput, buildingID int) (int, error)
	Create(ctx context.Context, poster *dto.PosterInput, propertyID int) (int, error)
	InsertFlat(ctx context.Context, flat *dto.FlatInput) error
	InsertFacilities(ctx context.Context, propertyID int, facilities []string) error
	InsertPhotos(ctx context.Context, posterID int, photos []dto.PhotoInput) error
	InsertMainPhoto(ctx context.Context, posterID int, avatarURL string) error

	GetUpdateIDsByAlias(ctx context.Context, alias string) (*dto.PosterUpdateIDs, error)
	Update(ctx context.Context, posterID int, poster *dto.PosterInput) error
	UpdateProperty(ctx context.Context, propertyID int, poster *dto.PosterInput) error
	UpdateBuilding(ctx context.Context, buildingID int, poster *dto.PosterInput) error
	UpdateFlat(ctx context.Context, flat *dto.FlatInput) error

	GetPhotoPathsByPosterID(ctx context.Context, posterID int) ([]string, error)

	DeleteFacilitiesByPropertyID(ctx context.Context, propertyID int) error
	DeletePhotosByPosterID(ctx context.Context, posterID int) error
	GetCityByName(ctx context.Context, name string) (*entity.City, error)
	CreateCity(ctx context.Context, name string) (*entity.City, error)

	Delete(ctx context.Context, posterID int) error
	DeleteFlat(ctx context.Context, propertyID int) error
	DeleteProperty(ctx context.Context, propertyID int) error
	DeleteBuilding(ctx context.Context, buildingID int) error

	AddView(ctx context.Context, userID int, posterID int)
	GetViewsCount(ctx context.Context, posterID int) (int, error)
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
	Order() OrderRepo
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

type SearchRepo interface {
	SearchPosters(ctx context.Context, filters dto.PostersFiltersDTO) (*dto.PostersResponse, error)
}

type SupportAgent interface {
	Chat(ctx context.Context, systemPrompt string, userPrompt string) (*dto.ChatResult, error)
}

type OrderRepo interface {
	Create(ctx context.Context, order *dto.Order) (int, error)
	InsertPhotos(ctx context.Context, orderID int, photos []dto.PhotoInput) error
	GetByUserID(ctx context.Context, userID int) ([]entity.Order, error)
	GetAll(ctx context.Context) ([]entity.Order, error)
}
