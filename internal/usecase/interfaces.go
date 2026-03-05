package usecase

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

//go:generate mockgen -source=interfaces.go -destination=../mocks/mocks_repo.go -package=mocks

type UserRepo interface {
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(dto dto.CreateUserDTO) (*entity.User, error)
}

type PosterRepo interface {
	GetPosters(limit, offset int) ([]*entity.Poster, error)
	CountPosters() int
}
type AuthRepo interface {
	GetToken(tokenID string, userID int) (*entity.RefreshToken, error)
	CreateToken(dto dto.CreateRefreshTokenDTO) error
	DeleteToken(tokenID string, userID int) error
}
