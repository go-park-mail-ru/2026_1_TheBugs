package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/domains"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/oauth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/tokens"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	uow    usecase.UnitOfWork
	cache  usecase.Сache
	sender usecase.MailSender
}

func NewAuthUseCase(uow usecase.UnitOfWork, cache usecase.Сache, sender usecase.MailSender) *AuthUseCase {
	return &AuthUseCase{
		uow:    uow,
		cache:  cache,
		sender: sender,
	}
}
func (uc AuthUseCase) RegisterUseCase(ctx context.Context, data dto.CreateUserDTO) error {
	if err := validator.ValidateCred(data.Email, data.Password); err != nil {
		return err
	}
	if err := validator.ValidateProfile(data.Phone, data.FirstName, data.LastName); err != nil {
		return err
	}
	existing, err := uc.uow.Users().GetByEmail(ctx, data.Email)
	if existing != nil {
		return entity.AlredyExitError
	}
	if err != nil {
		if !errors.Is(err, entity.NotFoundError) {
			return fmt.Errorf("uc.userRepo.GetByEmail: %w", err)
		}

	}
	salt, err := pwd.GenerateSalt()
	if err != nil {
		return fmt.Errorf("pwd.GenerateSalt: %w", err)
	}
	hashedPwd := pwd.HashPassword(data.Password, []byte(salt))
	_, err = uc.uow.Users().Create(ctx, dto.CreateUserDTO{
		Email:          data.Email,
		HashedPassword: &hashedPwd,
		Salt:           &salt,
		FirstName:      data.FirstName,
		LastName:       data.LastName,
		Phone:          validator.NormolizePhoneNumber(data.Phone),
	})
	if err != nil {
		return fmt.Errorf("uc.userRepo.Create: %w", err)
	}
	return nil
}

func (uc AuthUseCase) LoginUseCase(ctx context.Context, email string, passwod string) (*dto.UserAccessCredDTO, error) {
	var cred dto.UserAccessCredDTO
	if err := validator.ValidateCred(email, passwod); err != nil {
		return &cred, err
	}
	user, err := uc.uow.Users().GetByEmail(ctx, email)

	if err != nil {
		return &cred, entity.NotFoundError
	}
	if user.HashedPassword == nil || user.Salt == nil {
		return nil, entity.BadCredentials
	}

	ok := pwd.VerifyPassword(passwod, []byte(*user.Salt), *user.HashedPassword)

	if !ok {
		return &cred, entity.BadCredentials

	}
	accessToken, err := tokens.GenerateAccessToken(user.ID, config.Config.JWT.AccessExp)

	if err != nil {
		return &cred, entity.ServiceError
	}

	refreshToken, _, err := uc.createAndSaveRefreshToken(ctx, user.ID)
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

func (uc AuthUseCase) createAndSaveRefreshToken(ctx context.Context, userID int) (string, time.Time, error) {
	refreshTokenID := uuid.NewString()
	refreshTokenExp := config.Config.JWT.RefreshExp
	refreshTokenExpAt := time.Now().Add(refreshTokenExp)

	refreshToken, err := tokens.GenerateRefreshToken(refreshTokenID, userID, refreshTokenExp)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokens.GenerateRefreshToken: %w", entity.ServiceError)
	}

	err = uc.uow.Autho().CreateToken(ctx, dto.CreateRefreshTokenDTO{
		TokenID:   refreshTokenID,
		UserID:    userID,
		ExpiresAt: refreshTokenExpAt,
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("uc.authRepo.CreateToken: %w", err)
	}

	return refreshToken, refreshTokenExpAt, nil
}

func (uc AuthUseCase) RefreshTokenUseCase(ctx context.Context, refreshToken string) (*dto.UserAccessCredDTO, error) {
	var cred dto.UserAccessCredDTO

	tokenData, userID, err := uc.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return &cred, fmt.Errorf("uc.ValidateRefreshToken: %w", err)
	}
	accessToken, err := tokens.GenerateAccessToken(userID, config.Config.JWT.AccessExp)
	if err != nil {
		return &cred, fmt.Errorf("tokens.GenerateAccessToken: %w", err)
	}
	refreshTokenID := uuid.NewString()
	refreshTokenExp := config.Config.JWT.RefreshExp
	refreshTokenExpAt := time.Now().Add(refreshTokenExp)

	refreshToken, err = tokens.GenerateRefreshToken(refreshTokenID, userID, refreshTokenExp)
	if err != nil {
		return &cred, fmt.Errorf("tokens.GenerateRefreshToken: %w", entity.ServiceError)
	}
	err = uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {
		err = uc.uow.Autho().DeleteToken(ctx, tokenData.ID, userID)
		if err != nil {
			return fmt.Errorf("uc.authRepo.DeleteToken: %w", err)
		}
		err = uc.uow.Autho().CreateToken(ctx, dto.CreateRefreshTokenDTO{
			TokenID:   refreshTokenID,
			UserID:    userID,
			ExpiresAt: refreshTokenExpAt,
		})
		if err != nil {
			return fmt.Errorf("uc.authRepo.CreateToken: %w", err)
		}
		return nil
	})
	if err != nil {
		return &cred, fmt.Errorf("uc.uow.Do: %w", err)
	}

	cred = dto.UserAccessCredDTO{
		AccessToken:     accessToken,
		AccessTokenExp:  int(config.Config.JWT.AccessExp.Seconds()),
		RefreshToken:    refreshToken,
		RefreshTokenExp: int(config.Config.JWT.RefreshExp.Seconds()),
	}
	return &cred, nil
}
func (uc AuthUseCase) ValidateAccessToken(ctx context.Context, accessToken string) (*tokens.Claims, error) {
	accessData, err := tokens.ParseToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("tokens.ParseToken: %w", entity.JWTError)
	}
	if accessData.Type != entity.AccessTokenType {
		return nil, entity.JWTError
	}
	if ok, err := uc.cache.IsBlacklisted(ctx, accessData.ID); err != nil || ok {
		return nil, fmt.Errorf("uc.authRepo.IsBlacklisted: %w", entity.JWTError)
	}

	return accessData, nil
}

