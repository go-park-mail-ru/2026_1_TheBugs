package mapper

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PosterCardDTOToProto(p *dto.PosterCardDTO) *poster.PosterCard {
	if p == nil {
		return nil
	}

	return &poster.PosterCard{
		Id:           int64(p.ID),
		Alias:        p.Alias,
		Price:        p.Price,
		ImageUrl:     p.ImgURL,
		Address:      p.Address,
		Metro:        p.Metro,
		Area:         p.Area,
		FlatCategory: p.FlatCategory,
	}
}

func PosterCardsDTOToProto(list []dto.PosterCardDTO) []*poster.PosterCard {
	result := make([]*poster.PosterCard, 0, len(list))

	for i := range list {
		result = append(result, PosterCardDTOToProto(&list[i]))
	}

	return result
}

func PosterDTOToProto(p *dto.PosterDTO) *poster.Poster {
	if p == nil {
		return nil
	}

	var metroGeo *poster.Geography
	if p.MetroGeo != nil {
		metroGeo = &poster.Geography{
			Lat: p.MetroGeo.Lat,
			Lon: p.MetroGeo.Lon,
		}
	}

	var flat *poster.Flat
	if p.Flat != nil {
		flat = &poster.Flat{
			FlatCategory: p.Flat.FlatCategory,
			Number:       int32(p.Flat.Number),
			Floor:        int32(p.Flat.Floor),
			RoomCount:    int32(p.Flat.RoomCount),
		}
	}

	var house *poster.House
	if p.House != nil {
		house = &poster.House{}
	}

	var company *poster.UtilityCompanyCard
	if p.Company != nil {
		company = &poster.UtilityCompanyCard{
			Id:          int64(p.Company.ID),
			CompanyName: p.Company.CompanyName,
			AvatarUrl:   p.Company.AvatarURL,
			Alias:       p.Company.Alias,
		}
	}

	images := make([]*poster.Photo, 0, len(p.Images))
	for _, img := range p.Images {
		images = append(images, &poster.Photo{
			ImgUrl: img.ImgURL,
			Order:  int32(img.Order),
		})
	}

	facilities := make([]*poster.Facility, 0, len(p.Facilities))
	for _, f := range p.Facilities {
		facilities = append(facilities, &poster.Facility{
			Alias: f.Alias,
			Name:  f.Name,
		})
	}

	return &poster.Poster{
		Id:    int64(p.ID),
		Alias: p.Alias,
		Price: p.Price,
		Category: &poster.Category{
			Alias: p.Category.Alias,
			Name:  p.Category.Name,
		},
		Description: p.Description,
		Area:        p.Area,
		BuildingGeo: &poster.Geography{
			Lat: p.Geo.Lat,
			Lon: p.Geo.Lon,
		},
		Address:    p.Address,
		District:   p.District,
		Metro:      p.Metro,
		MetroGeo:   metroGeo,
		City:       p.City,
		FloorCount: int32(p.FloorCount),
		Images:     images,
		Seller: &poster.PosterSeller{
			FirstName: p.Seller.SellerFirstName,
			LastName:  p.Seller.SellerLastName,
			Phone:     p.Seller.SellerPhone,
			AvatarUrl: p.Seller.SellerAvatarURL,
		},
		Flat:       flat,
		House:      house,
		Facilities: facilities,
		Company:    company,
	}
}

func MyPosterDTOToProto(p *dto.MyPosterDTO) *poster.MyPoster {
	if p == nil {
		return nil
	}

	return &poster.MyPoster{
		Id:        int64(p.ID),
		Alias:     p.Alias,
		Address:   p.Address,
		Area:      p.Area,
		Price:     p.Price,
		AvatarUrl: p.AvatarURl,
		Category: &poster.Category{
			Name:  p.Category.Name,
			Alias: p.Category.Alias,
		},
	}
}

