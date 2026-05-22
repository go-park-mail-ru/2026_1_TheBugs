package user

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userpb "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/user"
)

func TestUserServiceServer_GetMe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int32(42)
	expectedUser := &dto.UserDTO{
		ID:        42,
		Email:     "user@example.com",
		Phone:     "+1234567890",
		AvatarURL: lo.ToPtr("https://example.com/avatar.png"),
		FirstName: "John",
		LastName:  "Doe",
	}

	tests := []struct {
		name      string
		req       *userpb.GetMeRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  &userpb.GetMeRequest{UserId: userID},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetByID(ctx, int(userID)).
					Return(expectedUser, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req:  &userpb.GetMeRequest{UserId: userID},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetByID(ctx, int(userID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "internal error",
			req:  &userpb.GetMeRequest{UserId: userID},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetByID(ctx, int(userID)).
					Return(nil, errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.GetMe(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedUser.ID, int(resp.Id))
				require.Equal(t, expectedUser.Email, resp.Email)
				require.Equal(t, expectedUser.Phone, resp.Phone)
				require.Equal(t, expectedUser.AvatarURL, resp.AvatarUrl)
				require.Equal(t, expectedUser.FirstName, resp.Firstname)
				require.Equal(t, expectedUser.LastName, resp.Lastname)
			}
		})
	}
}

