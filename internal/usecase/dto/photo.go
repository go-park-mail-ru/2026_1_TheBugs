package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PhotoDTO struct {
	ImgURL string `json:"img_url"`
	Order  int    `json:"order"`
}

func posterImagesToPosterImagesDTO(imgs []entity.PosterImage) []PhotoDTO {
	images := make([]PhotoDTO, 0, len(imgs))
	for _, image := range imgs {
		var imageDTO PhotoDTO
		imageDTO.ImgURL = image.ImgURL
		imageDTO.Order = image.Order
		images = append(images, imageDTO)
	}

	return images
}
