package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PosterDTO struct {
	ID      int      `json:"id"`
	Price   float64  `json:"price"`
	ImgURL  *string  `json:"imageUrl"`
	Address string   `json:"address"`
	Metro   *string  `json:"metro"`
	Area    float64  `json:"area"`
	Rating  *float64 `json:"rating"`
	Beds    *int     `json:"beds"`
}

type PostersFiltersDTO struct {
	Limit          int
	Offset         int
	UtilityCompany *string
}

func PostersToPostersDTO(posters []entity.Poster) []PosterDTO {
	listPosters := make([]PosterDTO, 0, len(posters))
	for i, poster := range posters {
		rating := float64((i*3/2)%10 + 1)
		count := (i%5 + 1)
		posterDTO := PosterDTO{
			ID:      poster.Id,
			Price:   poster.Price,
			ImgURL:  poster.ImgURL,
			Address: poster.Address,
			Metro:   poster.Metro,
			Area:    poster.Area,
			Rating:  &rating,
			Beds:    &count,
		}

		listPosters = append(listPosters, posterDTO)
	}
	return listPosters
}
