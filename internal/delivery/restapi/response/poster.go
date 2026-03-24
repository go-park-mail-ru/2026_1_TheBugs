package response

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"

type PostersResponse struct {
	Len     int                 `json:"len"`
	Posters []dto.PosterCardDTO `json:"posters"`
}

type PosterResponse struct {
	Poster *dto.PosterDTO `json:"poster"`
}

type MetroResponse struct {
	Len           int                   `json:"len"`
	MetroStations []dto.MetroStationDTO `json:"metro_stations"`
}
