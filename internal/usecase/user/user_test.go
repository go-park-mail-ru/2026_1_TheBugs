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

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)

			tc.setupMock(uowMock, fileMock)

			uc := NewUserUseCase(uowMock, fileMock)
			_, err := uc.GetByID(ctx, tc.userID)

			require.ErrorIs(t, err, tc.wantErr)
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

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			fileMock := mocks.NewMockFileRepo(ctrl)
			userRepoMock := mocks.NewMockUserRepo(ctrl)

			uowMock.EXPECT().Users().Return(userRepoMock).AnyTimes()

			tc.setupMock(uowMock, fileMock, userRepoMock)

			uc := NewUserUseCase(uowMock, fileMock)
			_, err := uc.UpdateProfile(ctx, tc.data)

			if tc.wantErr {
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
