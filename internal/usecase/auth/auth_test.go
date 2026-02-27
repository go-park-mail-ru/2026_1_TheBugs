package auth

// import (
// 	"testing"
// 	"time"

// 	"bou.ke/monkey"

// 	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
// 	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
// 	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
// 	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
// 	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/tokens"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )

// func TestRegisterUseCase_OK(t *testing.T) {
// 	t.Parallel()

// 	testEmail := "test@gmail.com"
// 	testPwd := "dpofdOPOOo12"

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mock := mocks.NewMockRepo(ctrl)
// 	gomock.InOrder(
// 		mock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, nil).Times(1),
// 		mock.EXPECT().CreateUser(gomock.Any()).Return(nil, nil).Times(1),
// 	)

// 	uc := NewAuthUseCase(mock)
// 	err := uc.RegisterUseCase(testEmail, testPwd)
// 	require.NoError(t, err)
// }

// func TestRegisterUseCase_ConflictErr(t *testing.T) {
// 	t.Parallel()
// 	testEmail := "test@gmail.com"
// 	testPwd := "dpofdOPOOo12"

// 	existingUser := entity.User{Email: "test"}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mock := mocks.NewMockRepo(ctrl)
// 	gomock.InOrder(
// 		mock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil),
// 		mock.EXPECT().CreateUser(gomock.Any()).Return(nil, nil).Times(0),
// 	)

// 	uc := NewAuthUseCase(mock)
// 	err := uc.RegisterUseCase(testEmail, testPwd)
// 	require.ErrorIs(t, err, entity.AlredyExitError)
// }

// func TestRegisterUseCase_ValidateErr(t *testing.T) {
// 	t.Parallel()
// 	testEmail := "wrong_email"
// 	testPwd := "dpofdOPOOo121"

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mock := mocks.NewMockRepo(ctrl)

// 	uc := NewAuthUseCase(mock)
// 	err := uc.RegisterUseCase(testEmail, testPwd)
// 	require.ErrorIs(t, err, entity.InvalidInput)
// }
// func TestLoginUseCase_OK(t *testing.T) {

// 	testEmail := "test@gmail.com"
// 	testPwd := "testgmailCom122"
// 	testToken := "dummy.token.value"

// 	config.JWTKeys.PrivateKey, _ = config.LoadPrivateKey("private.pem")

// 	salt, _ := pwd.GenerateSalt()
// 	testHashedPwd := pwd.HashPassword(testPwd, []byte(salt))
// 	existingUser := entity.User{
// 		ID:             0,
// 		Email:          testEmail,
// 		HashedPassword: testHashedPwd,
// 		Salt:           salt,
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mock := mocks.NewMockRepo(ctrl)

// 	mock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil)

// 	patch := monkey.Patch(tokens.GenerateAccessToken, func(userID int, exp time.Duration) (string, error) {
// 		return testToken, nil
// 	})
// 	defer patch.Unpatch()

// 	uc := NewAuthUseCase(mock)
// 	tokens, err := uc.LoginUseCase(testEmail, testPwd)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, tokens)
// 	require.Equal(t, testToken, tokens.AccessToken)
// }

// func TestLoginUseCase_BadCredentials(t *testing.T) {
// 	t.Parallel()

// 	testEmail := "test@gmail.com"
// 	testPwd := "correctPwd123"
// 	wrongPwd := "wrongPwd12323"

// 	salt, _ := pwd.GenerateSalt()
// 	hashed := pwd.HashPassword(testPwd, []byte(salt))
// 	existingUser := entity.User{
// 		ID:             0,
// 		Email:          testEmail,
// 		HashedPassword: hashed,
// 		Salt:           salt,
// 	}

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mock := mocks.NewMockRepo(ctrl)
// 	mock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil)

// 	uc := NewAuthUseCase(mock)
// 	_, err := uc.LoginUseCase(testEmail, wrongPwd)
// 	require.ErrorIs(t, err, entity.BadCredentials)
// }

// func TestLoginUseCase_NotFound(t *testing.T) {
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mock := mocks.NewMockRepo(ctrl)
// 	gomock.InOrder(
// 		mock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, entity.NotFoundError),
// 	)
// 	uc := NewAuthUseCase(mock)
// 	_, err := uc.LoginUseCase("not@exists.com", "somePwd123")
// 	require.ErrorIs(t, err, entity.NotFoundError)
// }

// func TestLoginUseCase_InvalidInput(t *testing.T) {
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	mock := mocks.NewMockRepo(ctrl)
// 	uc := NewAuthUseCase(mock)
// 	_, err := uc.LoginUseCase("wrongsyntax.com", "somePwd123")
// 	require.ErrorIs(t, err, entity.InvalidInput)
// }
