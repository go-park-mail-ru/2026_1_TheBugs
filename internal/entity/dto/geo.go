package dto

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"

type GeographyDTO struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GeographyPointToGeographyDTO(point geo.GeographyPoint) GeographyDTO {
	return GeographyDTO{
		Lat: point.Lat,
		Lon: point.Lon,
	}
}

func GeographyPointPtrToGeographyDTO(point *geo.GeographyPoint) *GeographyDTO {
	if point == nil {
		return nil
	}

	return &GeographyDTO{
		Lat: point.Lat,
		Lon: point.Lon,
	}
}