func (uc AuthUseCase) ValidateRefreshToken(ctx context.Context, refreshToken string) (*tokens.Claims, int, error) {
	refreshData, err := tokens.ParseToken(refreshToken)
	if err != nil {
		return nil, 0, fmt.Errorf("tokens.ParseToken: %w", entity.JWTError)
	}
	if refreshData.Type != entity.RefreshTokenType {
		return nil, 0, entity.JWTError
	}
	userID, err := strconv.Atoi(refreshData.Sub)
	if err != nil {
		return nil, 0, fmt.Errorf("strconv.Atoi: %w", entity.JWTError)
	}
	storedToken, err := uc.uow.Autho().GetToken(ctx, refreshData.ID, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("uc.authRepo.GetToken: %w", entity.JWTError)
	}

	if storedToken == nil || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, 0, entity.JWTError
	}
	return refreshData, userID, nil
}

func (uc AuthUseCase) LogoutUseCase(ctx context.Context, logoutCred dto.LogoutDTO) error {

	accessData, err := uc.ValidateAccessToken(ctx, logoutCred.AccessToken)
	if err != nil {
		return fmt.Errorf("uc.ValidateAccessToken: %w", err)
	}
	refreshData, userID, err := uc.ValidateRefreshToken(ctx, logoutCred.RefreshToken)
	if err != nil {
		return fmt.Errorf("uc.ValidateRefreshToken: %w", err)
	}

	err = uc.uow.Autho().DeleteToken(ctx, refreshData.ID, userID)
	if err != nil {
		return fmt.Errorf("uc.authRepo.DeleteToken: %w", entity.JWTError)
	}
	ttl := config.Config.JWT.AccessExp

	if err := uc.cache.SetBlacklist(ctx, accessData.ID, ttl); err != nil {
		return fmt.Errorf("uc.authRepo.BlacklistToken: %w", entity.JWTError)
	}
	return nil
}

func (uc AuthUseCase) LoginUserFromVKUseCase(ctx context.Context, flow dto.OAuthCodeFlow) (*dto.UserAccessCredDTO, error) {
	vkCred, err := oauth.ChangeCodeToAccessToken(ctx, flow)
	if err != nil {
		return nil, fmt.Errorf("uc.oauthRepo.ChangeCodeToCred: %w", err)
	}
	data, err := oauth.GetUserPublicInfoVK(ctx, vkCred.IDToken)
	if err != nil {
		return nil, fmt.Errorf("oauth.GetUserPublicInfoVK: %w", err)
	}
	log.Printf("VK claims: sub=%s, email=%s, name=%s", data.User.UserID, data.User.Email, data.User.FirstName)
	user, err := uc.uow.Users().GetByProvider(ctx, data.User.Email, "vk")
	if err != nil {
		if !errors.Is(err, entity.NotFoundError) {
			return nil, err
		} else {
			user, err = uc.uow.Users().CreateByProvider(ctx, dto.CreateUserByProviderDTO{
				Provider: "vk",
				Email:    data.User.Email,
			})
			if err != nil {
				return nil, fmt.Errorf("uc.userRepo.CreateByProvider: %w", err)
			}
		}
	}
	accessToken, err := tokens.GenerateAccessToken(user.ID, config.Config.JWT.AccessExp)

	if err != nil {
		return nil, entity.ServiceError
	}

	refreshToken, _, err := uc.createAndSaveRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, entity.ServiceError
	}
	cred := dto.UserAccessCredDTO{
		AccessToken:     accessToken,
		AccessTokenExp:  int(config.Config.JWT.AccessExp.Seconds()),
		RefreshToken:    refreshToken,
		RefreshTokenExp: int(config.Config.JWT.RefreshExp.Seconds()),
	}
	return &cred, nil
}

