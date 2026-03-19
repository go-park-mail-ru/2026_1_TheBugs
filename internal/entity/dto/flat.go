package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type FlatDTO struct {
	FlatCategory string `json:"flat_catigory"`
	Number       int    `json:"flat_number"`
	Floor        int    `json:"floor"`
}

func FlatToFlatFlatDTO(flat *entity.Flat) *FlatDTO {
	if flat == nil {
		return nil
	}

	return &FlatDTO{
		FlatCategory: flat.FlatCategory,
		Number:       flat.Number,
		Floor:        flat.Floor,
	}
}
