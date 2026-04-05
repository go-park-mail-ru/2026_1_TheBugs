package dto

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
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
	URL        *string
}

func posterPhotosInputFlatDTOtoPhotosInput(poster *PosterInputFlatDTO) []PhotoInput {
	photos := make([]PhotoInput, 0, len(poster.Images))
	for _, photo := range poster.Images {
		var photoInput PhotoInput
		photoInput.FileHeader = photo.FileHeader
		photoInput.Order = photo.Order
		photoInput.URL = photo.URL
		photos = append(photos, photoInput)
	}

	return photos
}

type PhotoInput struct {
	FileHeader *FileInput
	Path       string
	Order      int
	URL        *string
}

func GeneratePhotoPathForPoster(alias string, order int) string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	reqId := fmt.Sprintf("%016x", seed.Int())[:10]
	reqId += fmt.Sprintf("_ord_%d", order)
	return fmt.Sprintf("/poster/img/%s/%s.jpg", alias, reqId)
}

func MakePhotoPathsForPoster(poster *PosterInput) {
	for i, image := range poster.Images {
		if image.FileHeader != nil {
			path := GeneratePhotoPathForPoster(poster.Alias, image.Order)
			poster.Images[i].Path = path
			continue
		}
		url := *image.URL
		poster.Images[i].Path = strings.TrimPrefix(url, fmt.Sprintf("%s/%s", config.Config.PublicHost, config.Config.Bucket))

	}
}

func MakeUrlsFromPaths(poster *PosterDTO, publicHost string, bucket string) {
	for i, image := range poster.Images {
		url := photo.MakeUrlFromPath(image.ImgURL, publicHost, bucket)
		poster.Images[i].ImgURL = url
	}
}
