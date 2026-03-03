package response

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"

type PostersResponse struct {
	Len     int             `json:"len"`
	Posters []dto.PosterDTO `json:"posters"`
}