func (uc AuthUseCase) LoginUserFromYandexUseCase(ctx context.Context, flow dto.OAuthCodeFlow) (*dto.UserAccessCredDTO, error) {
	yandexCred, err := oauth.ChangeYandexCodeToAccessToken(ctx, flow)
	if err != nil {
		return nil, fmt.Errorf("uc.oauthRepo.ChangeCodeToCred: %w", err)
	}
	data, err := oauth.GetYandexUserPublicInfo(ctx, yandexCred.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("oauth.GetYandexUserPublicInfo: %w", err)
	}
	log.Printf("Yandex claims: sub=%s, email=%s, name=%s", data.ID, data.DefaultEmail, data.FirstName)
	user, err := uc.uow.Users().GetByEmail(ctx, data.DefaultEmail)
	if err != nil {
		if !errors.Is(err, entity.NotFoundError) {
			return nil, err
		} else {
			user, err = uc.uow.Users().CreateByProvider(ctx, dto.CreateUserByProviderDTO{
				Provider:   "yandex",
				Email:      data.DefaultEmail,
				Phone:      validator.NormolizePhoneNumber(data.DefaultPhone.Number),
				LastName:   data.LastName,
				FirstName:  data.FirstName,
				ProviderID: &data.ID,
			})
			if err != nil {
				return nil, fmt.Errorf("uc.userRepo.CreateByProvider: %w", err)
			}
		}
	}
	accessToken, err := tokens.GenerateAccessToken(user.ID, config.Config.JWT.AccessExp)

	if err != nil {
		return nil, entity.ServiceError
	}

	refreshToken, _, err := uc.createAndSaveRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, entity.ServiceError
	}
	cred := dto.UserAccessCredDTO{
		AccessToken:     accessToken,
		AccessTokenExp:  int(config.Config.JWT.AccessExp.Seconds()),
		RefreshToken:    refreshToken,
		RefreshTokenExp: int(config.Config.JWT.RefreshExp.Seconds()),
	}
	return &cred, nil
}

const MaxAttemptsRecovery = 5

func (uc AuthUseCase) SendVerificationCode(ctx context.Context, email string) (string, error) {
	op := "AuthUseCase.SendVerificationCode"
	log := middleware.GetLogger(ctx).WithField("op", op)

	_, err := uc.uow.Users().GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("uc.uow.Users().GetByEmail: %w", err)
	}
	blocked, err := uc.cache.IsBlacklisted(ctx, email)
	if err != nil {
		return "", fmt.Errorf("uc.cache.IsBlacklisted: %w", err)
	}
	if blocked {
		return "", fmt.Errorf("max limit offset: %w", entity.ToManyRequest)
	}
	err = uc.cache.SetBlacklist(ctx, email, time.Duration(1*time.Minute))
	if err != nil {
		return "", fmt.Errorf("uc.cache.SetBlacklist: %w", err)
	}
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%05d", seed.Intn(100000))
	sessionId := fmt.Sprintf("%010x", seed.Int())[:10]
	err = uc.cache.CreateRecoverSession(ctx, sessionId, domains.RecoverSession{Email: email, Code: code, Attempts: 0, Verified: false}, config.Config.JWT.RecoverExp)
	if err != nil {
		return "", fmt.Errorf("uc.cache.CreateRecoverSession: %w", err)
	}
	log.Info("send code")
	if err := uc.sender.SendCode(ctx, email, code); err != nil {
		log.Errorf("send code: %v", err)
	}

	return sessionId, nil
}

func (uc AuthUseCase) CheckRecoveryCode(ctx context.Context, sessionID string, code string) error {
	session, err := uc.cache.GetRecoverSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("uc.cache.GetRecoverSession: %w", err)
	}
	if session == nil {
		return fmt.Errorf("session is empty: %w", entity.BadCredentials)
	}

	if session.Code != code {
		attempts, _ := uc.cache.IncrementRecoverAttempts(ctx, sessionID)
		if attempts > MaxAttemptsRecovery {
			_ = uc.cache.DeleteRecoverSession(ctx, sessionID)
			return fmt.Errorf("max limit offset: %w", entity.ToManyRequest)
		}
		return fmt.Errorf("bad code: %w", entity.BadCredentials)
	}
	err = uc.cache.SetRecoverVerified(ctx, sessionID, true)
	if err != nil {
		return fmt.Errorf("uc.cache.SetRecoverVerified: %w", err)
	}
	return nil
}

func (uc AuthUseCase) UpdateUserPassword(ctx context.Context, sessionID string, password string) error {

	session, err := uc.cache.GetRecoverSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("uc.cache.GetRecoverSession: %w", err)
	}
	if session == nil {
		return fmt.Errorf("session empty: %w", entity.BadCredentials)
	}
	if !session.Verified {
		return fmt.Errorf("session unvirified: %w", entity.BadCredentials)
	}
	salt, _ := pwd.GenerateSalt()
	hashedPwd := pwd.HashPassword(password, []byte(salt))
	err = uc.uow.Users().UpdatePwd(ctx, session.Email, hashedPwd, salt)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().UpdatePwd: %w", err)
	}
	err = uc.cache.DeleteRecoverSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("uc.cache.DeleteRecoverSession: %w", err)
	}
	return nil
}
