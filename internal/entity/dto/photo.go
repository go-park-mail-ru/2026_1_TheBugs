package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PhotoDTO struct {
	ImgURL string `json:"img_url"`
	Order  int    `json:"order"`
}

func posterImagesToPosterImagesDTO(poster *entity.PosterById) []PhotoDTO {
	images := make([]PhotoDTO, 0, len(poster.Images))
	for _, image := range poster.Images {
		var imageDTO PhotoDTO
		imageDTO.ImgURL = image.ImgURL
		imageDTO.Order = image.Order
		images = append(images, imageDTO)
	}

	return images
}
