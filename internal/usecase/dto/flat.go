package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type FlatDTO struct {
	FlatCategory string `json:"flat_category"`
	Number       int    `json:"flat_number"`
	Floor        int    `json:"floor"`
	RoomCount    int    `json:"room_count"`
}

func FlatToFlatFlatDTO(flat *entity.Flat) *FlatDTO {
	if flat == nil {
		return nil
	}

	return &FlatDTO{
		FlatCategory: flat.FlatCategory,
		Number:       flat.Number,
		Floor:        flat.Floor,
		RoomCount:    flat.RoomCount,
	}
}

type FlatInputDTO struct {
	FlatCategoryID int `schema:"category_id"`
	Number         int `schema:"flat_number"`
	Floor          int `schema:"floor"`
}

func PosterInputFlatDTOtoFlatInput(poster *PosterInputFlatDTO) *entity.FlatInput {
	return &entity.FlatInput{
		CategoryID: poster.Flat.FlatCategoryID,
		Floor:      poster.Flat.Floor,
		Number:     poster.Flat.Number,
	}
}
