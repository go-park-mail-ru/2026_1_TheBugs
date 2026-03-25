package auth

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/domains"
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

			uc := NewAuthUseCase(uowMock, nil, nil)
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

			uc := NewAuthUseCase(uowMock, nil, nil)
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
			tokenMock := mocks.NewMockСache(ctrl)

			uowMock.EXPECT().Autho().Return(authMock).AnyTimes()

			if tc.setupMocks != nil {
				tc.setupMocks(uowMock, authMock)
			}

			unpatch := tc.patchParseFunc()
			defer unpatch()

			uc := NewAuthUseCase(uowMock, tokenMock, nil)
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
func TestValidateAccessToken(t *testing.T) {
	tokenID := "token-id"

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name      string
		token     string
		setupMock func(c *mocks.MockСache)
		patchFunc func() func()
		wantErr   bool
	}

	tests := []testCase{
		{
			name:  "success",
			token: "access.token",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  "sub",
						Type: entity.AccessTokenType,
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						}}, nil
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache) {
				c.EXPECT().
					IsBlacklisted(ctx, tokenID).
					Return(false, nil)
			},
			wantErr: false,
		},
		{
			name:  "blacklisted",
			token: "access.token",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  "sub",
						Type: entity.AccessTokenType,
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						}}, nil
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache) {
				c.EXPECT().
					IsBlacklisted(ctx, tokenID).
					Return(true, nil)
			},
			wantErr: true,
		},
		{
			name:  "invalid token",
			token: "bad-token",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return nil, errors.New("bad token")
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache) {},
			wantErr:   true,
		},
		{
			name: "cache error",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  "sub",
						Type: entity.AccessTokenType,
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						}}, nil
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache) {
				c.EXPECT().
					IsBlacklisted(ctx, tokenID).
					Return(false, errors.New("redis down"))
			},
			wantErr: true,
		},
		{
			name: "invalid token type",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{Sub: "sub", Type: "bad type"}, nil
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := mocks.NewMockСache(ctrl)
			tt.setupMock(cache)
			unpatch := tt.patchFunc()
			defer unpatch()

			uc := AuthUseCase{
				cache: cache,
			}

			_, err := uc.ValidateAccessToken(ctx, tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {
	tokenID := "refresh-id"
	userID := 1

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type testCase struct {
		name      string
		token     string
		setupMock func(r *mocks.MockAuthRepo)
		patchFunc func() func()
		wantErr   bool
	}

	tests := []testCase{
		{
			name:  "success",
			token: "refresh.token",
			patchFunc: func() func() {
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
			setupMock: func(r *mocks.MockAuthRepo) {
				r.EXPECT().
					GetToken(ctx, tokenID, userID).
					Return(&entity.RefreshToken{
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)
			},
			wantErr: false,
		},
		{
			name:  "expired",
			token: "refresh.token",
			patchFunc: func() func() {
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
			setupMock: func(r *mocks.MockAuthRepo) {
				r.EXPECT().
					GetToken(ctx, tokenID, userID).
					Return(&entity.RefreshToken{
						ExpiresAt: time.Now().Add(-time.Hour),
					}, nil)
			},
			wantErr: true,
		},
		{
			name:  "not found",
			token: "refresh.token",
			patchFunc: func() func() {
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
			setupMock: func(r *mocks.MockAuthRepo) {
				r.EXPECT().
					GetToken(ctx, tokenID, userID).
					Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:  "repo error",
			token: "refresh.token",
			patchFunc: func() func() {
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
			setupMock: func(r *mocks.MockAuthRepo) {
				r.EXPECT().
					GetToken(ctx, tokenID, userID).
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:  "invalid token",
			token: "bad-token",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return nil, errors.New("bad token")
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(r *mocks.MockAuthRepo) {},
			wantErr:   true,
		},
		{
			name:  "invalid type",
			token: "refresh.token",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  strconv.Itoa(userID),
						Type: "wrong",
					}, nil
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(r *mocks.MockAuthRepo) {},
			wantErr:   true,
		},
		{
			name:  "invalid userID",
			token: "refresh.token",
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return &tokens.Claims{
						Sub:  "not-int",
						Type: entity.RefreshTokenType,
						RegisteredClaims: jwt.RegisteredClaims{
							ID: tokenID,
						},
					}, nil
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(r *mocks.MockAuthRepo) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockUOW := mocks.NewMockUnitOfWork(ctl)
			mockAuth := mocks.NewMockAuthRepo(ctl)
			mockCache := mocks.NewMockСache(ctl)
			mockSender := mocks.NewMockMailSender(ctl)

			mockUOW.EXPECT().Autho().Return(mockAuth).AnyTimes()

			tt.setupMock(mockAuth)
			unpatch := tt.patchFunc()
			defer unpatch()

			uc := NewAuthUseCase(mockUOW, mockCache, mockSender)

			_, _, err := uc.ValidateRefreshToken(ctx, tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogoutUseCase(t *testing.T) {
	accessID := "access-id"
	refreshID := "refresh-id"
	userID := 1

	ctx := context.Background()

	type testCase struct {
		name      string
		dto       dto.LogoutDTO
		setupMock func(c *mocks.MockСache, r *mocks.MockAuthRepo)
		patchFunc func() func()
		wantErr   bool
	}

	tests := []testCase{
		{
			name: "success",
			dto: dto.LogoutDTO{
				AccessToken:  "access.token",
				RefreshToken: "refresh.token",
			},
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					switch token {
					case "access.token":
						return &tokens.Claims{
							Type: entity.AccessTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: accessID,
							},
						}, nil
					case "refresh.token":
						return &tokens.Claims{
							Sub:  strconv.Itoa(userID),
							Type: entity.RefreshTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: refreshID,
							},
						}, nil
					default:
						return nil, errors.New("bad token")
					}
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache, r *mocks.MockAuthRepo) {
				// access validation
				c.EXPECT().
					IsBlacklisted(ctx, accessID).
					Return(false, nil)

				// refresh validation
				r.EXPECT().
					GetToken(ctx, refreshID, userID).
					Return(&entity.RefreshToken{
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)

				// delete refresh
				r.EXPECT().
					DeleteToken(ctx, refreshID, userID).
					Return(nil)

				// blacklist access
				c.EXPECT().
					SetBlacklist(ctx, accessID, gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "access invalid",
			dto: dto.LogoutDTO{
				AccessToken:  "bad",
				RefreshToken: "refresh.token",
			},
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					return nil, errors.New("bad token")
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache, r *mocks.MockAuthRepo) {},
			wantErr:   true,
		},
		{
			name: "refresh invalid",
			dto: dto.LogoutDTO{
				AccessToken:  "access.token",
				RefreshToken: "bad",
			},
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					if token == "access.token" {
						return &tokens.Claims{
							Type: entity.AccessTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: accessID,
							},
						}, nil
					}
					return nil, errors.New("bad token")
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache, r *mocks.MockAuthRepo) {
				c.EXPECT().
					IsBlacklisted(ctx, accessID).
					Return(false, nil)
			},
			wantErr: true,
		},
		{
			name: "delete error",
			dto: dto.LogoutDTO{
				AccessToken:  "access.token",
				RefreshToken: "refresh.token",
			},
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					switch token {
					case "access.token":
						return &tokens.Claims{
							Type: entity.AccessTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: accessID,
							},
						}, nil
					case "refresh.token":
						return &tokens.Claims{
							Sub:  strconv.Itoa(userID),
							Type: entity.RefreshTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: refreshID,
							},
						}, nil
					default:
						return nil, errors.New("bad token")
					}
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache, r *mocks.MockAuthRepo) {
				c.EXPECT().
					IsBlacklisted(ctx, accessID).
					Return(false, nil)

				r.EXPECT().
					GetToken(ctx, refreshID, userID).
					Return(&entity.RefreshToken{
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)

				r.EXPECT().
					DeleteToken(ctx, refreshID, userID).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "blacklist error",
			dto: dto.LogoutDTO{
				AccessToken:  "access.token",
				RefreshToken: "refresh.token",
			},
			patchFunc: func() func() {
				patch := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
					switch token {
					case "access.token":
						return &tokens.Claims{
							Type: entity.AccessTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: accessID,
							},
						}, nil
					case "refresh.token":
						return &tokens.Claims{
							Sub:  strconv.Itoa(userID),
							Type: entity.RefreshTokenType,
							RegisteredClaims: jwt.RegisteredClaims{
								ID: refreshID,
							},
						}, nil
					default:
						return nil, errors.New("bad token")
					}
				})
				return func() { patch.Unpatch() }
			},
			setupMock: func(c *mocks.MockСache, r *mocks.MockAuthRepo) {
				c.EXPECT().
					IsBlacklisted(ctx, accessID).
					Return(false, nil)

				r.EXPECT().
					GetToken(ctx, refreshID, userID).
					Return(&entity.RefreshToken{
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)

				r.EXPECT().
					DeleteToken(ctx, refreshID, userID).
					Return(nil)

				c.EXPECT().
					SetBlacklist(ctx, accessID, gomock.Any()).
					Return(errors.New("redis down"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUOW := mocks.NewMockUnitOfWork(ctrl)
			mockAuth := mocks.NewMockAuthRepo(ctrl)
			mockCache := mocks.NewMockСache(ctrl)
			mockSender := mocks.NewMockMailSender(ctrl)

			mockUOW.EXPECT().
				Autho().
				Return(mockAuth).
				AnyTimes()

			tt.setupMock(mockCache, mockAuth)

			unpatch := tt.patchFunc()
			defer unpatch()

			uc := NewAuthUseCase(mockUOW, mockCache, mockSender)

			err := uc.LogoutUseCase(ctx, tt.dto)

			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendVerificationCode(t *testing.T) {
	ctx := context.Background()
	email := "test@email.com"

	tests := []struct {
		name  string
		setup func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache, mockSender *mocks.MockMailSender)
		err   error
	}{
		{
			name: "OK",
			setup: func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache, mockSender *mocks.MockMailSender) {
				gomock.InOrder(
					mockUser.EXPECT().GetByEmail(ctx, email).Return(nil, nil),
					mockCache.EXPECT().IsBlacklisted(ctx, email).Return(false, nil),
					mockCache.EXPECT().SetBlacklist(ctx, email, gomock.Any()).Return(nil),
					mockCache.EXPECT().CreateRecoverSession(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
				)
				mockSender.EXPECT().SendCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "Email not found",
			setup: func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache, mockSender *mocks.MockMailSender) {
				mockUser.EXPECT().GetByEmail(ctx, email).Return(nil, entity.NotFoundError)
			},
			err: entity.NotFoundError,
		},
		{
			name: "Email in blacklis",
			setup: func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache, mockSender *mocks.MockMailSender) {
				mockUser.EXPECT().GetByEmail(ctx, email).Return(nil, nil)
				mockCache.EXPECT().IsBlacklisted(ctx, email).Return(true, nil)
			},
			err: entity.ToManyRequest,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockUOW := mocks.NewMockUnitOfWork(ctl)
			mockUser := mocks.NewMockUserRepo(ctl)
			mockCache := mocks.NewMockСache(ctl)
			mockSender := mocks.NewMockMailSender(ctl)

			mockUOW.EXPECT().Users().Return(mockUser).AnyTimes()

			if tc.setup != nil {
				tc.setup(mockUOW, mockUser, mockCache, mockSender)
			}

			authUC := NewAuthUseCase(mockUOW, mockCache, mockSender)

			_, err := authUC.SendVerificationCode(ctx, email)
			require.ErrorIs(t, err, tc.err)
		})
	}

}

func TestCheckRecoveryCode(t *testing.T) {
	ctx := context.Background()
	sessionID := "test-session-id"

	tests := []struct {
		name  string
		code  string
		setup func(mockCache *mocks.MockСache)
		err   error
	}{
		{
			name: "OK",
			code: "code",
			setup: func(mockCache *mocks.MockСache) {
				session := &domains.RecoverSession{
					Email:    "test@email.com",
					Code:     "code",
					Attempts: 0,
					Verified: false,
				}

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
					mockCache.EXPECT().SetRecoverVerified(ctx, sessionID, gomock.Any()).Return(nil))
			},
			err: nil,
		},
		{
			name: "Session is empty",
			code: "code",
			setup: func(mockCache *mocks.MockСache) {
				var session *domains.RecoverSession

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
				)
			},
			err: entity.BadCredentials,
		},
		{
			name: "Invalid code",
			code: "wrong_code",
			setup: func(mockCache *mocks.MockСache) {
				session := &domains.RecoverSession{
					Email:    "test@email.com",
					Code:     "code",
					Attempts: 0,
					Verified: false,
				}

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
					mockCache.EXPECT().IncrementRecoverAttempts(ctx, sessionID).Return(int64(1), nil),
				)
			},
			err: entity.BadCredentials,
		},
		{
			name: "Invalid code and max limit",
			code: "wrong_code",
			setup: func(mockCache *mocks.MockСache) {
				session := &domains.RecoverSession{
					Email:    "test@email.com",
					Code:     "code",
					Attempts: MaxAttemptsRecovery,
					Verified: false,
				}

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
					mockCache.EXPECT().IncrementRecoverAttempts(ctx, sessionID).Return(int64(MaxAttemptsRecovery+1), nil),
					mockCache.EXPECT().DeleteRecoverSession(ctx, sessionID).Return(nil),
				)
			},
			err: entity.ToManyRequest,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockUOW := mocks.NewMockUnitOfWork(ctl)
			mockUser := mocks.NewMockUserRepo(ctl)
			mockCache := mocks.NewMockСache(ctl)
			mockSender := mocks.NewMockMailSender(ctl)

			mockUOW.EXPECT().Users().Return(mockUser).AnyTimes()

			if tc.setup != nil {
				tc.setup(mockCache)
			}

			authUC := NewAuthUseCase(mockUOW, mockCache, mockSender)

			err := authUC.CheckRecoveryCode(ctx, sessionID, tc.code)
			require.ErrorIs(t, err, tc.err)
		})
	}

}

func TestUpdateUserPassword(t *testing.T) {
	ctx := context.Background()
	sessionID := "test-session-id"

	tests := []struct {
		name  string
		pwd   string
		setup func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache)
		err   error
	}{
		{
			name: "OK",
			pwd:  "new_pwd",
			setup: func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache) {
				email := "test@email.com"
				session := &domains.RecoverSession{
					Email:    email,
					Verified: true,
				}

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
					mockUser.EXPECT().UpdatePwd(ctx, email, gomock.Any(), gomock.Any()).Return(nil),
					mockCache.EXPECT().DeleteRecoverSession(ctx, sessionID).Return(nil),
				)
			},
			err: nil,
		},
		{
			name: "Empty session",
			pwd:  "new_pwd",
			setup: func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache) {
				var session *domains.RecoverSession

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
				)
			},
			err: entity.BadCredentials,
		},

		{
			name: "Session is not verified",
			pwd:  "new_pwd",
			setup: func(uow *mocks.MockUnitOfWork, mockUser *mocks.MockUserRepo, mockCache *mocks.MockСache) {
				email := "test@email.com"
				session := &domains.RecoverSession{
					Email:    email,
					Verified: false,
				}

				gomock.InOrder(
					mockCache.EXPECT().GetRecoverSession(ctx, sessionID).Return(session, nil),
				)
			},
			err: entity.BadCredentials,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			mockUOW := mocks.NewMockUnitOfWork(ctl)
			mockUser := mocks.NewMockUserRepo(ctl)
			mockCache := mocks.NewMockСache(ctl)
			mockSender := mocks.NewMockMailSender(ctl)

			mockUOW.EXPECT().Users().Return(mockUser).AnyTimes()

			if tc.setup != nil {
				tc.setup(mockUOW, mockUser, mockCache)
			}

			authUC := NewAuthUseCase(mockUOW, mockCache, mockSender)

			err := authUC.UpdateUserPassword(ctx, sessionID, tc.pwd)
			require.ErrorIs(t, err, tc.err)
		})
	}

}
