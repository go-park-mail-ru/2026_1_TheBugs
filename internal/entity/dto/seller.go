package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PosterSellerDTO struct {
	SellerFirstName string `json:"seller_first_name"`
	SellerLastName  string `json:"seller_last_name"`
	SellerPhone     string `json:"seller_phone"`
}

func posterToPosterSellerDTO(poster entity.PosterById) PosterSellerDTO {
	return PosterSellerDTO{
		SellerFirstName: poster.SellerFirstName,
		SellerLastName:  poster.SellerLastName,
		SellerPhone:     poster.Phone,
	}
}
