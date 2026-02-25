package auth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/pwd"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/tokens"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/validator"
)

type AuthUseCase struct {
	repo usecase.Repo
}

func NewAuthUseCase(repo usecase.Repo) *AuthUseCase {
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
		if errors.Is(err, entity.NotFoundError) {
			return entity.NotFoundError
		}
		return fmt.Errorf("uc.repo.GetUserByEmail: %w", err)
	}
	salt, err := pwd.GenerateSalt()
	if err != nil {
		return fmt.Errorf("pwd.GenerateSalt: %w", err)
	}
	hashedPwd := pwd.HashPassword(password, []byte(salt))
	_, err = uc.repo.CreateUser(dto.CreateUserDTO{
		Email:          email,
		HashedPassword: hashedPwd,
		Salt:           salt,
	})
	if err != nil {
		return fmt.Errorf("uc.repo.CreateUser: %w", err)
	}
	return nil
}

func (uc AuthUseCase) LoginUseCase(email string, passwod string) (*dto.UserAccessCredDTO, error) {
	var cred dto.UserAccessCredDTO
	if !validator.ValidateEmail(email) || !validator.ValidatePwd(passwod) {
		return &cred, entity.InvalidInput
	}
	user, err := uc.repo.GetUserByEmail(email)
	if err != nil {
		return &cred, entity.NotFoundError
	}
	ok := pwd.VerifyPassword(passwod, []byte(user.Salt), user.HashedPassword)
	if !ok {
		return &cred, entity.BadCredentials
	}
	accessToken, err := tokens.GenerateJWT(strconv.Itoa(user.ID), "access", config.Config.JWT.AccessExp)
	if err != nil {
		return &cred, entity.ServiceError
	}
	cred.AccessToken = accessToken
	cred.AccessTokenExp = int(config.Config.JWT.AccessExp.Seconds())
	return &cred, nil
}
