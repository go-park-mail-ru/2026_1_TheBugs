package auth

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/auth"
	jwtUtils "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/jwt"
)

func TestAuthServiceServer_RegisterUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "ValidPass123",
		Phone:     "+1234567890",
		Firstname: "John",
		Lastname:  "Doe",
	}

	tests := []struct {
		name      string
		req       *authpb.RegisterUserRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					RegisterUseCase(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, cred dto.CreateUserDTO) error {
						// verify DTO fields
						require.Equal(t, validReq.Email, cred.Email)
						require.Equal(t, validReq.Password, cred.Password)
						require.Equal(t, validReq.Phone, cred.Phone)
						require.Equal(t, validReq.Firstname, cred.FirstName)
						require.Equal(t, validReq.Lastname, cred.LastName)
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: &authpb.RegisterUserRequest{
				Email:    "",
				Password: "pass",
				Phone:    "phone",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing password",
			req: &authpb.RegisterUserRequest{
				Email:    "a@a.ru",
				Password: "",
				Phone:    "phone",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing phone",
			req: &authpb.RegisterUserRequest{
				Email:    "a@a.ru",
				Password: "pass",
				Phone:    "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - already exists",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					RegisterUseCase(ctx, gomock.Any()).
					Return(entity.AlredyExitError)
			},
			wantErr:  true,
			wantCode: codes.AlreadyExists, // assuming TranslateDomainsError maps to AlreadyExists
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.RegisterUser(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "registered", resp.Status)
			}
		})
	}
}

func TestAuthServiceServer_LoginUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.LoginUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedCred := &dto.UserAccessCredDTO{
		AccessToken:     "access-token",
		AccessTokenExp:  3600,
		RefreshToken:    "refresh-token",
		RefreshTokenExp: 86400,
	}

	tests := []struct {
		name      string
		req       *authpb.LoginUserRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUseCase(ctx, validReq.Email, validReq.Password).
					Return(expectedCred, nil)
			},
			wantErr: false,
		},
		{
			name: "missing email",
			req: &authpb.LoginUserRequest{
				Email:    "",
				Password: "pass",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing password",
			req: &authpb.LoginUserRequest{
				Email:    "a@a.ru",
				Password: "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUseCase(ctx, validReq.Email, validReq.Password).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - bad credentials",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUseCase(ctx, validReq.Email, validReq.Password).
					Return(nil, entity.BadCredentials)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.LoginUser(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedCred.AccessToken, resp.AccessToken)
				require.Equal(t, int64(expectedCred.AccessTokenExp), resp.AccessTokenExp)
				require.Equal(t, expectedCred.RefreshToken, resp.RefreshToken)
				require.Equal(t, int64(expectedCred.RefreshTokenExp), resp.RefreshTokenExp)
			}
		})
	}
}

func TestAuthServiceServer_RefreshToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	refreshToken := "valid-refresh-token"
	expectedCred := &dto.UserAccessCredDTO{
		AccessToken:     "new-access-token",
		AccessTokenExp:  3600,
		RefreshToken:    "new-refresh-token",
		RefreshTokenExp: 86400,
	}

	tests := []struct {
		name      string
		req       *authpb.RefreshTokenRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  &authpb.RefreshTokenRequest{RefreshToken: refreshToken},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					RefreshTokenUseCase(ctx, refreshToken).
					Return(expectedCred, nil)
			},
			wantErr: false,
		},
		{
			name:      "missing refresh token",
			req:       &authpb.RefreshTokenRequest{RefreshToken: ""},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - invalid token",
			req:  &authpb.RefreshTokenRequest{RefreshToken: refreshToken},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					RefreshTokenUseCase(ctx, refreshToken).
					Return(nil, entity.JWTError)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.RefreshToken(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedCred.AccessToken, resp.AccessToken)
				require.Equal(t, int64(expectedCred.AccessTokenExp), resp.AccessTokenExp)
				require.Equal(t, expectedCred.RefreshToken, resp.RefreshToken)
				require.Equal(t, int64(expectedCred.RefreshTokenExp), resp.RefreshTokenExp)
			}
		})
	}
}

