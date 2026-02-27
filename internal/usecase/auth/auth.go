package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/tokens"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/validator"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo usecase.UserRepo
	authRepo usecase.AuthRepo
}

func NewAuthUseCase(userRepo usecase.UserRepo, authRepo usecase.AuthRepo) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}
func (uc AuthUseCase) RegisterUseCase(email string, password string) error {
	if !validator.ValidateEmail(email) || !validator.ValidatePwd(password) {
		return entity.InvalidInput
	}
	existing, err := uc.userRepo.GetUserByEmail(email)
	if existing != nil {
		return entity.AlredyExitError
	}
	if err != nil {
		if errors.Is(err, entity.NotFoundError) {
			return entity.NotFoundError
		}
		return fmt.Errorf("uc.userRepo.GetUserByEmail: %w", err)
	}
	salt, err := pwd.GenerateSalt()
	if err != nil {
		return fmt.Errorf("pwd.GenerateSalt: %w", err)
	}
	hashedPwd := pwd.HashPassword(password, []byte(salt))
	_, err = uc.userRepo.CreateUser(dto.CreateUserDTO{
		Email:          email,
		HashedPassword: hashedPwd,
		Salt:           salt,
	})
	if err != nil {
		return fmt.Errorf("uc.userRepo.CreateUser: %w", err)
	}
	return nil
}

func (uc AuthUseCase) LoginUseCase(email string, passwod string) (*dto.UserAccessCredDTO, error) {
	var cred dto.UserAccessCredDTO
	if !validator.ValidateEmail(email) || !validator.ValidatePwd(passwod) {
		return &cred, entity.InvalidInput
	}
	user, err := uc.userRepo.GetUserByEmail(email)

	if err != nil {
		return &cred, entity.NotFoundError
	}

	ok := pwd.VerifyPassword(passwod, []byte(user.Salt), user.HashedPassword)

	if !ok {
		return &cred, entity.BadCredentials

	}
	accessToken, err := tokens.GenerateAccessToken(user.ID, config.Config.JWT.AccessExp)

	if err != nil {
		return &cred, entity.ServiceError
	}

	refreshToken, _, err := uc.createAndSaveRefreshToken(user.ID)
	if err != nil {
		return &cred, entity.ServiceError
	}
	cred = dto.UserAccessCredDTO{
		AccessToken:     accessToken,
		AccessTokenExp:  int(config.Config.JWT.AccessExp.Seconds()),
		RefreshToken:    refreshToken,
		RefreshTokenExp: int(config.Config.JWT.RefreshExp.Seconds()),
	}
	return &cred, nil
}

func (uc AuthUseCase) createAndSaveRefreshToken(userID int) (string, time.Time, error) {
	refreshTokenID := uuid.NewString()
	refreshTokenExp := config.Config.JWT.RefreshExp
	refreshTokenExpAt := time.Now().Add(refreshTokenExp)

	refreshToken, err := tokens.GenerateRefreshToken(refreshTokenID, userID, refreshTokenExp)
	if err != nil {
		return "", time.Time{}, entity.ServiceError
	}

	uc.authRepo.CreateToken(dto.CreateRefreshTokenDTO{
		TokenID:   refreshTokenID,
		UserID:    userID,
		ExpiresAt: refreshTokenExpAt,
	})

	return refreshToken, refreshTokenExpAt, nil
}

// TODO : сделать все в одной транзакции
func (uc AuthUseCase) RefreshTokenUseCase(refreshToken string) (*dto.UserAccessCredDTO, error) {
	var cred dto.UserAccessCredDTO

	tokenData, err := tokens.ParseToken(refreshToken)
	if err != nil {
		return &cred, entity.JWTError
	}
	userID, err := strconv.Atoi(tokenData.Sub)
	if err != nil {
		return &cred, entity.JWTError
	}

	if tokenData.Type != "refresh" {
		return &cred, entity.JWTError
	}
	storedToken, err := uc.authRepo.GetToken(tokenData.ID, userID)
	if err != nil {
		return &cred, fmt.Errorf("uc.authRepo.GetToken: %w", err)
	}

	if storedToken == nil || storedToken.ExpiresAt.Before(time.Now()) {
		return &cred, entity.JWTError
	}

	uc.authRepo.DeleteToken(tokenData.ID, userID)

	accessToken, err := tokens.GenerateAccessToken(userID, config.Config.JWT.AccessExp)
	if err != nil {
		return &cred, fmt.Errorf("tokens.GenerateAccessToken: %w", err)
	}
	refreshToken, _, err = uc.createAndSaveRefreshToken(userID)
	if err != nil {
		return &cred, fmt.Errorf("uc.createAndSaveRefreshToken: %w", err)
	}
	cred = dto.UserAccessCredDTO{
		AccessToken:     accessToken,
		AccessTokenExp:  int(config.Config.JWT.AccessExp.Seconds()),
		RefreshToken:    refreshToken,
		RefreshTokenExp: int(config.Config.JWT.RefreshExp.Seconds()),
	}
	return &cred, nil
}