func MyPostersDTOToProto(list []dto.MyPosterDTO) []*poster.MyPoster {
	result := make([]*poster.MyPoster, 0, len(list))

	for i := range list {
		result = append(result, MyPosterDTOToProto(&list[i]))
	}

	return result
}

func GeoJSONFeatureDTOToProto(f *dto.GeoJSONFeature) *poster.GeoJSONFeature {
	if f == nil {
		return nil
	}

	properties := make(map[string]string, len(f.Properties))
	for key, value := range f.Properties {
		properties[key] = fmt.Sprint(value)
	}

	return &poster.GeoJSONFeature{
		Type:       f.Type,
		Properties: properties,
		Geometry: &poster.Geometry{
			Type:        f.Geometry.Type,
			Coordinates: f.Geometry.Coordinates,
		},
	}
}

func GeoJSONFeaturesDTOToProto(list []dto.GeoJSONFeature) []*poster.GeoJSONFeature {
	result := make([]*poster.GeoJSONFeature, 0, len(list))

	for i := range list {
		result = append(result, GeoJSONFeatureDTOToProto(&list[i]))
	}

	return result
}

func PriceHistoryDTOToProto(p *dto.PriceHistoryDTO) *poster.PriceHistory {
	if p == nil {
		return nil
	}

	return &poster.PriceHistory{
		Date:  p.Date,
		Price: p.Price,
	}
}

func PriceHistoryDTOsToProto(list []dto.PriceHistoryDTO) []*poster.PriceHistory {
	result := make([]*poster.PriceHistory, 0, len(list))

	for i := range list {
		result = append(result, PriceHistoryDTOToProto(&list[i]))
	}

	return result
}

func FlatPosterProtoToDTO(meta *poster.FlatPosterMeta, photosMeta map[int]*poster.FlatPosterPhotoMeta, files map[int][]byte,
) (*dto.PosterInputFlatDTO, error) {
	if meta == nil {
		return nil, status.Error(codes.InvalidArgument, "poster meta is required")
	}

	req := &dto.PosterInputFlatDTO{
		UserID:         int(meta.UserId),
		Price:          meta.Price,
		Description:    meta.Description,
		CategoryAlias:  meta.CategoryAlias,
		Area:           meta.Area,
		GeoLat:         meta.GeoLat,
		GeoLon:         meta.GeoLon,
		FlatCategoryID: int(meta.FlatCategoryId),
		FlatFloor:      int(meta.FlatFloor),
		Address:        meta.Address,
		City:           meta.City,
		FloorCount:     int(meta.FloorCount),
		Features:       meta.Features,
	}

	if meta.FlatNumber != nil {
		v := int(*meta.FlatNumber)
		req.FlatNumber = &v
	}

	if meta.District != nil {
		req.District = meta.District
	}

	if meta.CompanyId != nil {
		v := int(*meta.CompanyId)
		req.CompanyID = &v
	}

	for order, photoMeta := range photosMeta {
		photo := dto.PhotoInputDTO{
			Order: order,
		}

		switch source := photoMeta.Source.(type) {
		case *poster.FlatPosterPhotoMeta_Url:
			photo.URL = &source.Url

		case *poster.FlatPosterPhotoMeta_File:
			data, ok := files[order]
			if !ok || len(data) == 0 {
				return nil, status.Errorf(codes.InvalidArgument, "photo %d file chunks are required", order)
			}

			photo.FileHeader = &dto.FileInput{
				Filename:    source.File.Filename,
				Size:        source.File.Size,
				ContentType: source.File.ContentType,
				File:        io.NopCloser(bytes.NewReader(data)),
			}

		default:
			return nil, status.Errorf(codes.InvalidArgument, "photo %d source is required", order)
		}

		req.Images = append(req.Images, photo)
	}

	sort.Slice(req.Images, func(i, j int) bool {
		return req.Images[i].Order < req.Images[j].Order
	})

	return req, nil
}
