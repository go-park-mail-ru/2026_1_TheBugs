package photo

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

func MakePhotoPathsForPoster(poster *entity.PosterInput) {
	for i, image := range poster.Images {
		path := fmt.Sprintf("/poster/img/%s/%d.jpg", poster.Alias, image.Order)
		poster.Images[i].Path = path
	}
}

func GetKeyFromPath(path string) string {
	return strings.TrimPrefix(path, "/")
}

func MakeUrlsFromPaths(poster *dto.PosterDTO, publicHost string, bucket string) {
	for i, image := range poster.Images {
		url := fmt.Sprintf("%s/%s%s", publicHost, bucket, image.ImgURL)
		poster.Images[i].ImgURL = url
	}
}
