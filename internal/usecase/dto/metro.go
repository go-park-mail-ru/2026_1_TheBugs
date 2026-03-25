package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type MetroStationDTO struct {
	StationName string       `json:"station_name"`
	StationGeo  GeographyDTO `json:"station_geo"`
}

func MetroToMetroStationDTO(metro entity.MetroStation) MetroStationDTO {
	return MetroStationDTO{
		StationName: metro.StationName,
		StationGeo:  GeographyDTO(metro.StationGEO),
	}
}
