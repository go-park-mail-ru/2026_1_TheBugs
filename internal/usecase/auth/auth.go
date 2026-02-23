package auth

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/validator"
)

type AuthUseCase struct {
	repo user.Repo
}

func NewAuthUseCase(repo user.Repo) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
	}
}
func (uc AuthUseCase) RegisterUseCase(email string, password string) error {
	if !validator.ValidateEmail(email) || !validator.ValidatePwd(password) {
		return entity.InvalidInput
	}
	existing, err := uc.repo.GetUserByEmail(email)
	if existing != nil {
		return entity.AlredyExitError
	}
	if err != nil {
		return entity.ServiceError
	}
	salt, err := pwd.GenerateSalt()
	if err != nil {
		return entity.ServiceError
	}
	hashedPwd := pwd.HashPassword(password, []byte(salt))
	_, err = uc.repo.CreateUser(dto.CreateUserDTO{
		Email:          email,
		HashedPassword: hashedPwd,
		Salt:           salt,
	})
	if err != nil {
		return entity.AlredyExitError
	}
	return nil
}

func (uc AuthUseCase) LoginUseCase(email string, passwod string) error {
	if !validator.ValidateEmail(email) || !validator.ValidatePwd(passwod) {
		return entity.InvalidInput
	}
	user, err := uc.repo.GetUserByEmail(email)
	if err != nil {
		return entity.NotFoundError
	}
	ok := pwd.VerifyPassword(passwod, []byte(user.Satl), user.HashedPassword)
	if !ok {
		return entity.BadCredentials
	}
	return nil
}
