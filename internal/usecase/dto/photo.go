package dto

import (
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

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

type PhotoInputDTO struct {
	FileHeader *FileInput
	Order      int
}

func posterPhotosInputFlatDTOtoPhotosInput(poster *PosterInputFlatDTO) []PhotoInput {
	photos := make([]PhotoInput, 0, len(poster.Images))
	for _, photo := range poster.Images {
		var photoInput PhotoInput
		photoInput.FileHeader = photo.FileHeader
		photoInput.Order = photo.Order
		photos = append(photos, photoInput)
	}

	return photos
}

type PhotoInput struct {
	FileHeader *FileInput
	Path       string
	Order      int
}

func GeneratePhotoPathForPoster(alias string, order int) string {
	return fmt.Sprintf("/poster/img/%s/%d.jpg", alias, order)
}

func MakePhotoPathsForPoster(poster *PosterInput) {
	for i, image := range poster.Images {
		path := GeneratePhotoPathForPoster(poster.Alias, image.Order)
		poster.Images[i].Path = path
	}
}

func MakeUrlsFromPaths(poster *PosterDTO, publicHost string, bucket string) {
	for i, image := range poster.Images {
		url := photo.MakeUrlFromPath(image.ImgURL, publicHost, bucket)
		poster.Images[i].ImgURL = url
	}
}
