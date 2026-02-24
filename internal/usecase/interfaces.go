package usecase

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

type Repo interface {
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(dto dto.CreateUserDTO) (*entity.User, error)
}
