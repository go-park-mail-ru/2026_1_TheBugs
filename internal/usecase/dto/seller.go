package dto

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type PosterSellerDTO struct {
	SellerFirstName string  `json:"first_name"`
	SellerLastName  string  `json:"last_name"`
	SellerPhone     string  `json:"phone"`
	SellerAvatarURL *string `json:"avatar_url"`
}

func posterToPosterSellerDTO(poster *entity.PosterById) PosterSellerDTO {
	var url *string
	if poster.SellerAvatarURL != nil {
		avatar := photo.MakeUrlFromPath(*poster.SellerAvatarURL, config.Config.PublicHost, config.Config.Bucket)
		url = &avatar
	}

	return PosterSellerDTO{
		SellerFirstName: poster.SellerFirstName,
		SellerLastName:  poster.SellerLastName,
		SellerPhone:     poster.Phone,
		SellerAvatarURL: url,
	}
}
