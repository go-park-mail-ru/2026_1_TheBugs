package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type DeveloperDTO struct {
	DeveloperID   int     `json:"developer_id"`
	DeveloperName string  `json:"developer_name"`
	AvatarURL     *string `json:"avatar_url"`
}

func DevelopersToDevelopersDTO(developer []entity.Developer) []DeveloperDTO {
	developerDTOs := make([]DeveloperDTO, len(developer))
	for i := range developer {
		developerDTOs[i] = DeveloperDTO{
			DeveloperID:   developer[i].ID,
			DeveloperName: developer[i].DeveloperName,
			AvatarURL:     developer[i].AvatarURL,
		}
	}
	return developerDTOs
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
		DeveloperID:   developer.ID,
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

func UtilityCompaniesToUtilityCompaniesDTO(companies []entity.UtilityCompanyCard) []UtilityCompanyCardDTO {
	dtos := make([]UtilityCompanyCardDTO, len(companies))
	for i, c := range companies {
		dtos[i] = UtilityCompanyCardDTO{
			ID:          c.ID,
			CompanyName: c.CompanyName,
			Alias:       c.Alias,
			AvatarURL:   c.AvatarURL,
		}
	}
	return dtos
}
