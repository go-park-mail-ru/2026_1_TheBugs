package auth

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/tokens"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUseCase_OK(t *testing.T) {
	t.Parallel()

	testEmail := "test@gmail.com"
	testPwd := "dpofdOPOOo12"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)
	gomock.InOrder(
		userMock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, nil).Times(1),
		userMock.EXPECT().CreateUser(gomock.Any()).Return(nil, nil).Times(1),
	)

	uc := NewAuthUseCase(userMock, authMock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.NoError(t, err)
}

func TestRegisterUseCase_ConflictErr(t *testing.T) {
	t.Parallel()
	testEmail := "test@gmail.com"
	testPwd := "dpofdOPOOo12"

	existingUser := entity.User{Email: "test"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)
	gomock.InOrder(
		userMock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil),
		userMock.EXPECT().CreateUser(gomock.Any()).Return(nil, nil).Times(0),
	)

	uc := NewAuthUseCase(userMock, authMock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.ErrorIs(t, err, entity.AlredyExitError)
}

func TestRegisterUseCase_ValidateErr(t *testing.T) {
	t.Parallel()
	testEmail := "wrong_email"
	testPwd := "dpofdOPOOo121"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	uc := NewAuthUseCase(userMock, authMock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.ErrorIs(t, err, entity.InvalidInput)
}
func TestLoginUseCase_OK(t *testing.T) {

	testEmail := "test@gmail.com"
	testPwd := "testgmailCom122"
	testAccessToken := "dummy.access.value"
	testRefreshToken := "dummy.refresh.value"

	config.JWTKeys.PrivateKey, _ = config.LoadPrivateKey("private.pem")

	salt, _ := pwd.GenerateSalt()
	testHashedPwd := pwd.HashPassword(testPwd, []byte(salt))
	existingUser := entity.User{
		ID:             0,
		Email:          testEmail,
		HashedPassword: testHashedPwd,
		Salt:           salt,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	uc := NewAuthUseCase(userMock, authMock)

	gomock.InOrder(
		userMock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil),
		authMock.EXPECT().CreateToken(gomock.Any()).Return(nil),
	)

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

	tokens, err := uc.LoginUseCase(testEmail, testPwd)
	require.NoError(t, err)
	require.NotEmpty(t, tokens)
	require.Equal(t, testAccessToken, tokens.AccessToken)
	require.Equal(t, testRefreshToken, tokens.RefreshToken)
}

func TestLoginUseCase_BadCredentials(t *testing.T) {
	t.Parallel()

	testEmail := "test@gmail.com"
	testPwd := "correctPwd123"
	wrongPwd := "wrongPwd12323"

	salt, _ := pwd.GenerateSalt()
	hashed := pwd.HashPassword(testPwd, []byte(salt))
	existingUser := entity.User{
		ID:             0,
		Email:          testEmail,
		HashedPassword: hashed,
		Salt:           salt,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	userMock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil)

	uc := NewAuthUseCase(userMock, authMock)
	_, err := uc.LoginUseCase(testEmail, wrongPwd)
	require.ErrorIs(t, err, entity.BadCredentials)
}

func TestLoginUseCase_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)
	gomock.InOrder(
		userMock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, entity.NotFoundError),
	)
	uc := NewAuthUseCase(userMock, authMock)
	_, err := uc.LoginUseCase("not@exists.com", "somePwd123")
	require.ErrorIs(t, err, entity.NotFoundError)
}

func TestLoginUseCase_InvalidInput(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)
	uc := NewAuthUseCase(userMock, authMock)
	_, err := uc.LoginUseCase("wrongsyntax.com", "somePwd123")
	require.ErrorIs(t, err, entity.InvalidInput)
}

func TestRefreshUseCase_OK(t *testing.T) {
	t.Parallel()

	testRefreshToken := "dummy.refresh.value"
	testAccessToken := "dummy.access.value"
	userID := 1
	tokenID := "token.id"
	storedToken := &entity.RefreshToken{
		ID:        0,
		TokenID:   tokenID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(100 * time.Minute),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	gomock.InOrder(
		authMock.EXPECT().GetToken(gomock.Any(), gomock.Any()).Return(storedToken, nil),
		authMock.EXPECT().DeleteToken(gomock.Any(), gomock.Any()).Return(nil),
		authMock.EXPECT().CreateToken(gomock.Any()).Return(nil),
	)

	patchParse := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
		c := tokens.Claims{
			Sub:  strconv.Itoa(userID),
			Type: entity.RefreshTokenType,
			RegisteredClaims: jwt.RegisteredClaims{
				ID: tokenID,
			},
		}
		return &c, nil
	})
	patchAccess := monkey.Patch(tokens.GenerateAccessToken, func(userID int, exp time.Duration) (string, error) {
		return testAccessToken, nil
	})

	patchRefresh := monkey.Patch(tokens.GenerateRefreshToken, func(tokenID string, userID int, exp time.Duration) (string, error) {
		return testRefreshToken, nil
	})

	defer func() {
		patchAccess.Unpatch()
		patchRefresh.Unpatch()
		patchParse.Unpatch()
	}()

	uc := NewAuthUseCase(userMock, authMock)

	dto, err := uc.RefreshTokenUseCase(testRefreshToken)

	require.NoError(t, err)
	require.Equal(t, dto.AccessToken, testAccessToken)
	require.Equal(t, dto.RefreshToken, testRefreshToken)
}

func TestRefreshUseCase_BadParseToken(t *testing.T) {
	t.Parallel()

	testRefreshToken := "dummy.refresh.value"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	patchParse := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
		return nil, errors.New("bad token")
	})

	defer func() {
		patchParse.Unpatch()
	}()

	uc := NewAuthUseCase(userMock, authMock)

	_, err := uc.RefreshTokenUseCase(testRefreshToken)

	require.ErrorIs(t, err, entity.JWTError)
}

func TestRefreshUseCase_InvalidType(t *testing.T) {
	t.Parallel()

	testRefreshToken := "dummy.refresh.value"
	userID := 1
	tokenID := "token.id"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	patchParse := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
		c := tokens.Claims{
			Sub:  strconv.Itoa(userID),
			Type: "invalid_type",
			RegisteredClaims: jwt.RegisteredClaims{
				ID: tokenID,
			},
		}
		return &c, nil
	})

	defer func() {
		patchParse.Unpatch()
	}()

	uc := NewAuthUseCase(userMock, authMock)

	_, err := uc.RefreshTokenUseCase(testRefreshToken)

	require.ErrorIs(t, err, entity.JWTError)
}

func TestRefreshUseCase_NotFoundToken(t *testing.T) {
	t.Parallel()

	testRefreshToken := "dummy.refresh.value"
	userID := 1
	tokenID := "token.id"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockUserRepo(ctrl)
	authMock := mocks.NewMockAuthRepo(ctrl)

	gomock.InOrder(
		authMock.EXPECT().GetToken(gomock.Any(), gomock.Any()).Return(nil, entity.NotFoundError),
	)

	patchParse := monkey.Patch(tokens.ParseToken, func(token string) (*tokens.Claims, error) {
		c := tokens.Claims{
			Sub:  strconv.Itoa(userID),
			Type: entity.RefreshTokenType,
			RegisteredClaims: jwt.RegisteredClaims{
				ID: tokenID,
			},
		}
		return &c, nil
	})

	defer func() {
		patchParse.Unpatch()
	}()

	uc := NewAuthUseCase(userMock, authMock)

	_, err := uc.RefreshTokenUseCase(testRefreshToken)

	require.ErrorIs(t, err, entity.JWTError)

}
