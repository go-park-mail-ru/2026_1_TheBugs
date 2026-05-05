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
