package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type UtilityCompanyCardDTO struct {
	ID          int     `json:"id"`
	CompanyName string  `json:"company_name"`
	AvatarURL   *string `json:"avatar_url"`
	Alias       string  `json:"alias"`
}

func posterToUtilityCompanyCardDTO(poster entity.PosterById) *UtilityCompanyCardDTO {
	if poster.CompanyID == nil {
		return nil
	}
	return &UtilityCompanyCardDTO{
		ID:          *poster.CompanyID,
		CompanyName: *poster.CompanyName,
		AvatarURL:   poster.CompanyAvatarURL,
		Alias:       *poster.CompanyAlias,
	}
}
