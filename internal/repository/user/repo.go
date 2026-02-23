package user

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/google/uuid"
)

type Repo interface {
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(dto dto.CreateUserDTO) (*entity.User, error)
}

type UserRepo struct {
	userSlice map[string]entity.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		userSlice: map[string]entity.User{},
	}
}

func (r *UserRepo) GetUserByEmail(email string) (*entity.User, error) {
	for _, u := range r.userSlice {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, entity.NotFoundError

}

func (r *UserRepo) CreateUser(dto dto.CreateUserDTO) (*entity.User, error) {
	id := uuid.New()
	newUser := entity.User{
		Id:             id.String(),
		Email:          dto.Email,
		HashedPassword: dto.HashedPassword,
		Satl:           dto.Salt,
	}
	r.userSlice[id.String()] = newUser
	log.Println(r.userSlice)
	return &newUser, nil
}
