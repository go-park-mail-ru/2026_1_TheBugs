package response

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"

type PosterResponse struct {
	Poster *dto.PosterDTO `json:"poster"`
}

type MetroResponse struct {
	Len           int                   `json:"len"`
	MetroStations []dto.MetroStationDTO `json:"metro_stations"`
}

type MyPostersResponse struct {
	Len     int               `json:"len"`
	Posters []dto.MyPosterDTO `json:"posters"`
}

type CreatedPosterResponse struct {
	Poster *dto.CreatedPoster `json:"poster"`
}

type PosterViewsResponse struct {
	Views int `json:"views"`
}

type GenerateDescriptionResponse struct {
	Description string `json:"description"`
}

type PosterFavoritesCountResponse struct {
	Count int `json:"favorites"`
}

type PriceHistoryResponse struct {
	History []dto.PriceHistoryDTO `json:"history"`
	Count   int                   `json:"count"`
}
