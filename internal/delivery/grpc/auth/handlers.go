package main

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServiceServer contains your use‑cases and handlers.
type AuthServiceServer struct {
	uc UseCase // your existing use‑case layer
}

// RegisterUser
func (s *AuthServiceServer) RegisterUser(
	ctx context.Context,
	req *auth.RegisterUserRequest,
) (*auth.LoginResponse, error) {
	if req.Email == "" || req.Password == "" || req.Phone == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	cred := dto.CreateUserDTO{
		Email:     req.Email,
		Password:  req.Password,
		Phone:     req.Phone,
		LastName:  req.Lastname,
		FirstName: req.Firstname,
	}

	accessCred, err := s.uc.RegisterUseCase(ctx, cred)
	if err != nil {
		// map your domain error → gRPC status (e.g. AlreadyExists, InvalidArgument, etc.)
		return nil, translateError(err)
	}

	return &auth.LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp.Unix(),
	}, nil
}

// LoginUser
func (s *AuthServiceServer) LoginUser(
	ctx context.Context,
	req *auth.LoginUserRequest,
) (*auth.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	accessCred, err := s.uc.LoginUseCase(ctx, req.Email, req.Password)
	if err != nil {
		return nil, translateError(err)
	}

	return &auth.LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp.Unix(),
	}, nil
}

// RefreshToken
func (s *AuthServiceServer) RefreshToken(
	ctx context.Context,
	req *auth.RefreshTokenRequest,
) (*auth.LoginResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token required")
	}

	accessCred, err := s.uc.RefreshTokenUseCase(ctx, req.RefreshToken)
	if err != nil {
		return nil, translateError(err)
	}

	return &auth.LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp.Unix(),
	}, nil
}

// Logout
func (s *AuthServiceServer) Logout(
	ctx context.Context,
	req *auth.LogoutRequest,
) (*auth.LoginResponse, error) {
	if req.AccessToken == "" || req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token and refresh_token required")
	}

	err := s.uc.LogoutUseCase(ctx, dto.LogoutDTO{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, translateError(err)
	}

	// Return something minimal, e.g. empty response or reuse LoginResponse if needed.
	// If you want 204‑like behavior, client can just interpret empty LoginResponse as “no content”.
	return &auth.LoginResponse{}, nil
}

// SendCodeOnEmail
func (s *AuthServiceServer) SendCodeOnEmail(
	ctx context.Context,
	req *auth.SendCodeOnEmailRequest,
) (*auth.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	sessionID, err := s.uc.SendVerificationEmailCode(ctx, req.Email)
	if err != nil {
		return nil, translateError(err)
	}

	// If you need to return session via payload instead of cookies, you can do:
	// e.g. attach session_id in some field or ignore it client‑side.
	// For now just return an empty LoginResponse to match your HTTP‑style.
	return &auth.LoginResponse{}, nil
}

// VerifyCode
func (s *AuthServiceServer) VerifyCode(
	ctx context.Context,
	req *auth.VerifyCodeRequest,
) (*auth.LoginResponse, error) {
	if req.SessionId == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id and code required")
	}

	err := s.uc.VerifyUserEmail(ctx, req.SessionId, req.Code)
	if err != nil {
		return nil, translateError(err)
	}

	return &auth.LoginResponse{}, nil
}

// UpdatePassword
func (s *AuthServiceServer) UpdatePassword(
	ctx context.Context,
	req *auth.UpdatePasswordRequest,
) (*auth.LoginResponse, error) {
	if req.SessionId == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id and password required")
	}

	err := s.uc.UpdateUserPassword(ctx, req.SessionId, req.Password)
	if err != nil {
		return nil, translateError(err)
	}

	return &auth.LoginResponse{}, nil
}

// GetCSRFToken
func (s *AuthServiceServer) GetCSRFToken(
	ctx context.Context,
	req *auth.GetCSRFTokenRequest,
) (*auth.GetCSRFTokenResponse, error) {
	token := generateCSRFToken()
	return &auth.GetCSRFTokenResponse{
		CsrfToken: token,
	}, nil
}

// ------------------------------------------------------
// Helpers

func translateError(err error) error {
	// example mapping; adjust to your domain errors
	switch {
	case codes.IsAlreadyExists(err):
		return status.Error(codes.AlreadyExists, err.Error())
	case codes.IsInvalidArgument(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case codes.IsUnauthenticated(err):
		return status.Error(codes.Unauthenticated, err.Error())
	case codes.IsInternal(err):
		fallthrough
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