func TestAuthServiceServer_Logout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.LogoutRequest{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	tests := []struct {
		name      string
		req       *authpb.LogoutRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LogoutUseCase(ctx, dto.LogoutDTO{
						AccessToken:  validReq.AccessToken,
						RefreshToken: validReq.RefreshToken,
					}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing access token",
			req: &authpb.LogoutRequest{
				AccessToken:  "",
				RefreshToken: "refresh",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing refresh token",
			req: &authpb.LogoutRequest{
				AccessToken:  "access",
				RefreshToken: "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LogoutUseCase(ctx, gomock.Any()).
					Return(entity.JWTError)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.Logout(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "logged_out", resp.Status)
			}
		})
	}
}

func TestAuthServiceServer_VKLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.VKLoginRequest{
		Code:         "vk-auth-code",
		DeviceId:     "device123",
		State:        "state",
		CodeVerifier: "verifier",
	}
	expectedCred := &dto.UserAccessCredDTO{
		AccessToken:     "access-token",
		AccessTokenExp:  3600,
		RefreshToken:    "refresh-token",
		RefreshTokenExp: 86400,
	}

	tests := []struct {
		name      string
		req       *authpb.VKLoginRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUserFromVKUseCase(ctx, dto.OAuthCodeFlow{
						Code:         validReq.Code,
						DeviceID:     &validReq.DeviceId,
						State:        &validReq.State,
						CodeVerifier: &validReq.CodeVerifier,
					}).
					Return(expectedCred, nil)
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: &authpb.VKLoginRequest{
				Code: "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUserFromVKUseCase(ctx, gomock.Any()).
					Return(nil, entity.BadCredentials)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.VKLogin(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedCred.AccessToken, resp.AccessToken)
			}
		})
	}
}

func TestAuthServiceServer_YandexLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.YandexLoginRequest{
		Code:         "ya-auth-code",
		DeviceId:     "device456",
		State:        "state",
		CodeVerifier: "verifier",
	}
	expectedCred := &dto.UserAccessCredDTO{
		AccessToken:     "access-token",
		AccessTokenExp:  3600,
		RefreshToken:    "refresh-token",
		RefreshTokenExp: 86400,
	}

	tests := []struct {
		name      string
		req       *authpb.YandexLoginRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUserFromYandexUseCase(ctx, dto.OAuthCodeFlow{
						Code:         validReq.Code,
						DeviceID:     &validReq.DeviceId,
						State:        &validReq.State,
						CodeVerifier: &validReq.CodeVerifier,
					}).
					Return(expectedCred, nil)
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: &authpb.YandexLoginRequest{
				Code: "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					LoginUserFromYandexUseCase(ctx, gomock.Any()).
					Return(nil, entity.BadCredentials)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.YandexLogin(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedCred.AccessToken, resp.AccessToken)
			}
		})
	}
}

func TestAuthServiceServer_SendCodeOnEmail(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "user@example.com"
	sessionID := "session-123"

	tests := []struct {
		name      string
		req       *authpb.SendCodeOnEmailRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
		wantSess  string
	}{
		{
			name: "success",
			req:  &authpb.SendCodeOnEmailRequest{Email: email},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					SendRecoveryCode(ctx, email).
					Return(sessionID, nil)
			},
			wantErr:  false,
			wantSess: sessionID,
		},
		{
			name:      "missing email",
			req:       &authpb.SendCodeOnEmailRequest{Email: ""},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - not found",
			req:  &authpb.SendCodeOnEmailRequest{Email: email},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					SendRecoveryCode(ctx, email).
					Return("", entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - too many requests",
			req:  &authpb.SendCodeOnEmailRequest{Email: email},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					SendRecoveryCode(ctx, email).
					Return("", entity.ToManyRequest)
			},
			wantErr:  true,
			wantCode: codes.ResourceExhausted,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.SendCodeOnEmail(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.wantSess, resp.SessionId)
			}
		})
	}
}

func TestAuthServiceServer_SendVerifyCodeOnEmail(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	email := "user@example.com"
	sessionID := "verify-session-123"

	tests := []struct {
		name      string
		req       *authpb.SendVerifyCodeOnEmailRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
		wantSess  string
	}{
		{
			name: "success",
			req:  &authpb.SendVerifyCodeOnEmailRequest{Email: email},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					SendVerificationEmailCode(ctx, email).
					Return(sessionID, nil)
			},
			wantErr:  false,
			wantSess: sessionID,
		},
		{
			name:      "missing email",
			req:       &authpb.SendVerifyCodeOnEmailRequest{Email: ""},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - user already verified",
			req:  &authpb.SendVerifyCodeOnEmailRequest{Email: email},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					SendVerificationEmailCode(ctx, email).
					Return("", entity.AlredyExitError) // or appropriate error
			},
			wantErr:  true,
			wantCode: codes.AlreadyExists,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.SendVerifyCodeOnEmail(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.wantSess, resp.SessionId)
			}
		})
	}
}

