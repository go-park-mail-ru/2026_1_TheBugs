package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type UtilityCompanyDTO struct {
	ID          int          `json:"id"`
	Phone       string       `json:"phone"`
	CompanyName string       `json:"company_name"`
	GEO         GeographyDTO `json:"geo"`
	Address     string       `json:"address"`
	AvatarURL   *string      `json:"avatar_url"`
	Alias       string       `json:"alias"`
	Photos      []PhotoDTO   `json:"photos"`
}

func ToUtilityCompanyDTO(complex *entity.UtilityCompany, photos []entity.UtilityCompanyPhoto) *UtilityCompanyDTO {
	photoDTOs := make([]PhotoDTO, len(photos))
	for i, photo := range photos {
		photoDTOs[i] = PhotoDTO{
			ImgURL: *photo.ImgURL,
			Order:  *photo.Order,
		}
	}
	return &UtilityCompanyDTO{
		ID:          complex.ID,
		Phone:       complex.Phone,
		CompanyName: complex.CompanyName,
		GEO:         GeographyDTO{Lat: complex.GEO.Lat, Lon: complex.GEO.Lon},
		Address:     complex.Address,
		AvatarURL:   complex.AvatarURL,
		Photos:      photoDTOs,
		Alias:       complex.Alias,
	}
}
