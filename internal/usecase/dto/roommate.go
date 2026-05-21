package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type RoommateUserDTO struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type PosterRoommatesResponse struct {
	Users []RoommateUserDTO `json:"users"`
	Len   int               `json:"len"`
}

type RoommateTagDTO struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type RoommateUserProfileDTO struct {
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	AvatarURL   *string          `json:"avatar_url,omitempty"`
	Gender      string           `json:"gender"`
	Birthday    string           `json:"birthday"`
	Description *string          `json:"description,omitempty"`
	Tags        []RoommateTagDTO `json:"tags"`
}

func RoommateUserToDTO(user *entity.RoommateUser, tags []entity.RoommateTag) *RoommateUserProfileDTO {
	dtoTags := make([]RoommateTagDTO, 0, len(tags))
	for _, tag := range tags {
		dtoTags = append(dtoTags, RoommateTagDTO{
			Name:  tag.Name,
			Alias: tag.Alias,
		})
	}

	return &RoommateUserProfileDTO{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AvatarURL:   user.AvatarURL,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		Description: user.Description,
		Tags:        dtoTags,
	}
}

type RoommateContactsDTO struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type CreateRoommateFormRequest struct {
	UserID      int      `json:"-"`
	Gender      string   `json:"gender"`
	Birthday    string   `json:"birthday"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type RoommateFormDTO struct {
	Gender      string   `json:"gender"`
	Birthday    string   `json:"birthday"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func RoommateFormToDTO(form *entity.RoommateForm, tags []string) *RoommateFormDTO {
	return &RoommateFormDTO{
		Gender:      form.Gender,
		Birthday:    form.Birthday,
		Description: form.Description,
		Tags:        tags,
	}
}

type RoommateMatchesResponse struct {
	Users []RoommateUserDTO `json:"users"`
	Len   int               `json:"len"`
}
