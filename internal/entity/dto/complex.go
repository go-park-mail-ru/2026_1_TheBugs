package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type DeveloperDTO struct {
	DeveloperName string  `json:"developer_name"`
	AvatarURL     *string `json:"avatar_url"`
}

type UtilityCompanyDTO struct {
	ID          int          `json:"id"`
	Phone       string       `json:"phone"`
	CompanyName string       `json:"company_name"`
	Description string       `json:"description"`
	GEO         GeographyDTO `json:"geo"`
	Address     string       `json:"address"`
	AvatarURL   *string      `json:"avatar_url"`
	Alias       string       `json:"alias"`
	Photos      []PhotoDTO   `json:"photos"`
	Developer   DeveloperDTO `json:"developer"`
}

func ToUtilityCompanyDTO(complex *entity.UtilityCompany, photos []entity.UtilityCompanyPhoto, developer *entity.Developer) *UtilityCompanyDTO {
	photoDTOs := make([]PhotoDTO, len(photos))
	for i, photo := range photos {
		photoDTOs[i] = PhotoDTO{
			ImgURL: *photo.ImgURL,
			Order:  *photo.Order,
		}
	}
	developerDTO := DeveloperDTO{
		DeveloperName: developer.DeveloperName,
		AvatarURL:     developer.AvatarURL,
	}
	return &UtilityCompanyDTO{
		ID:          complex.ID,
		Phone:       complex.Phone,
		CompanyName: complex.CompanyName,
		Description: complex.Description,
		GEO:         GeographyDTO{Lat: complex.GEO.Lat, Lon: complex.GEO.Lon},
		Address:     complex.Address,
		AvatarURL:   complex.AvatarURL,
		Photos:      photoDTOs,
		Alias:       complex.Alias,
		Developer:   developerDTO,
	}
}

type UtilityCompanyCardDTO struct {
	ID          int     `json:"id"`
	CompanyName string  `json:"company_name"`
	AvatarURL   *string `json:"avatar_url"`
	Alias       string  `json:"alias"`
}

func posterToUtilityCompanyCardDTO(poster *entity.PosterById) *UtilityCompanyCardDTO {
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
