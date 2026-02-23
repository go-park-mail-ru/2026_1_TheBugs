package auth

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
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
	_, err := uc.repo.GetUserByEmail(email)
	if err == nil {
		return usecase.AlredyExitError
	}
	salt, err := pwd.GenerateSalt()
	if err != nil {
		return usecase.ServiceError
	}
	hashedPwd := pwd.HashPassword(password, []byte(salt))
	log.Println(hashedPwd)
	u, err := uc.repo.CreateUser(dto.CreateUserDTO{
		Email:          email,
		HashedPassword: hashedPwd,
		Salt:           salt,
	})
	log.Println(u)
	if err != nil {
		return usecase.AlredyExitError
	}
	return nil
}
