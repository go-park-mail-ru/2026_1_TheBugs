package usecase

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

type UserRepo interface {
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(dto dto.CreateUserDTO) (*entity.User, error)
}

type AuthRepo interface {
	GetToken(tokenID string, userID int) (*entity.RefreshToken, error)
	CreateToken(dto dto.CreateRefreshTokenDTO) error
	DeleteToken(tokenID string, userID int) error
}