func TestAuthServiceServer_VerifyCode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.VerifyCodeRequest{
		SessionId: "session-123",
		Code:      "123456",
	}

	tests := []struct {
		name      string
		req       *authpb.VerifyCodeRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					VerifyUserEmail(ctx, validReq.SessionId, validReq.Code).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing session id",
			req: &authpb.VerifyCodeRequest{
				SessionId: "",
				Code:      "123",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing code",
			req: &authpb.VerifyCodeRequest{
				SessionId: "sid",
				Code:      "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - invalid code",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					VerifyUserEmail(ctx, validReq.SessionId, validReq.Code).
					Return(entity.BadCredentials)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.VerifyCode(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "email_verified", resp.Status)
			}
		})
	}
}

func TestAuthServiceServer_VerifyRecoveryCode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.VerifyRecoveryCodeRequest{
		SessionId: "recovery-session",
		Code:      "654321",
	}

	tests := []struct {
		name      string
		req       *authpb.VerifyRecoveryCodeRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					CheckRecoveryCode(ctx, validReq.SessionId, validReq.Code).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing session id",
			req: &authpb.VerifyRecoveryCodeRequest{
				SessionId: "",
				Code:      "123",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing code",
			req: &authpb.VerifyRecoveryCodeRequest{
				SessionId: "sid",
				Code:      "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - too many attempts",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					CheckRecoveryCode(ctx, validReq.SessionId, validReq.Code).
					Return(entity.ToManyRequest)
			},
			wantErr:  true,
			wantCode: codes.ResourceExhausted,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.VerifyRecoveryCode(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "verified", resp.Status)
			}
		})
	}
}

func TestAuthServiceServer_UpdatePassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validReq := &authpb.UpdatePasswordRequest{
		SessionId: "recovery-session",
		Password:  "NewStrongP@ssw0rd",
	}

	tests := []struct {
		name      string
		req       *authpb.UpdatePasswordRequest
		setupMock func(mockUC *mocks.MockAuthUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					UpdateUserPassword(ctx, validReq.SessionId, validReq.Password).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing session id",
			req: &authpb.UpdatePasswordRequest{
				SessionId: "",
				Password:  "pass",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing password",
			req: &authpb.UpdatePasswordRequest{
				SessionId: "sid",
				Password:  "",
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - bad credentials",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					UpdateUserPassword(ctx, validReq.SessionId, validReq.Password).
					Return(entity.BadCredentials)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.UpdatePassword(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "password_updated", resp.Status)
			}
		})
	}
}

func TestAuthServiceServer_CheckAccessToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	accessToken := "valid-access-token"
	// userID := 42

	tests := []struct {
		name       string
		req        *authpb.CheckAccessTokenRequest
		setupMock  func(mockUC *mocks.MockAuthUseCase)
		wantErr    bool
		wantCode   codes.Code
		wantUserID string
	}{
		{
			name: "success",
			req:  &authpb.CheckAccessTokenRequest{AccessToken: accessToken},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					ValidateAccessToken(ctx, accessToken).
					Return(&jwtUtils.Claims{Sub: "42"}, nil)
			},
			wantErr:    false,
			wantUserID: "42",
		},
		{
			name:      "missing access token",
			req:       &authpb.CheckAccessTokenRequest{AccessToken: ""},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "invalid token",
			req:  &authpb.CheckAccessTokenRequest{AccessToken: accessToken},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					ValidateAccessToken(ctx, accessToken).
					Return(nil, entity.JWTError)
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
		{
			name: "blacklisted token",
			req:  &authpb.CheckAccessTokenRequest{AccessToken: accessToken},
			setupMock: func(mockUC *mocks.MockAuthUseCase) {
				mockUC.EXPECT().
					ValidateAccessToken(ctx, accessToken).
					Return(nil, entity.NotFoundError) // assuming blacklisted maps to NotFound
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockAuthUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewAuthServiceServer(mockUC)
			resp, err := server.CheckAccessToken(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.wantUserID, resp.UserId)
			}
		})
	}
}
