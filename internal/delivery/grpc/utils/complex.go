package utils

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

func MapComplex(cmp *dto.UtilityCompanyDTO) *complex.GetComplexResponse {
	if cmp == nil {
		return nil
	}
	return &complex.GetComplexResponse{
		Id:          int32(cmp.ID),
		Alias:       cmp.Alias,
		CompanyName: cmp.CompanyName,
		Description: cmp.Description,
		Address:     cmp.Address,
		AvatarUrl:   cmp.AvatarURL,
		Geo: &complex.GeoPoint{
			Lat: cmp.GEO.Lat,
			Lon: cmp.GEO.Lon,
		},
		Photos:    mapPhotos(cmp.Photos),
		Developer: MapDeveloper(cmp.Developer),
		Phone:     cmp.Phone,
	}
}

func MapDevelopers(devs []dto.DeveloperDTO) []*complex.Developer {
	res := make([]*complex.Developer, 0, len(devs))

	for _, d := range devs {
		res = append(res, MapDeveloper(d))
	}

	return res
}

func mapPhotos(photos []dto.PhotoDTO) []*complex.Photo {
	res := make([]*complex.Photo, 0, len(photos))

	for _, p := range photos {
		res = append(res, &complex.Photo{
			Url:   p.ImgURL,
			Order: int32(p.Order),
		})
	}

	return res
}

func MapDeveloper(d dto.DeveloperDTO) *complex.Developer {
	return &complex.Developer{
		Id:            int32(d.DeveloperID),
		DeveloperName: d.DeveloperName,
		AvatarUrl:     d.AvatarURL,
	}
}

func MapComplexes(comps []dto.UtilityCompanyCardDTO) []*complex.Complex {
	res := make([]*complex.Complex, 0, len(comps))

	for _, c := range comps {
		res = append(res, mapComplexItem(c))
	}

	return res
}

func mapComplexItem(c dto.UtilityCompanyCardDTO) *complex.Complex {
	return &complex.Complex{
		Id:          int32(c.ID),
		Alias:       c.Alias,
		CompanyName: c.CompanyName,
		AvatarUrl:   c.AvatarURL,
	}
}
