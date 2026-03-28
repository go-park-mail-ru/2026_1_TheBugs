package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type FacilityDTO struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}

func FacilitiesToFacilitiesDTO(fs []entity.Facility) []FacilityDTO {
	res := make([]FacilityDTO, 0, len(fs))
	for _, f := range fs {
		res = append(res, FacilityDTO{Alias: f.Alias, Name: f.Name})
	}
	return res

}
