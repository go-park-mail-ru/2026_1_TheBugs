package auth

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	jwtUtils "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/jwt"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UseCase defines the interface for auth use cases
type UseCase interface {
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

type AuthServiceServer struct {
	auth.UnimplementedAuthServiceServer
	uc UseCase
}

// NewAuthServiceServer creates a new AuthServiceServer
func NewAuthServiceServer(uc UseCase) *AuthServiceServer {
	return &AuthServiceServer{
		uc: uc,
	}
}

// RegisterUser handles user registration
func (s *AuthServiceServer) RegisterUser(
	ctx context.Context,
	req *auth.RegisterUserRequest,
) (*auth.StatusResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "RegisterUser")
	if req.Email == "" || req.Password == "" || req.Phone == "" {
		log.Error("missing required fields: email, password, phone")
		return nil, status.Error(codes.InvalidArgument, "missing required fields: email, password, phone")
	}

	cred := dto.CreateUserDTO{
		Email:     req.Email,
		Password:  req.Password,
		Phone:     req.Phone,
		LastName:  req.Lastname,
		FirstName: req.Firstname,
	}

	err := s.uc.RegisterUseCase(ctx, cred)
	if err != nil {
		log.Errorf("s.uc.RegisterUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.StatusResponse{Status: "registered"}, nil
}

// LoginUser handles user login
func (s *AuthServiceServer) LoginUser(
	ctx context.Context,
	req *auth.LoginUserRequest,
) (*auth.LoginResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "LoginUser")
	if req.Email == "" || req.Password == "" {
		log.Error("missing required fields: email, password")
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	accessCred, err := s.uc.LoginUseCase(ctx, req.Email, req.Password)
	if err != nil {
		log.Errorf("s.uc.LoginUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.LoginResponse{
		AccessToken:     accessCred.AccessToken,
		AccessTokenExp:  int64(accessCred.AccessTokenExp),
		RefreshToken:    accessCred.RefreshToken,
		RefreshTokenExp: int64(accessCred.RefreshTokenExp),
	}, nil
}

// RefreshToken handles token refresh
func (s *AuthServiceServer) RefreshToken(
	ctx context.Context,
	req *auth.RefreshTokenRequest,
) (*auth.LoginResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "RefreshToken")
	if req.RefreshToken == "" {
		log.Error("missing required field: refresh_token")
		return nil, status.Error(codes.InvalidArgument, "refresh_token required")
	}

	accessCred, err := s.uc.RefreshTokenUseCase(ctx, req.RefreshToken)
	if err != nil {
		log.Errorf("s.uc.RefreshTokenUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.LoginResponse{
		AccessToken:     accessCred.AccessToken,
		AccessTokenExp:  int64(accessCred.AccessTokenExp),
		RefreshToken:    accessCred.RefreshToken,
		RefreshTokenExp: int64(accessCred.RefreshTokenExp),
	}, nil
}

// Logout handles user logout
func (s *AuthServiceServer) Logout(
	ctx context.Context,
	req *auth.LogoutRequest,
) (*auth.StatusResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "Logout")
	if req.AccessToken == "" || req.RefreshToken == "" {
		log.Error("missing required fields: access_token and refresh_token")
		return nil, status.Error(codes.InvalidArgument, "access_token and refresh_token required")
	}

	err := s.uc.LogoutUseCase(ctx, dto.LogoutDTO{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		log.Errorf("s.uc.LogoutUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.StatusResponse{Status: "logged_out"}, nil
}

// VKLogin handles VK OAuth login
func (s *AuthServiceServer) VKLogin(
	ctx context.Context,
	req *auth.VKLoginRequest,
) (*auth.LoginResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "VKLogin")
	if req.Code == "" {
		log.Error("missing required field: code")
		return nil, status.Error(codes.InvalidArgument, "code required")
	}

	flow := dto.OAuthCodeFlow{
		Code:         req.Code,
		DeviceID:     &req.DeviceId,
		State:        &req.State,
		CodeVerifier: &req.CodeVerifier,
	}

	accessCred, err := s.uc.LoginUserFromVKUseCase(ctx, flow)
	if err != nil {
		log.Errorf("s.uc.LoginUserFromVKUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.LoginResponse{
		AccessToken:     accessCred.AccessToken,
		AccessTokenExp:  int64(accessCred.AccessTokenExp),
		RefreshToken:    accessCred.RefreshToken,
		RefreshTokenExp: int64(accessCred.RefreshTokenExp),
	}, nil
}

// YandexLogin handles Yandex OAuth login
func (s *AuthServiceServer) YandexLogin(
	ctx context.Context,
	req *auth.YandexLoginRequest,
) (*auth.LoginResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "YandexLogin")
	if req.Code == "" {
		log.Error("missing required field: code")
		return nil, status.Error(codes.InvalidArgument, "code required")
	}

	flow := dto.OAuthCodeFlow{
		Code:         req.Code,
		DeviceID:     &req.DeviceId,
		State:        &req.State,
		CodeVerifier: &req.CodeVerifier,
	}

	accessCred, err := s.uc.LoginUserFromYandexUseCase(ctx, flow)
	if err != nil {
		log.Errorf("s.uc.LoginUserFromYandexUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.LoginResponse{
		AccessToken:     accessCred.AccessToken,
		AccessTokenExp:  int64(accessCred.AccessTokenExp),
		RefreshToken:    accessCred.RefreshToken,
		RefreshTokenExp: int64(accessCred.RefreshTokenExp),
	}, nil
}

// SendCodeOnEmail sends recovery code to email
func (s *AuthServiceServer) SendCodeOnEmail(
	ctx context.Context,
	req *auth.SendCodeOnEmailRequest,
) (*auth.SessionResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendCodeOnEmail")
	if req.Email == "" {
		log.Error("missing required field: email")
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	sessionID, err := s.uc.SendRecoveryCode(ctx, req.Email)
	if err != nil {
		log.Errorf("s.uc.SendRecoveryCode: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.SessionResponse{SessionId: sessionID}, nil
}

// SendVerifyCodeOnEmail sends verification code to email
func (s *AuthServiceServer) SendVerifyCodeOnEmail(
	ctx context.Context,
	req *auth.SendVerifyCodeOnEmailRequest,
) (*auth.SessionResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "SendVerifyCodeOnEmail")
	if req.Email == "" {
		log.Error("missing required field: email")
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	sessionID, err := s.uc.SendVerificationEmailCode(ctx, req.Email)
	if err != nil {
		log.Errorf("s.uc.SendVerificationEmailCode: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.SessionResponse{SessionId: sessionID}, nil
}

// VerifyCode verifies email verification code
func (s *AuthServiceServer) VerifyCode(
	ctx context.Context,
	req *auth.VerifyCodeRequest,
) (*auth.StatusResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "VerifyCode")
	if req.SessionId == "" || req.Code == "" {
		log.Error("missing required fields: session_id and code")
		return nil, status.Error(codes.InvalidArgument, "session_id and code required")
	}

	err := s.uc.VerifyUserEmail(ctx, req.SessionId, req.Code)
	if err != nil {
		log.Errorf("s.uc.VerifyUserEmail: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.StatusResponse{Status: "email_verified"}, nil
}

// VerifyRecoveryCode verifies recovery code
func (s *AuthServiceServer) VerifyRecoveryCode(
	ctx context.Context,
	req *auth.VerifyRecoveryCodeRequest,
) (*auth.StatusResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "VerifyRecoveryCode")
	if req.SessionId == "" || req.Code == "" {
		log.Error("missing required fields: session_id and code")
		return nil, status.Error(codes.InvalidArgument, "session_id and code required")
	}

	err := s.uc.CheckRecoveryCode(ctx, req.SessionId, req.Code)
	if err != nil {
		log.Errorf("s.uc.CheckRecoveryCode: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.StatusResponse{Status: "verified"}, nil
}

// UpdatePassword updates user password
func (s *AuthServiceServer) UpdatePassword(
	ctx context.Context,
	req *auth.UpdatePasswordRequest,
) (*auth.StatusResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "UpdatePassword")
	if req.SessionId == "" || req.Password == "" {
		log.Error("missing required fields: session_id and password")
		return nil, status.Error(codes.InvalidArgument, "session_id and password required")
	}

	err := s.uc.UpdateUserPassword(ctx, req.SessionId, req.Password)
	if err != nil {
		log.Errorf("s.uc.UpdateUserPassword: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.StatusResponse{Status: "password_updated"}, nil
}

func (s *AuthServiceServer) CheckAccessToken(
	ctx context.Context,
	req *auth.CheckAccessTokenRequest,
) (*auth.CheckAccessTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token required")
	}

	claims, err := s.uc.ValidateAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &auth.CheckAccessTokenResponse{UserId: claims.Sub}, nil
}
