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
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
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
					UpdateProfile(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.UpdateProfileRequest, opt ...any) (*user.GetMeResponse, error) {
						require.NotNil(t, req.File)
						require.Equal(t, "avatar.jpg", req.File.Filename)
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

func TestUserHandler_GetRoommateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		pathID         string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			pathID: "42",
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateUserRequest, opts ...grpc.CallOption) (*user.GetRoommateUserResponse, error) {
						require.Equal(t, int64(42), req.UserId)

						return &user.GetRoommateUserResponse{
							FirstName:   "John",
							LastName:    "Doe",
							AvatarUrl:   lo.ToPtr("https://example.com/avatar.jpg"),
							Gender:      "male",
							Birthday:    "2000-01-01",
							Description: lo.ToPtr("good roommate"),
							Tags: []*user.RoommateTag{
								{
									Name:  "Не курю",
									Alias: "no_smoking",
								},
							},
						}, nil
					})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid user id",
			pathID:         "bad",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc error not found",
			pathID: "99",
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateUserRequest, opts ...grpc.CallOption) (*user.GetRoommateUserResponse, error) {
						require.Equal(t, int64(99), req.UserId)
						return nil, status.Error(codes.NotFound, "user not found")
					})
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc error internal",
			pathID: "42",
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateUserRequest, opts ...grpc.CallOption) (*user.GetRoommateUserResponse, error) {
						require.Equal(t, int64(42), req.UserId)
						return nil, status.Error(codes.Internal, "internal error")
					})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/users/"+test.pathID, nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": test.pathID,
			})
			rec := httptest.NewRecorder()

			handler.GetRoommateUser(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.expectedStatus == http.StatusOK {
				var resp user.GetRoommateUserResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, "John", resp.FirstName)
				require.Equal(t, "Doe", resp.LastName)
				require.Equal(t, "male", resp.Gender)
				require.Equal(t, "2000-01-01", resp.Birthday)
				require.NotNil(t, resp.AvatarUrl)
				require.Equal(t, "https://example.com/avatar.jpg", *resp.AvatarUrl)
				require.NotNil(t, resp.Description)
				require.Equal(t, "good roommate", *resp.Description)
				require.Len(t, resp.Tags, 1)
				require.Equal(t, "Не курю", resp.Tags[0].Name)
				require.Equal(t, "no_smoking", resp.Tags[0].Alias)
			}
		})
	}
}

func TestUserHandler_AddRoommateMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		pathID         string
		userID         int
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			pathID: "42",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddRoommateMatch(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.AddRoommateMatchRequest, opts ...grpc.CallOption) (*user.AddRoommateMatchResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(42), req.ToUserId)
						return &user.AddRoommateMatchResponse{}, nil
					})
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing user id in context",
			pathID:         "42",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid target user id",
			pathID:         "bad",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc error not found",
			pathID: "99",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddRoommateMatch(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.AddRoommateMatchRequest, opts ...grpc.CallOption) (*user.AddRoommateMatchResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(99), req.ToUserId)
						return nil, status.Error(codes.NotFound, "user not found")
					})
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc error invalid argument",
			pathID: "10",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddRoommateMatch(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.AddRoommateMatchRequest, opts ...grpc.CallOption) (*user.AddRoommateMatchResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(10), req.ToUserId)
						return nil, status.Error(codes.InvalidArgument, "invalid input")
					})
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc error internal",
			pathID: "42",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddRoommateMatch(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.AddRoommateMatchRequest, opts ...grpc.CallOption) (*user.AddRoommateMatchResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(42), req.ToUserId)
						return nil, status.Error(codes.Internal, "internal error")
					})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodPost, "/users/"+test.pathID+"/match", nil).WithContext(ctx)
			req = mux.SetURLVars(req, map[string]string{
				"id": test.pathID,
			})
			rec := httptest.NewRecorder()

			handler.AddRoommateMatch(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_GetRoommateContacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		pathID         string
		userID         int
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			pathID: "42",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateContacts(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateContactsRequest, opts ...grpc.CallOption) (*user.GetRoommateContactsResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(42), req.ToUserId)

						return &user.GetRoommateContactsResponse{
							Email: "target@example.com",
							Phone: "+79991234567",
						}, nil
					})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user id in context",
			pathID:         "42",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid target user id",
			pathID:         "bad",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc error not found",
			pathID: "99",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateContacts(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateContactsRequest, opts ...grpc.CallOption) (*user.GetRoommateContactsResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(99), req.ToUserId)
						return nil, status.Error(codes.NotFound, "contacts not found")
					})
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc error invalid argument",
			pathID: "10",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateContacts(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateContactsRequest, opts ...grpc.CallOption) (*user.GetRoommateContactsResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(10), req.ToUserId)
						return nil, status.Error(codes.InvalidArgument, "invalid input")
					})
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc error internal",
			pathID: "42",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateContacts(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateContactsRequest, opts ...grpc.CallOption) (*user.GetRoommateContactsResponse, error) {
						require.Equal(t, int64(10), req.FromUserId)
						require.Equal(t, int64(42), req.ToUserId)
						return nil, status.Error(codes.Internal, "internal error")
					})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodGet, "/user/"+test.pathID+"/contacts", nil).WithContext(ctx)
			req = mux.SetURLVars(req, map[string]string{
				"id": test.pathID,
			})
			rec := httptest.NewRecorder()

			handler.GetRoommateContacts(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.expectedStatus == http.StatusOK {
				var resp user.GetRoommateContactsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, "target@example.com", resp.Email)
				require.Equal(t, "+79991234567", resp.Phone)
			}
		})
	}
}

