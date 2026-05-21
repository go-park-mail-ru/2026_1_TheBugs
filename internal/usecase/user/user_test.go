package user

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestUserUseCase_GetByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		userID    int
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		wantErr   error
	}{
		{
			name:   "OK",
			userID: 1,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				expected := &dto.UserDTO{ID: 1, FirstName: "John"}
				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)
				userRepoMock.EXPECT().GetByID(ctx, 1).Return(expected, nil)
			},
			wantErr: nil,
		},
		{
			name:   "NotFound",
			userID: 999,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)
				userRepoMock.EXPECT().GetByID(ctx, 999).Return(nil, entity.NotFoundError)
			},
			wantErr: entity.NotFoundError,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			_, err := uc.GetByID(ctx, test.userID)

			require.ErrorIs(t, err, test.wantErr)
		})
	}
}

func TestUserUseCase_UpdateProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		data      dto.UpdateProfileRequest
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo)
		wantErr   bool
	}{
		{
			name: "OK_NoAvatar",
			data: dto.UpdateProfileRequest{
				ID:        1,
				FirstName: strPtr("John"),
				LastName:  strPtr("Doe"),
				Phone:     strPtr("88001234567"),
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo) {
				expected := &dto.UserDTO{ID: 1, FirstName: "John"}
				userRepoMock.EXPECT().UpdateProfile(ctx, gomock.Any()).Return(expected, nil)
			},
			wantErr: false,
		},
		{
			name: "OK_WithValidAvatar",
			data: dto.UpdateProfileRequest{
				ID: 1,
				Avatar: &dto.FileInput{
					File:        validImageFile(),
					Size:        1024,
					ContentType: "image/jpeg",
					Filename:    "ss.jpeg",
				},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo) {
				expected := &dto.UserDTO{ID: 1}

				fileMock.EXPECT().Upload(ctx, gomock.Any(), gomock.Any(), int64(1024), "image/jpeg").Return(nil)
				userRepoMock.EXPECT().GetByID(ctx, gomock.Any()).Return(expected, nil)
				userRepoMock.EXPECT().UpdateProfile(ctx, gomock.Any()).Return(expected, nil)
			},
			wantErr: false,
		},
		{
			name: "InvalidPhone",
			data: dto.UpdateProfileRequest{
				ID:    1,
				Phone: strPtr("invalid_phone"),
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo) {
			},
			wantErr: true,
		},
		{
			name: "InvalidName",
			data: dto.UpdateProfileRequest{
				ID:        1,
				FirstName: strPtr(string(make([]byte, validator.MaxNameLenght+1))), // слишком короткое
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo) {
			},
			wantErr: true,
		},
		{
			name: "InvalidPhoto",
			data: dto.UpdateProfileRequest{
				ID: 1,
				Avatar: &dto.FileInput{
					File:        invalidImageFile(),
					Size:        1024,
					ContentType: "plain/text",
					Filename:    "ss.jpeg",
				},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo) {
				// file.Upload НЕ вызывается
			},
			wantErr: true,
		},
		{
			name: "FileUploadError",
			data: dto.UpdateProfileRequest{
				ID: 1,
				Avatar: &dto.FileInput{
					File:        validImageFile(),
					Size:        1024,
					ContentType: "image/jpeg",
					Filename:    "ss.jpeg",
				},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, userRepoMock *mocks.MockUserRepo) {
				expected := &dto.UserDTO{ID: 1}
				userRepoMock.EXPECT().GetByID(ctx, gomock.Any()).Return(expected, nil)
				fileMock.EXPECT().Upload(ctx, gomock.Any(), gomock.Any(), int64(1024), "image/jpeg").Return(errors.New("upload failed"))
			},
			wantErr: true,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)
			userRepoMock := mocks.NewMockUserRepo(ctrl)

			uowMock.EXPECT().Users().Return(userRepoMock).AnyTimes()

			test.setupMock(uowMock, fileMock, userRepoMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			_, err := uc.UpdateProfile(ctx, test.data)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func validImageFile() *mockFile {
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
	return &mockFile{reader: bytes.NewReader(jpegHeader)}
}

func invalidImageFile() *mockFile {
	return &mockFile{reader: bytes.NewReader([]byte("not an image"))}
}

func strPtr(s string) *string {
	return &s
}

type mockFile struct {
	reader io.Reader
}

func (m *mockFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *mockFile) Close() error {
	return nil
}

func TestUserUseCase_GetRoommateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	userID := 42

	roommateUser := &entity.RoommateUser{
		FirstName:   "John",
		LastName:    "Doe",
		AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: lo.ToPtr("good roommate"),
	}

	roommateTags := []entity.RoommateTag{
		{
			Name:  "Без животных",
			Alias: "no_pets",
		},
		{
			Name:  "Не курю",
			Alias: "no_smoking",
		},
	}

	expected := &dto.RoommateUserProfileDTO{
		FirstName:   "John",
		LastName:    "Doe",
		AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: lo.ToPtr("good roommate"),
		Tags: []dto.RoommateTagDTO{
			{
				Name:  "Без животных",
				Alias: "no_pets",
			},
			{
				Name:  "Не курю",
				Alias: "no_smoking",
			},
		},
	}

	cases := []struct {
		name      string
		userID    int
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		want      *dto.RoommateUserProfileDTO
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					GetRoommateUser(ctx, userID).
					Return(roommateUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateTags(ctx, userID).
					Return(roommateTags, nil).
					Times(1)
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name:   "GetRoommateUserNotFound",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetRoommateUser(ctx, userID).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "GetRoommateUserServiceError",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetRoommateUser(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "GetRoommateTagsError",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					GetRoommateUser(ctx, userID).
					Return(roommateUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateTags(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "EmptyTags",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					GetRoommateUser(ctx, userID).
					Return(roommateUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateTags(ctx, userID).
					Return([]entity.RoommateTag{}, nil).
					Times(1)
			},
			want: &dto.RoommateUserProfileDTO{
				FirstName:   "John",
				LastName:    "Doe",
				AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: lo.ToPtr("good roommate"),
				Tags:        []dto.RoommateTagDTO{},
			},
			wantErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			got, err := uc.GetRoommateUser(ctx, test.userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestUserUseCase_AddRoommateMatсh(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	fromUserID := 10
	toUserID := 42

	toContacts := &dto.RoommateContactsDTO{
		Email: "target@example.com",
		Phone: "+79991234567",
	}

	fromUser := &dto.UserDTO{
		ID:        fromUserID,
		FirstName: "John",
		LastName:  "Doe",
	}

	poster := &entity.PosterFlat{
		Address: "Moscow, Arbat 1",
	}

	cases := []struct {
		name       string
		fromUserID int
		toUserID   int
		setupMock  func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender)
		wantErr    error
	}{
		{
			name:       "OK",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				posterRepoMock := mocks.NewMockPosterRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(3)
				uowMock.EXPECT().Posters().Return(posterRepoMock).Times(1)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				posterRepoMock.EXPECT().
					GetRoommatePoster(ctx, fromUserID).
					Return(poster, nil).
					Times(1)

				senderMock.EXPECT().
					SendRoommateMatch(ctx, toContacts.Email, fromUser.FirstName, fromUser.LastName, poster.Address).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name:       "SelfMatesth",
			fromUserID: fromUserID,
			toUserID:   fromUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
			},
			wantErr: entity.InvalidInput,
		},
		{
			name:       "AddRoommateMatchNotFound",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(entity.NotFoundError).
					Times(1)
			},
			wantErr: entity.NotFoundError,
		},
		{
			name:       "AddRoommateMatchServiceError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetRoommateContactsError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetByIDError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(3)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetRoommatePosterError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				posterRepoMock := mocks.NewMockPosterRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(3)
				uowMock.EXPECT().Posters().Return(posterRepoMock).Times(1)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				posterRepoMock.EXPECT().
					GetRoommatePoster(ctx, fromUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:       "SenderErrorIgnored",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				posterRepoMock := mocks.NewMockPosterRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(3)
				uowMock.EXPECT().Posters().Return(posterRepoMock).Times(1)

				userRepoMock.EXPECT().
					AddRoommateMatch(ctx, fromUserID, toUserID).
					Return(nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				posterRepoMock.EXPECT().
					GetRoommatePoster(ctx, fromUserID).
					Return(poster, nil).
					Times(1)

				senderMock.EXPECT().
					SendRoommateMatch(ctx, toContacts.Email, fromUser.FirstName, fromUser.LastName, poster.Address).
					Return(errors.New("smtp error")).
					Times(1)
			},
			wantErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)
			senderMock := mocks.NewMockMailSender(ctrl)

			test.setupMock(uowMock, fileMock, senderMock)

			uc := NewUserUseCase(uowMock, fileMock, senderMock)
			err := uc.AddRoommateMatch(ctx, test.fromUserID, test.toUserID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUserUseCase_GetRoommateContacts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	fromUserID := 10
	toUserID := 42

	fromContacts := &dto.RoommateContactsDTO{
		Email: "from@example.com",
		Phone: "+79990000000",
	}

	toContacts := &dto.RoommateContactsDTO{
		Email: "target@example.com",
		Phone: "+79991234567",
	}

	fromUser := &dto.UserDTO{
		ID:        fromUserID,
		FirstName: "John",
		LastName:  "Doe",
	}

	toUser := &dto.UserDTO{
		ID:        toUserID,
		FirstName: "Jane",
		LastName:  "Smith",
	}

	poster := &entity.PosterFlat{
		Address: "Moscow, Arbat 1",
	}

	cases := []struct {
		name       string
		fromUserID int
		toUserID   int
		setupMock  func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender)
		want       *dto.RoommateContactsDTO
		wantErr    error
	}{
		{
			name:       "OK",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				posterRepoMock := mocks.NewMockPosterRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(5)
				uowMock.EXPECT().Posters().Return(posterRepoMock).Times(1)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(fromContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, toUserID).
					Return(toUser, nil).
					Times(1)

				posterRepoMock.EXPECT().
					GetRoommatePoster(ctx, fromUserID).
					Return(poster, nil).
					Times(1)

				senderMock.EXPECT().
					SendRoommateContactsForRequester(
						ctx,
						fromContacts.Email,
						toUser.FirstName,
						toUser.LastName,
						toContacts.Email,
						toContacts.Phone,
						poster.Address,
					).
					Return(nil).
					Times(1)

				senderMock.EXPECT().
					SendRoommateContactsForAccepted(
						ctx,
						toContacts.Email,
						fromUser.FirstName,
						fromUser.LastName,
						fromContacts.Email,
						fromContacts.Phone,
						poster.Address,
					).
					Return(nil).
					Times(1)
			},
			want:    toContacts,
			wantErr: nil,
		},
		{
			name:       "SelfContact",
			fromUserID: fromUserID,
			toUserID:   fromUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
			},
			want:    nil,
			wantErr: entity.InvalidInput,
		},
		{
			name:       "NotMatesthed",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(false, nil).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:       "IsRoommateMatesthError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(false, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetFromContactsError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:       "GetToContactsError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(3)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(fromContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetFromUserError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(4)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(fromContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetToUserError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(5)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(fromContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, toUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:       "GetRoommatePosterError",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				posterRepoMock := mocks.NewMockPosterRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(5)
				uowMock.EXPECT().Posters().Return(posterRepoMock).Times(1)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(fromContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, toUserID).
					Return(toUser, nil).
					Times(1)

				posterRepoMock.EXPECT().
					GetRoommatePoster(ctx, fromUserID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:       "SenderErrorIgnored",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo, senderMock *mocks.MockMailSender) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))
				posterRepoMock := mocks.NewMockPosterRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(5)
				uowMock.EXPECT().Posters().Return(posterRepoMock).Times(1)

				userRepoMock.EXPECT().
					IsRoommateMatch(ctx, fromUserID, toUserID).
					Return(true, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, fromUserID).
					Return(fromContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateContacts(ctx, toUserID).
					Return(toContacts, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, fromUserID).
					Return(fromUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetByID(ctx, toUserID).
					Return(toUser, nil).
					Times(1)

				posterRepoMock.EXPECT().
					GetRoommatePoster(ctx, fromUserID).
					Return(poster, nil).
					Times(1)

				senderMock.EXPECT().
					SendRoommateContactsForRequester(
						ctx,
						fromContacts.Email,
						toUser.FirstName,
						toUser.LastName,
						toContacts.Email,
						toContacts.Phone,
						poster.Address,
					).
					Return(errors.New("smtp error")).
					Times(1)

				senderMock.EXPECT().
					SendRoommateContactsForAccepted(
						ctx,
						toContacts.Email,
						fromUser.FirstName,
						fromUser.LastName,
						fromContacts.Email,
						fromContacts.Phone,
						poster.Address,
					).
					Return(errors.New("smtp error")).
					Times(1)
			},
			want:    toContacts,
			wantErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)
			senderMock := mocks.NewMockMailSender(ctrl)

			test.setupMock(uowMock, fileMock, senderMock)

			uc := NewUserUseCase(uowMock, fileMock, senderMock)
			got, err := uc.GetRoommateContacts(ctx, test.fromUserID, test.toUserID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestUserUseCase_CreateRoommateForm(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	validData := dto.CreateRoommateFormRequest{
		UserID:      10,
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	cases := []struct {
		name      string
		data      dto.CreateRoommateFormRequest
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		wantErr   error
	}{
		{
			name: "OK",
			data: validData,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					CreateRoommateForm(ctx, validData).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name: "InvalidUserID",
			data: dto.CreateRoommateFormRequest{
				UserID:      0,
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			wantErr: entity.InvalidInput,
		},
		{
			name: "InvalidGender",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "unknown",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			wantErr: entity.InvalidInput,
		},
		{
			name: "EmptyBirthday",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "male",
				Birthday:    "",
				Description: "good roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			wantErr: entity.InvalidInput,
		},
		{
			name: "EmptyDescription",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: "",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			wantErr: entity.InvalidInput,
		},
		{
			name: "RepoNotFound",
			data: validData,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					CreateRoommateForm(ctx, validData).
					Return(entity.NotFoundError).
					Times(1)
			},
			wantErr: entity.NotFoundError,
		},
		{
			name: "RepoServiceError",
			data: validData,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					CreateRoommateForm(ctx, validData).
					Return(entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			err := uc.CreateRoommateForm(ctx, test.data)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUserUseCase_GetRoommateForm(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	userID := 10

	form := &entity.RoommateForm{
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
	}

	tags := []string{"no_smoking", "no_pets"}

	expected := &dto.RoommateFormDTO{
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	cases := []struct {
		name      string
		userID    int
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		want      *dto.RoommateFormDTO
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					GetRoommateForm(ctx, userID).
					Return(form, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateFormTags(ctx, userID).
					Return(tags, nil).
					Times(1)
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name:   "InvalidUserID",
			userID: 0,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			want:    nil,
			wantErr: entity.InvalidInput,
		},
		{
			name:   "GetRoommateFormNotFound",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetRoommateForm(ctx, userID).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "GetRoommateFormServiceError",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetRoommateForm(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "GetRoommateFormTagsError",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					GetRoommateForm(ctx, userID).
					Return(form, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateFormTags(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "EmptyTags",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(2)

				userRepoMock.EXPECT().
					GetRoommateForm(ctx, userID).
					Return(form, nil).
					Times(1)

				userRepoMock.EXPECT().
					GetRoommateFormTags(ctx, userID).
					Return([]string{}, nil).
					Times(1)
			},
			want: &dto.RoommateFormDTO{
				Gender:      "male",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{},
			},
			wantErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			got, err := uc.GetRoommateForm(ctx, test.userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestUserUseCase_UpdateRoommateForm(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	validData := dto.CreateRoommateFormRequest{
		UserID:      10,
		Gender:      "female",
		Birthday:    "2001-02-03",
		Description: "updated roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	cases := []struct {
		name      string
		data      dto.CreateRoommateFormRequest
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		wantErr   error
	}{
		{
			name: "OK",
			data: validData,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					UpdateRoommateForm(ctx, validData).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name: "InvalidUserID",
			data: dto.CreateRoommateFormRequest{
				UserID:      0,
				Gender:      "female",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {},
			wantErr:   entity.InvalidInput,
		},
		{
			name: "InvalidGender",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "unknown",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {},
			wantErr:   entity.InvalidInput,
		},
		{
			name: "EmptyBirthday",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "female",
				Birthday:    "",
				Description: "updated roommate",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {},
			wantErr:   entity.InvalidInput,
		},
		{
			name: "EmptyDescription",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "female",
				Birthday:    "2001-02-03",
				Description: "",
				Tags:        []string{"no_smoking"},
			},
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {},
			wantErr:   entity.InvalidInput,
		},
		{
			name: "RepoNotFound",
			data: validData,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					UpdateRoommateForm(ctx, validData).
					Return(entity.NotFoundError).
					Times(1)
			},
			wantErr: entity.NotFoundError,
		},
		{
			name: "RepoServiceError",
			data: validData,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					UpdateRoommateForm(ctx, validData).
					Return(entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			err := uc.UpdateRoommateForm(ctx, test.data)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestUserUseCase_GetIncomingRoommateMatches(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	userID := 10

	users := []dto.RoommateUserDTO{
		{
			ID:        42,
			FirstName: "John",
			LastName:  "Doe",
			AvatarURL: lo.ToPtr("https://example.com/avatar.jpg"),
		},
		{
			ID:        43,
			FirstName: "Jane",
			LastName:  "Smith",
			AvatarURL: nil,
		},
	}

	expected := &dto.RoommateMatchesResponse{
		Users: users,
		Len:   len(users),
	}

	cases := []struct {
		name      string
		userID    int
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		want      *dto.RoommateMatchesResponse
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetIncomingRoommateMatches(ctx, userID).
					Return(users, nil).
					Times(1)
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name:   "InvalidUserID",
			userID: 0,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			want:    nil,
			wantErr: entity.InvalidInput,
		},
		{
			name:   "RepoNotFound",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetIncomingRoommateMatches(ctx, userID).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "RepoServiceError",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetIncomingRoommateMatches(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "EmptyList",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetIncomingRoommateMatches(ctx, userID).
					Return([]dto.RoommateUserDTO{}, nil).
					Times(1)
			},
			want: &dto.RoommateMatchesResponse{
				Users: []dto.RoommateUserDTO{},
				Len:   0,
			},
			wantErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			got, err := uc.GetIncomingRoommateMatches(ctx, test.userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestUserUseCase_GetMatchedRoommateMatches(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	userID := 10

	users := []dto.RoommateUserDTO{
		{
			ID:        42,
			FirstName: "John",
			LastName:  "Doe",
			AvatarURL: lo.ToPtr("https://example.com/avatar.jpg"),
		},
		{
			ID:        43,
			FirstName: "Jane",
			LastName:  "Smith",
			AvatarURL: nil,
		},
	}

	expected := &dto.RoommateMatchesResponse{
		Users: users,
		Len:   len(users),
	}

	cases := []struct {
		name      string
		userID    int
		setupMock func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo)
		want      *dto.RoommateMatchesResponse
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetMatchedRoommateMatches(ctx, userID).
					Return(users, nil).
					Times(1)
			},
			want:    expected,
			wantErr: nil,
		},
		{
			name:   "InvalidUserID",
			userID: 0,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
			},
			want:    nil,
			wantErr: entity.InvalidInput,
		},
		{
			name:   "RepoNotFound",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetMatchedRoommateMatches(ctx, userID).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "RepoServiceError",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetMatchedRoommateMatches(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "EmptyList",
			userID: userID,
			setupMock: func(uowMock *mocks.MockUnitOfWork, fileMock *mocks.MockFileRepo) {
				userRepoMock := mocks.NewMockUserRepo(gomock.NewController(t))

				uowMock.EXPECT().Users().Return(userRepoMock).Times(1)

				userRepoMock.EXPECT().
					GetMatchedRoommateMatches(ctx, userID).
					Return([]dto.RoommateUserDTO{}, nil).
					Times(1)
			},
			want: &dto.RoommateMatchesResponse{
				Users: []dto.RoommateUserDTO{},
				Len:   0,
			},
			wantErr: nil,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			test.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock, nil)
			got, err := uc.GetMatchedRoommateMatches(ctx, test.userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
