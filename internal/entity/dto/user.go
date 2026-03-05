package dto

import (
	"time"
)

type CreateUserDTO struct {
	Email          string
	HashedPassword string
	Salt           string
}

type UserAccessCredDTO struct {
	AccessToken     string `json:"access_token"`
	AccessTokenExp  int    `json:"expire_at"`
	RefreshToken    string `json:"refresh_token"`
	RefreshTokenExp int    `json:"refresh_expire_at"`
}

type CreateRefreshTokenDTO struct {
	TokenID   string
	UserID    int
	ExpiresAt time.Time
}
