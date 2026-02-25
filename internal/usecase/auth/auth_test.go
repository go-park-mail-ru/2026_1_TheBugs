package auth

import (
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUseCaseOK(t *testing.T) {
	t.Parallel()

	testEmail := "test@gmail.com"
	testPwd := "dpofdOPOOo12"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockRepo(ctrl)
	gomock.InOrder(
		mock.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, nil).Times(1),
		mock.EXPECT().CreateUser(gomock.Any()).Return(nil, nil).Times(1),
	)

	uc := NewAuthUseCase(mock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.NoError(t, err)
}

func TestRegisterUseCaseErr(t *testing.T) {
	t.Parallel()
	testEmail := "test@gmail.com"
	testPwd := "dpofdOPOOo12"

	existingUser := entity.User{Email: "test"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockRepo(ctrl)
	gomock.InOrder(
		mock.EXPECT().GetUserByEmail(gomock.Any()).Return(&existingUser, nil),
		mock.EXPECT().CreateUser(gomock.Any()).Return(nil, nil).Times(0),
	)

	uc := NewAuthUseCase(mock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.ErrorIs(t, err, entity.AlredyExitError)
}

func TestRegisterUseCaseValidateErr(t *testing.T) {
	t.Parallel()
	testEmail := "wrong_email"
	testPwd := "dpofdOPOOo121"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockRepo(ctrl)

	uc := NewAuthUseCase(mock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.ErrorIs(t, err, entity.InvalidInput)
}