func TestUserServiceServer_UpdateProfile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int32(42)
	avatarData := []byte("fake image data")

	tests := []struct {
		name      string
		req       *userpb.UpdateProfileRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success without avatar",
			req: &userpb.UpdateProfileRequest{
				Id:        userID,
				Firstname: lo.ToPtr("Jane"),
				Lastname:  lo.ToPtr("Smith"),
				Phone:     lo.ToPtr("+987654321"),
				File:      nil,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateProfile(ctx, dto.UpdateProfileRequest{
						ID:        int(userID),
						FirstName: lo.ToPtr("Jane"),
						LastName:  lo.ToPtr("Smith"),
						Phone:     lo.ToPtr("+987654321"),
						Avatar:    nil,
					}).
					Return(&dto.UserDTO{
						ID:        int(userID),
						Email:     "user@example.com",
						Phone:     "+987654321",
						AvatarURL: nil,
						FirstName: "Jane",
						LastName:  "Smith",
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "success with avatar",
			req: &userpb.UpdateProfileRequest{
				Id:        userID,
				Firstname: lo.ToPtr("Jane"),
				Lastname:  lo.ToPtr("Smith"),
				Phone:     lo.ToPtr("+987654321"),
				File: &userpb.UploadFile{
					Filename:    "avatar.jpg",
					Avatar:      avatarData,
					Size:        int64(len(avatarData)),
					ContentType: "image/jpeg",
				},
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateProfile(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, req dto.UpdateProfileRequest) (*dto.UserDTO, error) {
						// verify FileInput fields
						require.NotNil(t, req.Avatar)
						require.Equal(t, "avatar.jpg", req.Avatar.Filename)
						require.Equal(t, int64(len(avatarData)), req.Avatar.Size)
						require.Equal(t, "image/jpeg", req.Avatar.ContentType)

						// read and verify content
						content, err := io.ReadAll(req.Avatar.File)
						require.NoError(t, err)
						require.Equal(t, avatarData, content)

						return &dto.UserDTO{
							ID:        int(userID),
							Email:     "user@example.com",
							Phone:     "+987654321",
							AvatarURL: lo.ToPtr("https://example.com/new-avatar.jpg"),
							FirstName: "Jane",
							LastName:  "Smith",
						}, nil
					})
			},
			wantErr: false,
		},
		{
			name: "use case error - not found",
			req: &userpb.UpdateProfileRequest{
				Id:        userID,
				Firstname: lo.ToPtr("Jane"),
				Lastname:  lo.ToPtr("Smith"),
				Phone:     lo.ToPtr("+987654321"),
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateProfile(ctx, gomock.Any()).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - invalid input",
			req: &userpb.UpdateProfileRequest{
				Id:        userID,
				Firstname: nil, // invalid
				Lastname:  nil,
				Phone:     nil,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateProfile(ctx, gomock.Any()).
					Return(nil, entity.InvalidInput)
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error - internal",
			req: &userpb.UpdateProfileRequest{
				Id:        userID,
				Firstname: lo.ToPtr("Jane"),
				Lastname:  lo.ToPtr("Smith"),
				Phone:     lo.ToPtr("+987654321"),
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateProfile(ctx, gomock.Any()).
					Return(nil, errors.New("some internal error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.UpdateProfile(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				if tt.req.Firstname != nil {
					require.Equal(t, *tt.req.Firstname, resp.Firstname)
				}
				if tt.req.Lastname != nil {
					require.Equal(t, *tt.req.Lastname, resp.Lastname)
				}
				if tt.req.Phone != nil {
					require.Equal(t, *tt.req.Phone, resp.Phone)
				}
				require.Equal(t, int(tt.req.Id), int(resp.Id))
			}
		})
	}
}

func TestUserServiceServer_GetRoommateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(42)
	expectedUser := &dto.RoommateUserProfileDTO{
		FirstName:   "John",
		LastName:    "Doe",
		AvatarURL:   lo.ToPtr("https://example.com/avatar.png"),
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: lo.ToPtr("good roommate"),
		Tags: []dto.RoommateTagDTO{
			{
				Name:  "Не курю",
				Alias: "no_smoking",
			},
		},
	}

	tests := []struct {
		name      string
		req       *userpb.GetRoommateUserRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  &userpb.GetRoommateUserRequest{UserId: userID},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateUser(ctx, int(userID)).
					Return(expectedUser, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req:  &userpb.GetRoommateUserRequest{UserId: userID},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateUser(ctx, int(userID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "internal error",
			req:  &userpb.GetRoommateUserRequest{UserId: userID},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateUser(ctx, int(userID)).
					Return(nil, errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.GetRoommateUser(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedUser.FirstName, resp.FirstName)
				require.Equal(t, expectedUser.LastName, resp.LastName)
				require.Equal(t, expectedUser.AvatarURL, resp.AvatarUrl)
				require.Equal(t, expectedUser.Gender, resp.Gender)
				require.Equal(t, expectedUser.Birthday, resp.Birthday)
				require.Equal(t, expectedUser.Description, resp.Description)
				require.Len(t, resp.Tags, 1)
				require.Equal(t, expectedUser.Tags[0].Name, resp.Tags[0].Name)
				require.Equal(t, expectedUser.Tags[0].Alias, resp.Tags[0].Alias)
			}
		})
	}
}

func TestUserServiceServer_AddRoommateMatch(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fromUserID := int64(10)
	toUserID := int64(42)
	posterAlias := "poster-alias"

	tests := []struct {
		name      string
		req       *userpb.AddRoommateMatchRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success with poster alias",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId:  fromUserID,
				ToUserId:    toUserID,
				PosterAlias: &posterAlias,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					AddRoommateMatch(ctx, int(fromUserID), int(toUserID), &posterAlias).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success without poster alias",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId: fromUserID,
				ToUserId:   toUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					AddRoommateMatch(ctx, int(fromUserID), int(toUserID), nil).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid from user id",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId: 0,
				ToUserId:   toUserID,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "invalid to user id",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId: fromUserID,
				ToUserId:   0,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - not found",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId: fromUserID,
				ToUserId:   toUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					AddRoommateMatch(ctx, int(fromUserID), int(toUserID), nil).
					Return(entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - invalid input",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId: fromUserID,
				ToUserId:   fromUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					AddRoommateMatch(ctx, int(fromUserID), int(fromUserID), nil).
					Return(entity.InvalidInput)
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error - internal",
			req: &userpb.AddRoommateMatchRequest{
				FromUserId: fromUserID,
				ToUserId:   toUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					AddRoommateMatch(ctx, int(fromUserID), int(toUserID), nil).
					Return(errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.AddRoommateMatch(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestUserServiceServer_GetRoommateContacts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fromUserID := int64(10)
	toUserID := int64(42)

	expectedContacts := &dto.RoommateContactsDTO{
		Email: "target@example.com",
		Phone: "+79991234567",
	}

	tests := []struct {
		name      string
		req       *userpb.GetRoommateContactsRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req: &userpb.GetRoommateContactsRequest{
				FromUserId: fromUserID,
				ToUserId:   toUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateContacts(ctx, int(fromUserID), int(toUserID)).
					Return(expectedContacts, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid from user id",
			req: &userpb.GetRoommateContactsRequest{
				FromUserId: 0,
				ToUserId:   toUserID,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "invalid to user id",
			req: &userpb.GetRoommateContactsRequest{
				FromUserId: fromUserID,
				ToUserId:   0,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - not found",
			req: &userpb.GetRoommateContactsRequest{
				FromUserId: fromUserID,
				ToUserId:   toUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateContacts(ctx, int(fromUserID), int(toUserID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - invalid input",
			req: &userpb.GetRoommateContactsRequest{
				FromUserId: fromUserID,
				ToUserId:   fromUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateContacts(ctx, int(fromUserID), int(fromUserID)).
					Return(nil, entity.InvalidInput)
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error - internal",
			req: &userpb.GetRoommateContactsRequest{
				FromUserId: fromUserID,
				ToUserId:   toUserID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateContacts(ctx, int(fromUserID), int(toUserID)).
					Return(nil, errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.GetRoommateContacts(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedContacts.Email, resp.Email)
				require.Equal(t, expectedContacts.Phone, resp.Phone)
			}
		})
	}
}

func TestUserServiceServer_CreateRoommateForm(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(10)

	validReq := &userpb.CreateRoommateFormRequest{
		UserId:      userID,
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	tests := []struct {
		name      string
		req       *userpb.CreateRoommateFormRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					CreateRoommateForm(ctx, dto.CreateRoommateFormRequest{
						UserID:      int(userID),
						Gender:      "male",
						Birthday:    "2000-01-01",
						Description: "good roommate",
						Tags:        []string{"no_smoking", "no_pets"},
					}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &userpb.CreateRoommateFormRequest{
				UserId:      0,
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing gender",
			req: &userpb.CreateRoommateFormRequest{
				UserId:      userID,
				Gender:      "",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing birthday",
			req: &userpb.CreateRoommateFormRequest{
				UserId:      userID,
				Gender:      "male",
				Birthday:    "",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing description",
			req: &userpb.CreateRoommateFormRequest{
				UserId:      userID,
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: "",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - invalid input",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					CreateRoommateForm(ctx, gomock.Any()).
					Return(entity.InvalidInput)
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error - internal",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					CreateRoommateForm(ctx, gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.CreateRoommateForm(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestUserServiceServer_GetRoommateForm(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(10)

	expectedForm := &dto.RoommateFormDTO{
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	tests := []struct {
		name      string
		req       *userpb.GetRoommateFormRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req: &userpb.GetRoommateFormRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateForm(ctx, int(userID)).
					Return(expectedForm, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &userpb.GetRoommateFormRequest{
				UserId: 0,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - not found",
			req: &userpb.GetRoommateFormRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateForm(ctx, int(userID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - internal",
			req: &userpb.GetRoommateFormRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetRoommateForm(ctx, int(userID)).
					Return(nil, errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.GetRoommateForm(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, expectedForm.Gender, resp.Gender)
				require.Equal(t, expectedForm.Birthday, resp.Birthday)
				require.Equal(t, expectedForm.Description, resp.Description)
				require.Equal(t, expectedForm.Tags, resp.Tags)
			}
		})
	}
}

func TestUserServiceServer_UpdateRoommateForm(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(10)

	validReq := &userpb.UpdateRoommateFormRequest{
		UserId:      userID,
		Gender:      "female",
		Birthday:    "2001-02-03",
		Description: "updated roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	tests := []struct {
		name      string
		req       *userpb.UpdateRoommateFormRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateRoommateForm(ctx, dto.CreateRoommateFormRequest{
						UserID:      int(userID),
						Gender:      "female",
						Birthday:    "2001-02-03",
						Description: "updated roommate",
						Tags:        []string{"no_smoking", "no_pets"},
					}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &userpb.UpdateRoommateFormRequest{
				UserId:      0,
				Gender:      "female",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing gender",
			req: &userpb.UpdateRoommateFormRequest{
				UserId:      userID,
				Gender:      "",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing birthday",
			req: &userpb.UpdateRoommateFormRequest{
				UserId:      userID,
				Gender:      "female",
				Birthday:    "",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "missing description",
			req: &userpb.UpdateRoommateFormRequest{
				UserId:      userID,
				Gender:      "female",
				Birthday:    "2001-02-03",
				Description: "",
				Tags:        []string{"no_smoking"},
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - invalid input",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateRoommateForm(ctx, gomock.Any()).
					Return(entity.InvalidInput)
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error - not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateRoommateForm(ctx, gomock.Any()).
					Return(entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "use case error - internal",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					UpdateRoommateForm(ctx, gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.UpdateRoommateForm(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}

func TestUserServiceServer_GetIncomingRoommateMatches(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(10)

	expectedResp := &dto.RoommateMatchesResponse{
		Users: []dto.RoommateUserDTO{
			{
				ID:          42,
				FirstName:   "John",
				LastName:    "Doe",
				AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
				PosterAlias: lo.ToPtr("flat-1"),
			},
		},
		Len: 1,
	}

	tests := []struct {
		name      string
		req       *userpb.GetIncomingRoommateMatchesRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req: &userpb.GetIncomingRoommateMatchesRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetIncomingRoommateMatches(ctx, int(userID)).
					Return(expectedResp, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &userpb.GetIncomingRoommateMatchesRequest{
				UserId: 0,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - internal",
			req: &userpb.GetIncomingRoommateMatchesRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetIncomingRoommateMatches(ctx, int(userID)).
					Return(nil, errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.GetIncomingRoommateMatches(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				require.Equal(t, int64(expectedResp.Len), resp.Len)
				require.Len(t, resp.Users, 1)
				require.Equal(t, int64(expectedResp.Users[0].ID), resp.Users[0].Id)
				require.Equal(t, expectedResp.Users[0].FirstName, resp.Users[0].FirstName)
				require.Equal(t, expectedResp.Users[0].LastName, resp.Users[0].LastName)
				require.Equal(t, expectedResp.Users[0].AvatarURL, resp.Users[0].AvatarUrl)
				require.Equal(t, expectedResp.Users[0].PosterAlias, resp.Users[0].PosterAlias)
			}
		})
	}
}

func TestUserServiceServer_GetMatchedRoommateMatches(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := int64(10)

	expectedResp := &dto.RoommateMatchesResponse{
		Users: []dto.RoommateUserDTO{
			{
				ID:        42,
				FirstName: "John",
				LastName:  "Doe",
				AvatarURL: lo.ToPtr("https://example.com/avatar.jpg"),
			},
		},
		Len: 1,
	}

	tests := []struct {
		name      string
		req       *userpb.GetMatchedRoommateMatchesRequest
		setupMock func(mockUC *mocks.MockUserUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req: &userpb.GetMatchedRoommateMatchesRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetMatchedRoommateMatches(ctx, int(userID)).
					Return(expectedResp, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &userpb.GetMatchedRoommateMatchesRequest{
				UserId: 0,
			},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "use case error - internal",
			req: &userpb.GetMatchedRoommateMatchesRequest{
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockUserUseCase) {
				mockUC.EXPECT().
					GetMatchedRoommateMatches(ctx, int(userID)).
					Return(nil, errors.New("database error"))
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockUserUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewUserServiceServer(mockUC)
			resp, err := server.GetMatchedRoommateMatches(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				require.Equal(t, int64(expectedResp.Len), resp.Len)
				require.Len(t, resp.Users, 1)
				require.Equal(t, int64(expectedResp.Users[0].ID), resp.Users[0].Id)
				require.Equal(t, expectedResp.Users[0].FirstName, resp.Users[0].FirstName)
				require.Equal(t, expectedResp.Users[0].LastName, resp.Users[0].LastName)
				require.Equal(t, expectedResp.Users[0].AvatarURL, resp.Users[0].AvatarUrl)
			}
		})
	}
}
