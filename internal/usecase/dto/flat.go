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
	FlatCategoryID int
	Number         *int
	Floor          int
}

func PosterInputFlatDTOtoFlatInput(poster *PosterInputFlatDTO) *FlatInput {
	return &FlatInput{
		CategoryID: poster.FlatCategoryID,
		Floor:      poster.FlatFloor,
		Number:     poster.FlatNumber,
	}
}

type FlatInput struct {
	PropertyID int
	CategoryID int
	Floor      int
	Number     *int
}
