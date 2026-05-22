package delivery

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	jwtUtils "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/jwt"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_uc.go -package=mocks

type AuthUseCase interface {
	RegisterUseCase(ctx context.Context, data dto.CreateUserDTO) error
	LoginUseCase(ctx context.Context, email, password string) (*dto.UserAccessCredDTO, error)
	RefreshTokenUseCase(ctx context.Context, refreshToken string) (*dto.UserAccessCredDTO, error)
	LogoutUseCase(ctx context.Context, logoutCred dto.LogoutDTO) error
	LoginUserFromVKUseCase(ctx context.Context, flow dto.OAuthCodeFlow) (*dto.UserAccessCredDTO, error)
	LoginUserFromYandexUseCase(ctx context.Context, flow dto.OAuthCodeFlow) (*dto.UserAccessCredDTO, error)
	SendRecoveryCode(ctx context.Context, email string) (string, error)
	SendVerificationEmailCode(ctx context.Context, email string) (string, error)
	VerifyUserEmail(ctx context.Context, sessionID string, code string) error
	CheckRecoveryCode(ctx context.Context, sessionID string, code string) error
	UpdateUserPassword(ctx context.Context, sessionID string, password string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (*jwtUtils.Claims, error)
}

type UserUseCase interface {
	GetByID(ctx context.Context, userID int) (*dto.UserDTO, error)
	UpdateProfile(ctx context.Context, data dto.UpdateProfileRequest) (*dto.UserDTO, error)

	GetRoommateUser(ctx context.Context, userID int) (*dto.RoommateUserProfileDTO, error)
	AddRoommateMatch(ctx context.Context, fromUserID int, toUserID int, posterAlias *string) error
	GetRoommateContacts(ctx context.Context, fromUserID int, toUserID int) (*dto.RoommateContactsDTO, error)
	CreateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error
	GetRoommateForm(ctx context.Context, userID int) (*dto.RoommateFormDTO, error)
	UpdateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error
	GetIncomingRoommateMatches(ctx context.Context, userID int) (*dto.RoommateMatchesResponse, error)
	GetMatchedRoommateMatches(ctx context.Context, userID int) (*dto.RoommateMatchesResponse, error)
}

type PostersUseCase interface {
	SearchPostersUseCase(ctx context.Context, filters dto.PostersFiltersDTO) (*dto.PostersResponse, error)
	GetPosterByAliasUseCase(ctx context.Context, alias string, userID *int) (*dto.PosterDTO, error)
	GetPosterByUserID(ctx context.Context, userID int) ([]dto.MyPosterDTO, error)

	AddViewPoster(ctx context.Context, alias string, userID int) error
	GetViewsPoster(ctx context.Context, alias string) (int, error)

	AddFavoritePoster(ctx context.Context, alias string, userID int) error
	GetFavoritesPoster(ctx context.Context, userID int) (*dto.PostersResponse, error)
	DeleteFavoritePoster(ctx context.Context, alias string, userID int) error
	GetFavoritesCountPoster(ctx context.Context, posterAlias string, userID *int) (int, bool, error)

	GetPostersByCoords(ctx context.Context, bounds dto.MapBounds, filters dto.PostersFiltersDTO) (*dto.GeoJSONFeatureResponse, error)
	GetPostersByRadius(ctx context.Context, point dto.GeographyDTO) ([]dto.MyPosterDTO, error)

	GenerateDescription(ctx context.Context, input dto.GenerateDescriptionDTO) (string, error)

	GetPriceHistoryPoster(ctx context.Context, posterAlias string) ([]dto.PriceHistoryDTO, error)

	CreateFlatPoster(ctx context.Context, poster *dto.PosterInputFlatDTO) (*dto.CreatedPoster, error)
	UpdateFlatPoster(ctx context.Context, alias string, poster *dto.PosterInputFlatDTO) (*dto.CreatedPoster, error)
	DeleteFlatPoster(ctx context.Context, alias string, userID int) (*dto.CreatedPoster, error)

	GetPosterRoommates(ctx context.Context, alias string) ([]dto.RoommateUserDTO, error)
	AddPosterRoommate(ctx context.Context, alias string, userID int) error
}
