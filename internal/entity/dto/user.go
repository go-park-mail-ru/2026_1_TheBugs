package dto

type CreateUserDTO struct {
	Email          string
	HashedPassword string
	Salt           string
}
