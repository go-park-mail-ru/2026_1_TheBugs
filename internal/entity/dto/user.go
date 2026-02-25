package dto

type CreateUserDTO struct {
	Email          string
	HashedPassword string
	Salt           string
}

type UserAccessCredDTO struct {
	AccessToken    string `json:"access_token"`
	AccessTokenExp int    `json:"expire_at"`
}
