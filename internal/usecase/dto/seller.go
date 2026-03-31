package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PosterSellerDTO struct {
	SellerFirstName string  `json:"first_name"`
	SellerLastName  string  `json:"last_name"`
	SellerPhone     string  `json:"phone"`
	SellerAvatarURL *string `json:"avatar_url"`
}

func posterToPosterSellerDTO(poster *entity.PosterById) PosterSellerDTO {
	return PosterSellerDTO{
		SellerFirstName: poster.SellerFirstName,
		SellerLastName:  poster.SellerLastName,
		SellerPhone:     poster.Phone,
		SellerAvatarURL: poster.SellerAvatarURL,
	}
}
