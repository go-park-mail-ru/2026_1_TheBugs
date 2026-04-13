package dto

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
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

type UpdateProfileDTO struct {
	ID         int
	Phone      *string
	FirstName  *string
	LastName   *string
	AvatarPath *string
}

type UpdateProfileRequest struct {
	ID        int        `schema:"-"`
	Phone     *string    `schema:"phone"`
	FirstName *string    `schema:"first_name"`
	LastName  *string    `schema:"last_name"`
	Avatar    *FileInput `schema:"-"`
}

type UpdateUserDTO struct {
	ID             int
	Email          *string
	HashedPassword *string
	Salt           *string
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
	Phone     string  `json:"phone"`
}

func UserToDTO(user *entity.UserDetails) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarURL: user.AvatarURL,
		Phone:     user.Phone,
	}
}

func (user *UserDTO) MakeAvatarPath() {
	if user == nil {
		return
	}
	if user.AvatarURL != nil {
		avatar := photo.MakeUrlFromPath(*user.AvatarURL, config.Config.PublicHost, config.Config.Bucket)
		user.AvatarURL = &avatar
	}
}

func GenerateAvatarPathForUser(id int) string {
	saltBytes := make([]byte, 10)
	rand.Read(saltBytes)
	return fmt.Sprintf("/user/%d/avatar-%s.jpg", id, base64.RawStdEncoding.EncodeToString(saltBytes))
}
