package user

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks/grpc_client"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserHandler_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		userID         int
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			userID: 42,
			setupMock: func() {
				mockClient.EXPECT().
					GetMe(gomock.Any(), &user.GetMeRequest{UserId: 42}).
					Return(&user.GetMeResponse{
						Id:        42,
						Email:     "test@example.com",
						Phone:     "+123456789",
						AvatarUrl: lo.ToPtr("https://example.com/avatar.jpg"),
						Firstname: "John",
						Lastname:  "Doe",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user id in context",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "grpc error not found",
			userID: 99,
			setupMock: func() {
				mockClient.EXPECT().
					GetMe(gomock.Any(), &user.GetMeRequest{UserId: 99}).
					Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = utils.SetUserID(ctx, tt.userID)
			}
			req := httptest.NewRequest(http.MethodGet, "/user/me", nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetMe(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp user.GetMeResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, int32(42), resp.Id)
			}
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		userID         int
		formData       map[string]string
		fileContent    []byte
		fileName       string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success without avatar",
			userID: 10,
			formData: map[string]string{
				"first_name": "Jane",
				"last_name":  "Smith",
				"phone":      "+987654321",
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdateProfile(gomock.Any(), &user.UpdateProfileRequest{
						Id:        10,
						Firstname: lo.ToPtr("Jane"),
						Lastname:  lo.ToPtr("Smith"),
						Phone:     lo.ToPtr("+987654321"),
						File:      nil,
					}).
					Return(&user.GetMeResponse{
						Id:        10,
						Email:     "jane@example.com",
						Phone:     "+987654321",
						AvatarUrl: lo.ToPtr(""),
						Firstname: "Jane",
						Lastname:  "Smith",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success with avatar",
			userID: 10,
			formData: map[string]string{
				"first_name": "Jane",
				"last_name":  "Smith",
				"phone":      "+987654321",
			},
			fileContent: []byte("fake image data"),
			fileName:    "avatar.jpg",
			setupMock: func() {
				mockClient.EXPECT().
					UpdateProfile(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.UpdateProfileRequest) (*user.GetMeResponse, error) {
						require.NotNil(t, req.File)
						require.Equal(t, "avatar.jpg", req.File.Filename)
						require.Equal(t, "image/jpeg", req.File.ContentType) // зависит от реализации utils.ParseFileInput
						require.Equal(t, int64(len("fake image data")), req.File.Size)
						require.Equal(t, []byte("fake image data"), req.File.Avatar)
						return &user.GetMeResponse{
							Id:        10,
							Email:     "jane@example.com",
							Phone:     "+987654321",
							AvatarUrl: lo.ToPtr("https://example.com/new-avatar.jpg"),
							Firstname: "Jane",
							Lastname:  "Smith",
						}, nil
					})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user id in context",
			userID:         0,
			formData:       map[string]string{},
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "grpc error",
			userID: 10,
			formData: map[string]string{
				"first_name": "Jane",
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdateProfile(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "invalid data"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			// Создаём multipart запрос
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Добавляем текстовые поля
			for k, v := range tt.formData {
				_ = writer.WriteField(k, v)
			}

			// Добавляем файл, если задан
			if tt.fileContent != nil {
				part, err := writer.CreateFormFile("avatar", tt.fileName)
				require.NoError(t, err)
				_, err = part.Write(tt.fileContent)
				require.NoError(t, err)
			}

			err := writer.Close()
			require.NoError(t, err)

			ctx := context.Background()
			if tt.userID != 0 {
				ctx = utils.SetUserID(ctx, tt.userID)
			}
			req := httptest.NewRequest(http.MethodPut, "/user/me/profile", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.UpdateProfile(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp user.GetMeResponse
				err = json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, int32(tt.userID), resp.Id)
			}
		})
	}
}
