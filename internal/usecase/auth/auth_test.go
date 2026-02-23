package auth

import (
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type MockUserRepo struct {
	Slice []entity.User
}

func NewMockUserRepo(slice []entity.User) *MockUserRepo {
	return &MockUserRepo{
		Slice: slice,
	}
}

func (r *MockUserRepo) GetUserByEmail(email string) (entity.User, error) {
	for _, u := range r.Slice {
		if u.Email == email {
			return u, nil
		}
	}
	return entity.User{}, repository.NotFoundRecord

}

func (r *MockUserRepo) CreateUser(dto dto.CreateUserDTO) (entity.User, error) {
	id := uuid.New()
	newUser := entity.User{
		Id:             id.String(),
		Email:          dto.Email,
		HashedPassword: dto.HashedPassword,
		Satl:           dto.Salt,
	}
	r.Slice = append(r.Slice, newUser)
	return newUser, nil
}

func TestRegisterUseCaseOK(t *testing.T) {
	t.Parallel()
	testEmail := "email"
	testPwd := "pwd"
	mock := NewMockUserRepo(make([]entity.User, 0))

	uc := NewAuthUseCase(mock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.NoError(t, err)

	u, err := mock.GetUserByEmail(testEmail)
	require.NoError(t, err)

	require.Equal(t, testEmail, u.Email)
}

func TestRegisterUseCaseErr(t *testing.T) {
	t.Parallel()
	testEmail := "email"
	testPwd := "pwd"
	mock := NewMockUserRepo([]entity.User{{Email: testEmail}})
	uc := NewAuthUseCase(mock)
	err := uc.RegisterUseCase(testEmail, testPwd)
	require.ErrorIs(t, err, usecase.AlredyExitError)
}
