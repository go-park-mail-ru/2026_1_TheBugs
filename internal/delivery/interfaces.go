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
}