func TestUserHandler_CreateRoommateForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		userID         int
		body           any
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			userID: 10,
			body: dto.CreateRoommateFormRequest{
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking", "no_pets"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					CreateRoommateForm(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.CreateRoommateFormRequest, opts ...grpc.CallOption) (*user.CreateRoommateFormResponse, error) {
						require.Equal(t, int64(10), req.UserId)
						require.Equal(t, "male", req.Gender)
						require.Equal(t, "2000-01-01", req.Birthday)
						require.Equal(t, "good roommate", req.Description)
						require.Equal(t, []string{"no_smoking", "no_pets"}, req.Tags)

						return &user.CreateRoommateFormResponse{}, nil
					})
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing user id in context",
			userID:         0,
			body:           dto.CreateRoommateFormRequest{},
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid body",
			userID:         10,
			body:           "{bad json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc invalid argument",
			userID: 10,
			body: dto.CreateRoommateFormRequest{
				Gender:      "",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					CreateRoommateForm(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "invalid form"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc internal error",
			userID: 10,
			body: dto.CreateRoommateFormRequest{
				Gender:      "female",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					CreateRoommateForm(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			var body *bytes.Buffer
			switch v := test.body.(type) {
			case string:
				body = bytes.NewBufferString(v)
			default:
				body = &bytes.Buffer{}
				err := json.NewEncoder(body).Encode(v)
				require.NoError(t, err)
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodPost, "/users/me/roommate-form", body).WithContext(ctx)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateRoommateForm(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_GetRoommateForm(t *testing.T) {
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
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateForm(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateFormRequest, opts ...grpc.CallOption) (*user.GetRoommateFormResponse, error) {
						require.Equal(t, int64(10), req.UserId)

						return &user.GetRoommateFormResponse{
							Gender:      "male",
							Birthday:    "2000-01-01",
							Description: "good roommate",
							Tags:        []string{"no_smoking", "no_pets"},
						}, nil
					})
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
			name:   "grpc not found",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateForm(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateFormRequest, opts ...grpc.CallOption) (*user.GetRoommateFormResponse, error) {
						require.Equal(t, int64(10), req.UserId)
						return nil, status.Error(codes.NotFound, "form not found")
					})
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetRoommateForm(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetRoommateFormRequest, opts ...grpc.CallOption) (*user.GetRoommateFormResponse, error) {
						require.Equal(t, int64(10), req.UserId)
						return nil, status.Error(codes.Internal, "internal error")
					})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodGet, "/users/me/roommate-form", nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetRoommateForm(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.expectedStatus == http.StatusOK {
				var resp user.GetRoommateFormResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, "male", resp.Gender)
				require.Equal(t, "2000-01-01", resp.Birthday)
				require.Equal(t, "good roommate", resp.Description)
				require.Equal(t, []string{"no_smoking", "no_pets"}, resp.Tags)
			}
		})
	}
}

