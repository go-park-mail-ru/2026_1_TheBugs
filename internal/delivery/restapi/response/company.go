package response

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

type DevelopersResponse struct {
	Len        int                `json:"len"`
	Developers []dto.DeveloperDTO `json:"developers"`
}

type CompaniesResponse struct {
	Len       int                         `json:"len"`
	Companies []dto.UtilityCompanyCardDTO `json:"utility_companies"`
}

func MapGRPCComplexToDTO(grpcComplex *complex.GetComplexResponse) *dto.UtilityCompanyDTO {
	return &dto.UtilityCompanyDTO{
		ID:          int(grpcComplex.Id),
		Alias:       grpcComplex.Alias,
		CompanyName: grpcComplex.CompanyName,
		Description: grpcComplex.Description,
		Address:     grpcComplex.Address,
		AvatarURL:   grpcComplex.AvatarUrl,
		GEO: dto.GeographyDTO{
			Lat: grpcComplex.Geo.Lat,
			Lon: grpcComplex.Geo.Lon,
		},
		Photos:    mapGRPCPhotosToDTO(grpcComplex.Photos),
		Developer: MapGRPCDeveloperToDTO(grpcComplex.Developer),
	}
}

func MapGRPCDeveloperToDTO(grpcDev *complex.Developer) dto.DeveloperDTO {
	if grpcDev == nil {
		return dto.DeveloperDTO{}
	}
	return dto.DeveloperDTO{
		DeveloperID:   int(grpcDev.Id),
		DeveloperName: grpcDev.DeveloperName,
		AvatarURL:     grpcDev.AvatarUrl,
	}
}

// Преобразование gRPC Photos в HTTP DTO
func mapGRPCPhotosToDTO(grpcPhotos []*complex.Photo) []dto.PhotoDTO {
	res := make([]dto.PhotoDTO, 0, len(grpcPhotos))
	for _, p := range grpcPhotos {
		res = append(res, dto.PhotoDTO{
			ImgURL: p.Url,
			Order:  int(p.Order),
		})
	}
	return res
}

// Преобразование gRPC DevelopersResponse в HTTP DTO
func MapGRPCDevelopersToDTO(grpcDevs *complex.GetAllDevelopersResponse) *DevelopersResponse {
	res := make([]dto.DeveloperDTO, 0, len(grpcDevs.Developers))
	for _, dev := range grpcDevs.Developers {
		dtoDev := MapGRPCDeveloperToDTO(dev)
		if &dtoDev != nil {
			res = append(res, dtoDev)
		}
	}
	return &DevelopersResponse{
		Developers: res,
	}
}

func MapGRPCComplexesToDTO(grpcComplexes []*complex.Complex) *CompaniesResponse {
	res := make([]dto.UtilityCompanyCardDTO, 0, len(grpcComplexes))
	for _, grpcComp := range grpcComplexes {
		res = append(res, mapGRPCComplexItemToDTO(grpcComp))
	}
	return &CompaniesResponse{
		Len:       len(res),
		Companies: res,
	}
}

func mapGRPCComplexItemToDTO(grpcComp *complex.Complex) dto.UtilityCompanyCardDTO {
	return dto.UtilityCompanyCardDTO{
		ID:          int(grpcComp.Id),
		Alias:       grpcComp.Alias,
		CompanyName: grpcComp.CompanyName,
		AvatarURL:   grpcComp.AvatarUrl,
	}
}
