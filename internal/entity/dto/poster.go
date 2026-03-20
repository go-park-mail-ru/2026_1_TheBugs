package dto

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

type PosterCardDTO struct {
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
	Limit  int
	Offset int
}

func PostersToPostersDTO(posters []entity.Poster) []PosterCardDTO {
	listPosters := make([]PosterCardDTO, 0, len(posters))
	for i, poster := range posters {
		rating := float64((i*3/2)%10 + 1)
		count := (i%5 + 1)
		posterDTO := PosterCardDTO{
			ID:      poster.ID,
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

type PosterDTO struct {
	ID          int             `json:"id"`
	Alias       string          `json:"alias"`
	Price       float64         `json:"price"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Area        float64         `json:"area"`
	Geo         GeographyDTO    `json:"building_geo"`
	Address     string          `json:"address"`
	District    *string         `json:"district"`
	Metro       *string         `json:"metro"`
	MetroGeo    *GeographyDTO   `json:"metro_geo"`
	City        string          `json:"city"`
	FloorCount  int             `json:"floor_count"`
	Images      []PhotoDTO      `json:"images"`
	Seller      PosterSellerDTO `json:"seller"`
	Flat        *FlatDTO        `json:"flat,omitempty"`
	House       *HouseDTO       `json:"house,omitempty"`

	Company *UtilityCompanyCardDTO `json:"company,omitempty"`
}

func PosterToPosterDTO(poster entity.PosterById) PosterDTO {
	return PosterDTO{
		ID:          poster.ID,
		Alias:       poster.Alias,
		Price:       poster.Price,
		Category:    poster.Category,
		Description: poster.Description,
		Area:        poster.Area,
		Geo:         GeographyPointToGeographyDTO(poster.Geo),
		Address:     poster.Address,
		District:    poster.District,
		Metro:       poster.Metro,
		MetroGeo:    GeographyPointPtrToGeographyDTO(poster.MetroGeo),
		City:        poster.City,
		FloorCount:  poster.FloorCount,
		Images:      posterImagesToPosterImagesDTO(poster),
		Seller:      posterToPosterSellerDTO(poster),
		Company:     posterToUtilityCompanyCardDTO(poster),
	}
}
