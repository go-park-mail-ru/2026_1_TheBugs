package auth

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/tokens"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUseCase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		data      dto.CreateUserDTO
		setupMock func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo)
		wantErr   error
	}{
		{
			name: "OK",
			data: dto.CreateUserDTO{
				Email:     "test@gmail.com",
				Password:  "dpofdOPOOo12",
				FirstName: "Mark",
				LastName:  "Mini",
				Phone:     "8 800 955 12 12",
			},
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
				gomock.InOrder(
					userMock.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, nil).Times(1),
					userMock.EXPECT().Create(ctx, gomock.Any()).Return(nil, nil).Times(1),
				)
			},
			wantErr: nil,
		},
		{
			name: "Conflict",
			data: dto.CreateUserDTO{
				Email:     "test@gmail.com",
				Password:  "dpofdOPOOo12",
				FirstName: "Mark",
				LastName:  "Mini",
				Phone:     "8 800 955 12 12",
			},
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
				existingUser := &entity.User{Email: "test"}
				userMock.EXPECT().GetByEmail(ctx, gomock.Any()).Return(existingUser, nil)
				userMock.EXPECT().Create(ctx, gomock.Any()).Times(0)
			},
			wantErr: entity.AlredyExitError,
		},
		{
			name: "InvalidInput",
			data: dto.CreateUserDTO{
				Email:     "wrong_email",
				Password:  "dpofdOPOOo121",
				FirstName: "Mark",
				LastName:  "Mini",
				Phone:     "8 800 955 12 12",
			},
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
			},
			wantErr: entity.InvalidInput,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			userMock := mocks.NewMockUserRepo(ctrl)
			authMock := mocks.NewMockAuthRepo(ctrl)

			uowMock.EXPECT().Users().Return(userMock).AnyTimes()
			uowMock.EXPECT().Autho().Return(authMock).AnyTimes()

			if tc.setupMock != nil {
				tc.setupMock(userMock, authMock)
			}

			uc := NewAuthUseCase(uowMock, nil)
			err := uc.RegisterUseCase(ctx, tc.data)

			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestLoginUseCase(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		email     string
		password  string
		setupMock func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo)
		wantErr   error
	}{
		{
			name:     "OK",
			email:    "test@gmail.com",
			password: "testgmailCom122",
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
				salt, _ := pwd.GenerateSalt()
				hashed := pwd.HashPassword("testgmailCom122", []byte(salt))
				existingUser := entity.User{
					ID:             0,
					Email:          "test@gmail.com",
					HashedPassword: &hashed,
					Salt:           &salt,
				}

				gomock.InOrder(
					userMock.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&existingUser, nil),
					authMock.EXPECT().CreateToken(ctx, gomock.Any()).Return(nil),
				)
			},
			wantErr: nil,
		},
		{
			name:     "Wrong Cred",
			email:    "test@gmail.com",
			password: "wrongPwd12323",
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
				salt, _ := pwd.GenerateSalt()
				hashed := pwd.HashPassword("testgmailCom122", []byte(salt))
				existingUser := entity.User{
					ID:             0,
					Email:          "test@gmail.com",
					HashedPassword: &hashed,
					Salt:           &salt,
				}
				userMock.EXPECT().GetByEmail(ctx, gomock.Any()).Return(&existingUser, nil)
			},
			wantErr: entity.BadCredentials,
		},
		{
			name:     "Invalid input pwd",
			email:    "test@gmail.com",
			password: "111111111",
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
			},
			wantErr: entity.InvalidInput,
		},
		{
			name:     "Not found",
			email:    "test@gmail.com",
			password: "wrongPwd12323",
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
				userMock.EXPECT().GetByEmail(ctx, gomock.Any()).Return(nil, entity.NotFoundError)
			},
			wantErr: entity.NotFoundError,
		},
		{
			name:     "Invalid input email",
			email:    "wrong_email",
			password: "dpofdOPOOo121",
			setupMock: func(userMock *mocks.MockUserRepo, authMock *mocks.MockAuthRepo) {
			},
			wantErr: entity.InvalidInput,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uowMock := mocks.NewMockUnitOfWork(ctrl)
			userMock := mocks.NewMockUserRepo(ctrl)
			authMock := mocks.NewMockAuthRepo(ctrl)

			uowMock.EXPECT().Users().Return(userMock).AnyTimes()
			uowMock.EXPECT().Autho().Return(authMock).AnyTimes()

			if tc.setupMock != nil {
				tc.setupMock(userMock, authMock)
			}
			testAccessToken := "dummy.access.value"
			testRefreshToken := "dummy.refresh.value"

			patchAccess := monkey.Patch(tokens.GenerateAccessToken, func(userID int, exp time.Duration) (string, error) {
				return testAccessToken, nil
			})
			patchRefresh := monkey.Patch(tokens.GenerateRefreshToken, func(tokenID string, userID int, exp time.Duration) (string, error) {
				return testRefreshToken, nil
			})
			defer func() {
				patchAccess.Unpatch()
				patchRefresh.Unpatch()
			}()

			uc := NewAuthUseCase(uowMock, nil)
			cred, err := uc.LoginUseCase(ctx, tc.email, tc.password)

			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				require.NotEmpty(t, cred)
				require.Equal(t, testAccessToken, cred.AccessToken)
				require.Equal(t, testRefreshToken, cred.RefreshToken)
			}
		})
	}
}