func TestUserHandler_UpdateRoommateForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockUserServiceClient(ctrl)
	handler := &UserHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		userID         int
		body           any
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			userID: 10,
			body: dto.RoommateFormDTO{
				Gender:      "female",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking", "no_pets"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdateRoommateForm(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.UpdateRoommateFormRequest, opts ...grpc.CallOption) (*user.UpdateRoommateFormResponse, error) {
						require.Equal(t, int64(10), req.UserId)
						require.Equal(t, "female", req.Gender)
						require.Equal(t, "2001-02-03", req.Birthday)
						require.Equal(t, "updated roommate", req.Description)
						require.Equal(t, []string{"no_smoking", "no_pets"}, req.Tags)

						return &user.UpdateRoommateFormResponse{}, nil
					})
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing user id in context",
			userID:         0,
			body:           dto.RoommateFormDTO{},
			setupMock:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid body",
			userID:         10,
			body:           "{bad json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc invalid argument",
			userID: 10,
			body: dto.RoommateFormDTO{
				Gender:      "",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdateRoommateForm(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "invalid form"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc not found",
			userID: 10,
			body: dto.RoommateFormDTO{
				Gender:      "male",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdateRoommateForm(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "form not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			userID: 10,
			body: dto.RoommateFormDTO{
				Gender:      "female",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func() {
				mockClient.EXPECT().
					UpdateRoommateForm(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			var body *bytes.Buffer
			switch v := test.body.(type) {
			case string:
				body = bytes.NewBufferString(v)
			default:
				body = &bytes.Buffer{}
				err := json.NewEncoder(body).Encode(v)
				require.NoError(t, err)
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodPut, "/users/me/roommate-form", body).WithContext(ctx)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.UpdateRoommateForm(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestUserHandler_GetIncomingRoommateMatches(t *testing.T) {
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
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetIncomingRoommateMatches(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetIncomingRoommateMatchesRequest, opts ...grpc.CallOption) (*user.GetIncomingRoommateMatchesResponse, error) {
						require.Equal(t, int64(10), req.UserId)

						return &user.GetIncomingRoommateMatchesResponse{
							Users: []*user.RoommateUser{
								{
									Id:          42,
									FirstName:   "John",
									LastName:    "Doe",
									AvatarUrl:   lo.ToPtr("https://example.com/avatar.jpg"),
									PosterAlias: lo.ToPtr("flat-1"),
								},
							},
							Len: 1,
						}, nil
					})
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
			name:   "grpc error internal",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetIncomingRoommateMatches(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetIncomingRoommateMatchesRequest, opts ...grpc.CallOption) (*user.GetIncomingRoommateMatchesResponse, error) {
						require.Equal(t, int64(10), req.UserId)
						return nil, status.Error(codes.Internal, "internal error")
					})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodGet, "/user/me/roommate-matches/incoming", nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetIncomingRoommateMatches(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.expectedStatus == http.StatusOK {
				var resp user.GetIncomingRoommateMatchesResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, int64(1), resp.Len)
				require.Len(t, resp.Users, 1)
				require.Equal(t, int64(42), resp.Users[0].Id)
				require.Equal(t, "John", resp.Users[0].FirstName)
				require.Equal(t, "Doe", resp.Users[0].LastName)
				require.NotNil(t, resp.Users[0].AvatarUrl)
				require.Equal(t, "https://example.com/avatar.jpg", *resp.Users[0].AvatarUrl)
				require.NotNil(t, resp.Users[0].PosterAlias)
				require.Equal(t, "flat-1", *resp.Users[0].PosterAlias)
			}
		})
	}
}

func TestUserHandler_GetMatchedRoommateMatches(t *testing.T) {
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
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetMatchedRoommateMatches(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetMatchedRoommateMatchesRequest, opts ...grpc.CallOption) (*user.GetMatchedRoommateMatchesResponse, error) {
						require.Equal(t, int64(10), req.UserId)

						return &user.GetMatchedRoommateMatchesResponse{
							Users: []*user.RoommateUser{
								{
									Id:        42,
									FirstName: "John",
									LastName:  "Doe",
									AvatarUrl: lo.ToPtr("https://example.com/avatar.jpg"),
								},
							},
							Len: 1,
						}, nil
					})
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
			name:   "grpc error internal",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetMatchedRoommateMatches(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *user.GetMatchedRoommateMatchesRequest, opts ...grpc.CallOption) (*user.GetMatchedRoommateMatchesResponse, error) {
						require.Equal(t, int64(10), req.UserId)
						return nil, status.Error(codes.Internal, "internal error")
					})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			ctx := context.Background()
			if test.userID != 0 {
				ctx = utils.SetUserID(ctx, test.userID)
			}

			req := httptest.NewRequest(http.MethodGet, "/user/me/roommate-matches/matched", nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			handler.GetMatchedRoommateMatches(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.expectedStatus == http.StatusOK {
				var resp user.GetMatchedRoommateMatchesResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, int64(1), resp.Len)
				require.Len(t, resp.Users, 1)
				require.Equal(t, int64(42), resp.Users[0].Id)
				require.Equal(t, "John", resp.Users[0].FirstName)
				require.Equal(t, "Doe", resp.Users[0].LastName)
				require.NotNil(t, resp.Users[0].AvatarUrl)
				require.Equal(t, "https://example.com/avatar.jpg", *resp.Users[0].AvatarUrl)
			}
		})
	}
}
