package user

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

type UserRepo struct {
	userSlice []entity.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		userSlice: []entity.User{},
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
	newId := 0
	if len(r.userSlice) > 0 {
		lastUser := r.userSlice[len(r.userSlice)-1]
		newId = lastUser.ID + 1
	}

	newUser := entity.User{
		ID:             newId,
		Email:          dto.Email,
		HashedPassword: dto.HashedPassword,
		Salt:           dto.Salt,
	}
	r.userSlice = append(r.userSlice, newUser)
	log.Println(r.userSlice)
	return &newUser, nil
}

//func (r *UserRepo) CreateUserRefreshToken(userID string, tokenID string)
