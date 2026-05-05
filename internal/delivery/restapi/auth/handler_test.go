package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks/grpc_client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {}

func TestAuthHandler_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)

	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		formData       map[string]string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			formData: map[string]string{
				"email":     "test@example.com",
				"password":  "pass123",
				"phone":     "123456789",
				"firstname": "John",
				"lastname":  "Doe",
			},
			setupMock: func() {
				mockClient.EXPECT().
					RegisterUser(gomock.Any(), &auth.RegisterUserRequest{
						Email:     "test@example.com",
						Password:  "pass123",
						Phone:     "123456789",
						Firstname: "John",
						Lastname:  "Doe",
					}).
					Return(&auth.StatusResponse{Status: "registered"}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "missing fields",
			formData: map[string]string{
				"email": "test@example.com",
			},
			setupMock: func() {
				mockClient.EXPECT().
					RegisterUser(gomock.Any(), &auth.RegisterUserRequest{
						Email:     "test@example.com",
						Password:  "",
						Phone:     "",
						Firstname: "",
						Lastname:  "",
					}).
					Return(nil, status.Error(codes.InvalidArgument, "missing required fields: password, phone"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "grpc error - already exists",
			formData: map[string]string{
				"email":     "test@example.com",
				"password":  "pass123",
				"phone":     "123456789",
				"firstname": "John",
				"lastname":  "Doe",
			},
			setupMock: func() {
				mockClient.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.AlreadyExists, "user already exists"))
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			form := make(url.Values)
			for k, v := range tt.formData {
				form.Set(k, v)
			}
			req := httptest.NewRequest(http.MethodPost, "/auth/reg", bytes.NewBufferString(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			handler.RegisterUser(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
func TestAuthHandler_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)

	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		formData       map[string]string
		setupMock      func()
		expectedStatus int
		checkCookie    bool
	}{
		{
			name: "success",
			formData: map[string]string{
				"email":    "test@example.com",
				"password": "pass123",
			},
			setupMock: func() {
				mockClient.EXPECT().
					LoginUser(gomock.Any(), &auth.LoginUserRequest{
						Email:    "test@example.com",
						Password: "pass123",
					}).
					Return(&auth.LoginResponse{
						AccessToken:     "access-token",
						AccessTokenExp:  3600,
						RefreshToken:    "refresh-token",
						RefreshTokenExp: 86400,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkCookie:    true,
		},
		{
			name: "missing password",
			formData: map[string]string{
				"email":    "test@example.com",
				"password": "",
			},
			setupMock: func() {
				mockClient.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "bad credentials"))
			},
			expectedStatus: http.StatusBadRequest,
			checkCookie:    false,
		},
		{
			name: "grpc error",
			formData: map[string]string{
				"email":    "test@example.com",
				"password": "wrong",
			},
			setupMock: func() {
				mockClient.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Unauthenticated, "bad credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
			checkCookie:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			form := make(url.Values)
			for k, v := range tt.formData {
				form.Set(k, v)
			}
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			handler.LoginUser(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.checkCookie {
				cookie := rec.Result().Cookies()
				found := false
				for _, c := range cookie {
					if c.Name == "refresh_token" {
						found = true
						require.Equal(t, "refresh-token", c.Value)
						break
					}
				}
				require.True(t, found)

				var resp LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "access-token", resp.AccessToken)
				require.Equal(t, 3600, resp.AccessTokenExp)
			}
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)

	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		refreshToken   string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:         "success",
			refreshToken: "valid-refresh",
			setupMock: func() {
				mockClient.EXPECT().
					RefreshToken(gomock.Any(), &auth.RefreshTokenRequest{RefreshToken: "valid-refresh"}).
					Return(&auth.LoginResponse{
						AccessToken:     "new-access",
						AccessTokenExp:  3600,
						RefreshToken:    "new-refresh",
						RefreshTokenExp: 86400,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "grpc error",
			refreshToken: "",
			setupMock: func() {
				mockClient.EXPECT().
					RefreshToken(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "invalid token"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "bad refresh token",
			refreshToken:   "invalid",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
			if tt.refreshToken != "" {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.refreshToken})
			}
			rec := httptest.NewRecorder()

			handler.RefreshToken(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				cookie := rec.Result().Cookies()
				var found bool
				for _, c := range cookie {
					if c.Name == "refresh_token" {
						found = true
						require.Equal(t, "new-refresh", c.Value)
						break
					}
				}
				require.True(t, found)

				var resp LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "new-access", resp.AccessToken)
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)

	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		accessToken    string
		refreshToken   string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:         "success",
			accessToken:  "access-token",
			refreshToken: "refresh-token",
			setupMock: func() {
				mockClient.EXPECT().
					Logout(gomock.Any(), &auth.LogoutRequest{
						AccessToken:  "access-token",
						RefreshToken: "refresh-token",
					}).
					Return(&auth.StatusResponse{Status: "logged_out"}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing access token",
			accessToken:    "",
			refreshToken:   "refresh",
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing refresh cookie",
			accessToken:    "access",
			refreshToken:   "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tt.accessToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.accessToken)
			}
			if tt.refreshToken != "" {
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: tt.refreshToken})
			}
			rec := httptest.NewRecorder()

			handler.Logout(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusNoContent {
				cookie := rec.Result().Cookies()
				for _, c := range cookie {
					if c.Name == "refresh_token" {
						require.Equal(t, "", c.Value)
						require.Equal(t, -1, c.MaxAge)
					}
				}
			}
		})
	}
}

func TestAuthHandler_SendCodeOnEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	reqBody := `{"email":"user@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/recover", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockClient.EXPECT().
		SendCodeOnEmail(gomock.Any(), &auth.SendCodeOnEmailRequest{Email: "user@example.com"}).
		Return(&auth.SessionResponse{SessionId: "test-session"}, nil)

	handler.SendCodeOnEmail(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	cookie := rec.Result().Cookies()
	found := false
	for _, c := range cookie {
		if c.Name == "session_id" {
			found = true
			require.Equal(t, "test-session", c.Value)
			break
		}
	}
	require.True(t, found)

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "ok", resp["status"])
}

func TestAuthHandler_VKLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: map[string]string{
				"code":          "vk_code_123",
				"device_id":     "device_1",
				"state":         "state",
				"code_verifier": "verifier",
			},
			setupMock: func() {
				mockClient.EXPECT().
					VKLogin(gomock.Any(), &auth.VKLoginRequest{
						Code:         "vk_code_123",
						DeviceId:     "device_1",
						State:        "state",
						CodeVerifier: "verifier",
					}).
					Return(&auth.LoginResponse{
						AccessToken:     "access-token",
						AccessTokenExp:  3600,
						RefreshToken:    "refresh-token",
						RefreshTokenExp: 86400,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing code",
			requestBody: map[string]string{
				"code": "",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "grpc error unauthorized",
			requestBody: map[string]string{
				"code": "invalid",
			},
			setupMock: func() {
				mockClient.EXPECT().
					VKLogin(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Unauthenticated, "invalid vk code"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/vk/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.VKLogin(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "access-token", resp.AccessToken)
				require.Equal(t, 3600, resp.AccessTokenExp)

				cookies := rec.Result().Cookies()
				found := false
				for _, c := range cookies {
					if c.Name == "refresh_token" {
						found = true
						require.Equal(t, "refresh-token", c.Value)
						break
					}
				}
				require.True(t, found)
			}
		})
	}
}

func TestAuthHandler_YandexLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: map[string]string{
				"code":          "ya_code_123",
				"device_id":     "device_1",
				"state":         "state",
				"code_verifier": "verifier",
			},
			setupMock: func() {
				mockClient.EXPECT().
					YandexLogin(gomock.Any(), &auth.YandexLoginRequest{
						Code:         "ya_code_123",
						DeviceId:     "device_1",
						State:        "state",
						CodeVerifier: "verifier",
					}).
					Return(&auth.LoginResponse{
						AccessToken:     "access-yandex",
						AccessTokenExp:  7200,
						RefreshToken:    "refresh-yandex",
						RefreshTokenExp: 86400,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing code",
			requestBody: map[string]string{
				"code": "",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/yandex", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.YandexLogin(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "access-yandex", resp.AccessToken)
			}
		})
	}
}

func TestAuthHandler_SendVerifyCodeOnEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: map[string]string{
				"email": "user@example.com",
			},
			setupMock: func() {
				mockClient.EXPECT().
					SendVerifyCodeOnEmail(gomock.Any(), &auth.SendVerifyCodeOnEmailRequest{Email: "user@example.com"}).
					Return(&auth.SessionResponse{SessionId: "verify-session-123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid json",
			requestBody: map[string]string{
				"email": "",
			},
			setupMock: func() {
				mockClient.EXPECT().
					SendVerifyCodeOnEmail(gomock.Any(), &auth.SendVerifyCodeOnEmailRequest{Email: ""}).
					Return(nil, status.Error(codes.InvalidArgument, "email required"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "grpc error already verified",
			requestBody: map[string]string{
				"email": "verified@example.com",
			},
			setupMock: func() {
				mockClient.EXPECT().
					SendVerifyCodeOnEmail(gomock.Any(), &auth.SendVerifyCodeOnEmailRequest{Email: "verified@example.com"}).
					Return(nil, status.Error(codes.AlreadyExists, "email already verified"))
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/email", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.SendVerifyCodeOnEmail(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				cookies := rec.Result().Cookies()
				found := false
				for _, c := range cookies {
					if c.Name == "session_id" {
						found = true
						require.Equal(t, "verify-session-123", c.Value)
						break
					}
				}
				require.True(t, found)

				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "ok", resp["status"])
			}
		})
	}
}

func TestAuthHandler_VerifyUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		sessionID      string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:      "success",
			sessionID: "valid-session",
			requestBody: map[string]string{
				"code": "123456",
			},
			setupMock: func() {
				mockClient.EXPECT().
					VerifyCode(gomock.Any(), &auth.VerifyCodeRequest{
						SessionId: "valid-session",
						Code:      "123456",
					}).
					Return(&auth.StatusResponse{Status: "email_verified"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "missing session cookie",
			sessionID: "",
			requestBody: map[string]string{
				"code": "123456",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "invalid code",
			sessionID: "session",
			requestBody: map[string]string{
				"code": "wrong",
			},
			setupMock: func() {
				mockClient.EXPECT().
					VerifyCode(gomock.Any(), &auth.VerifyCodeRequest{SessionId: "session", Code: "wrong"}).
					Return(nil, status.Error(codes.Unauthenticated, "invalid code"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/email/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.sessionID})
			}
			rec := httptest.NewRecorder()

			handler.VerifyUserEmail(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "email verified", resp["status"])
			}
		})
	}
}

func TestAuthHandler_VerifyRecoveryCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		sessionID      string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:      "success",
			sessionID: "recovery-session",
			requestBody: map[string]string{
				"code": "654321",
			},
			setupMock: func() {
				mockClient.EXPECT().
					VerifyRecoveryCode(gomock.Any(), &auth.VerifyRecoveryCodeRequest{
						SessionId: "recovery-session",
						Code:      "654321",
					}).
					Return(&auth.StatusResponse{Status: "verified"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "missing cookie",
			sessionID: "",
			requestBody: map[string]string{
				"code": "123",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "grpc error too many attempts",
			sessionID: "session",
			requestBody: map[string]string{
				"code": "wrong",
			},
			setupMock: func() {
				mockClient.EXPECT().
					VerifyRecoveryCode(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.ResourceExhausted, "too many attempts"))
			},
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/recover/verify", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.sessionID})
			}
			rec := httptest.NewRecorder()

			handler.VerifyRecoveryCode(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "verified", resp["status"])
			}
		})
	}
}

func TestAuthHandler_UpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockAuthServiceClient(ctrl)
	handler := AuthHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		sessionID      string
		requestBody    map[string]string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:      "success",
			sessionID: "verified-session",
			requestBody: map[string]string{
				"password": "NewPass123!",
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdatePassword(gomock.Any(), &auth.UpdatePasswordRequest{
						SessionId: "verified-session",
						Password:  "NewPass123!",
					}).
					Return(&auth.StatusResponse{Status: "password_updated"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "missing cookie",
			sessionID: "",
			requestBody: map[string]string{
				"password": "newpass",
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "missing password in body",
			sessionID: "session",
			requestBody: map[string]string{
				"password": "",
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdatePassword(gomock.Any(), &auth.UpdatePasswordRequest{SessionId: "session", Password: ""}).
					Return(nil, status.Error(codes.InvalidArgument, "password required"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "session not verified",
			sessionID: "unverified",
			requestBody: map[string]string{
				"password": "newpass",
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdatePassword(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Unauthenticated, "session not verified"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/recover/reset", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.sessionID != "" {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.sessionID})
			}
			rec := httptest.NewRecorder()

			handler.UpdatePassword(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "password_updated", resp["status"])
			}
		})
	}
}
