package complex

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

type UtilityCompanyDTO struct {
	ID          int              `json:"id"`
	Phone       string           `json:"phone"`
	CompanyName string           `json:"company_name"`
	GEO         dto.GeographyDTO `json:"geo"`
	Address     string           `json:"address"`
	AvatarURL   *string          `json:"avatar_url"`
	Alias       string           `json:"alias"`
	Photos      []dto.PhotoDTO   `json:"photos"`
}

func ToUtilityCompanyDTO(complex *entity.UtilityCompany, photos []entity.UtilityCompanyPhoto) *UtilityCompanyDTO {
	photoDTOs := make([]dto.PhotoDTO, len(photos))
	for i, photo := range photos {
		photoDTOs[i] = dto.PhotoDTO{
			ImgURL: *photo.ImgURL,
			Order:  *photo.Order,
		}
	}
	return &UtilityCompanyDTO{
		ID:          complex.ID,
		Phone:       complex.Phone,
		CompanyName: complex.CompanyName,
		GEO:         dto.GeographyDTO{Lat: complex.GEO.Lat, Lon: complex.GEO.Lon},
		Address:     complex.Address,
		AvatarURL:   complex.AvatarURL,
		Photos:      photoDTOs,
		Alias:       complex.Alias,
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
