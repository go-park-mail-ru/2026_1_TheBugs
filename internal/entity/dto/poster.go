package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PosterDTO struct {
	Price   float64 `json:"price"`
	ImgURL  *string `json:"img_url"`
	Address string  `json:"address"`
	Metro   string  `json:"metro"`
	Area    float64 `json:"area"`
	Floor   int     `json:"floor"`
	Type    string  `json:"type"`
}

type PostersFiltersDTO struct {
	Limit  int
	Offset int
}

func PostersToPostersDTO(posters []entity.Poster) []PosterDTO {
	listPosters := make([]PosterDTO, 0, len(posters))
	for _, poster := range posters {
		posterDTO := PosterDTO{
			Price:   poster.Price,
			ImgURL:  poster.ImgURL,
			Address: poster.Address,
			Metro:   poster.Metro,
			Area:    poster.Area,
			Floor:   poster.Floor,
			Type:    poster.Type,
		}

		listPosters = append(listPosters, posterDTO)
	}
	return listPosters
}