func TestRefreshUseCase(t *testing.T) {
	//t.Parallel()
	ctx := context.Background()

	userID := 1
	tokenID := "token.id"

	cases := []struct {
		name           string
		refreshToken   string
		setupMocks     func(uow *mocks.MockUnitOfWork, authMock *mocks.MockAuthRepo)
		patchParseFunc func() func()
		wantErr        error
		wantAccess     string
		wantRefresh    string
	}{
		{
			name:         "OK",
			refreshToken: "dummy.refresh.value",
			setupMocks: func(uow *mocks.MockUnitOfWork, authMock *mocks.MockAuthRepo) {
				storedToken := &entity.RefreshToken{
					ID:        0,
					TokenID:   tokenID,
					UserID:    userID,
					ExpiresAt: time.Now().Add(100 * time.Minute),
				}
				authMock.EXPECT().GetToken(ctx, gomock.Any(), gomock.Any()).Return(storedToken, nil)
				uow.EXPECT().
					Do(ctx, gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, fn func(usecase.UnitOfWork) error) error {
						return fn(uow)
					})
				gomock.InOrder(
					authMock.EXPECT().DeleteToken(ctx, gomock.Any(), gomock.Any()).Return(nil),
					authMock.EXPECT().CreateToken(ctx, gomock.Any()).Return(nil),
				)
			},
			patchParseFunc: func() func() {
				patch1 := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  strconv.Itoa(userID),
						Type: entity.RefreshTokenType,
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						},
					}, nil
				})
				patch2 := monkey.Patch(tokens.GenerateAccessToken, func(userID int, exp time.Duration) (string, error) {
					return "dummy.access.value", nil
				})
				patch3 := monkey.Patch(tokens.GenerateRefreshToken, func(tokenID string, userID int, exp time.Duration) (string, error) {
					return "dummy.refresh.value", nil
				})
				return func() {
					patch1.Unpatch()
					patch2.Unpatch()
					patch3.Unpatch()
				}
			},
			wantErr:     nil,
			wantAccess:  "dummy.access.value",
			wantRefresh: "dummy.refresh.value",
		},
		{
			name:         "BadParseToken",
			refreshToken: "dummy.refresh.value",
			setupMocks:   func(uow *mocks.MockUnitOfWork, authMock *mocks.MockAuthRepo) {},
			patchParseFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return nil, errors.New("bad token")
				})
				return func() { patch.Unpatch() }
			},
			wantErr: entity.JWTError,
		},
		{
			name:         "InvalidType",
			refreshToken: "dummy.refresh.value",
			setupMocks:   func(uow *mocks.MockUnitOfWork, authMock *mocks.MockAuthRepo) {},
			patchParseFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  strconv.Itoa(userID),
						Type: "invalid_type",
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						},
					}, nil
				})
				return func() { patch.Unpatch() }
			},
			wantErr: entity.JWTError,
		},
		{
			name:         "NotFoundToken",
			refreshToken: "dummy.refresh.value",
			setupMocks: func(uow *mocks.MockUnitOfWork, authMock *mocks.MockAuthRepo) {
				gomock.InOrder(
					authMock.EXPECT().GetToken(ctx, gomock.Any(), gomock.Any()).Return(nil, entity.NotFoundError),
				)
			},
			patchParseFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  strconv.Itoa(userID),
						Type: entity.RefreshTokenType,
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						},
					}, nil
				})
				return func() { patch.Unpatch() }
			},
			wantErr: entity.JWTError,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			//t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uowMock := mocks.NewMockUnitOfWork(ctrl)
			authMock := mocks.NewMockAuthRepo(ctrl)
			tokenMock := mocks.NewMockTokenRepo(ctrl)

			uowMock.EXPECT().Autho().Return(authMock).AnyTimes()

			if tc.setupMocks != nil {
				tc.setupMocks(uowMock, authMock)
			}

			unpatch := tc.patchParseFunc()
			defer unpatch()

			uc := NewAuthUseCase(uowMock, tokenMock)
			dto, err := uc.RefreshTokenUseCase(ctx, tc.refreshToken)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantAccess, dto.AccessToken)
				require.Equal(t, tc.wantRefresh, dto.RefreshToken)
			}
		})
	}
}
