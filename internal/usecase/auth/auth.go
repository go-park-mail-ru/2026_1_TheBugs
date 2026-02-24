package auth

import (
	"errors"

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
		if !errors.Is(err, entity.NotFoundError) {
			return entity.ServiceError
		}

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

func (uc AuthUseCase) LoginUseCase(email string, passwod string) (*dto.UserAccessCredDTO, error) {
	var cred dto.UserAccessCredDTO
	if !validator.ValidateEmail(email) || !validator.ValidatePwd(passwod) {
		return &cred, entity.InvalidInput
	}
	user, err := uc.repo.GetUserByEmail(email)
	if err != nil {
		return &cred, entity.NotFoundError
	}
	ok := pwd.VerifyPassword(passwod, []byte(user.Satl), user.HashedPassword)
	if !ok {
		return &cred, entity.BadCredentials
	}
	accessToken, err := tokens.GenerateJWT(user.Id, "access", config.Config.JWT.AccessExp)
	if err != nil {
		return &cred, entity.ServiceError
	}
	cred.AccessToken = accessToken
	cred.AccessTokenExp = int(config.Config.JWT.AccessExp.Seconds())
	return &cred, nil
}
