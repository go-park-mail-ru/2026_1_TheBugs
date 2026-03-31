package dto

import (
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

type CreateUserDTO struct {
	Email          string
	HashedPassword *string
	Salt           *string
	Password       string
	Phone          string
	FirstName      string
	LastName       string
}

type CreateUserByProviderDTO struct {
	Email      string
	Provider   entity.ProviderType
	Phone      string
	FirstName  string
	LastName   string
	ProviderID *string
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

type LogoutDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserDTO struct {
	ID        int     `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func UserToDTO(user *entity.UserDetails) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
	}
}
