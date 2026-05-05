package utils

import (
	"fmt"
	"io"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

func intPtrToInt32Ptr(v *int) *int32 {
	if v == nil {
		return nil
	}

	res := int32(*v)
	return &res
}

func intPtrToInt64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}

	res := int64(*v)
	return &res
}

func MapFiltersDTOToProto(params dto.PostersFiltersDTO) *poster.SearchPostersRequest {
	return &poster.SearchPostersRequest{
		Limit:            int32(params.Limit),
		Offset:           int32(params.Offset),
		SearchQuery:      params.SearchQuery,
		UtilityCompany:   params.UtilityCompany,
		Category:         params.Category,
		MaxPrice:         intPtrToInt32Ptr(params.MaxPrice),
		MinPrice:         intPtrToInt32Ptr(params.MinPrice),
		RoomCount:        intPtrToInt32Ptr(params.RoomCount),
		MaxSquare:        intPtrToInt32Ptr(params.MaxSquare),
		MinSquare:        intPtrToInt32Ptr(params.MinSquare),
		Facilities:       params.Facilities,
		MaxFlatFloor:     intPtrToInt32Ptr(params.MaxFlatFloor),
		MinFlatFloor:     intPtrToInt32Ptr(params.MinFlatFloor),
		IsNotFirstFloor:  params.IsNotFirstFloor,
		IsNotLastFloor:   params.IsNotLastFloor,
		MaxBuildingFloor: intPtrToInt32Ptr(params.MaxBuildingFloor),
		MinBuildingFloor: intPtrToInt32Ptr(params.MinBuildingFloor),
	}
}

func MapSearchPostersResponseToDTO(resp *poster.SearchPostersResponse) *dto.PostersResponse {
	if resp == nil {
		return &dto.PostersResponse{Len: 0, Posters: []dto.PosterCardDTO{}}
	}

	result := &dto.PostersResponse{
		Len: int(resp.Len),
	}
	list := make([]dto.PosterCardDTO, 0, len(resp.Posters))

	for _, p := range resp.Posters {
		card := dto.PosterCardDTO{
			ID:      int(p.Id),
			Alias:   p.Alias,
			Price:   p.Price,
			Address: p.Address,
			Area:    p.Area,
		}

		if p.ImageUrl != nil {
			card.ImgURL = p.ImageUrl
		}
		if p.Metro != nil {
			card.Metro = p.Metro
		}
		if p.FlatCategory != nil {
			card.FlatCategory = p.FlatCategory
		}

		list = append(list, card)
	}
	result.Posters = list

	return result
}

func MapProtoPosterToDTO(p *poster.Poster) *dto.PosterDTO {
	if p == nil {
		return nil
	}

	result := &dto.PosterDTO{
		ID:          int(p.Id),
		Alias:       p.Alias,
		Price:       p.Price,
		Description: p.Description,
		Area:        p.Area,
		Address:     p.Address,
		District:    p.District,
		Metro:       p.Metro,
		City:        p.City,
		FloorCount:  int(p.FloorCount),
	}

	if p.Category != nil {
		result.Category = dto.CategoryDTO{
			Alias: p.Category.Alias,
			Name:  p.Category.Name,
		}
	}

	if p.BuildingGeo != nil {
		result.Geo = dto.GeographyDTO{
			Lat: p.BuildingGeo.Lat,
			Lon: p.BuildingGeo.Lon,
		}
	}

	if p.MetroGeo != nil {
		result.MetroGeo = &dto.GeographyDTO{
			Lat: p.MetroGeo.Lat,
			Lon: p.MetroGeo.Lon,
		}
	}

	if p.Seller != nil {
		result.Seller = dto.PosterSellerDTO{
			SellerFirstName: p.Seller.FirstName,
			SellerLastName:  p.Seller.LastName,
			SellerPhone:     p.Seller.Phone,
			SellerAvatarURL: p.Seller.AvatarUrl,
		}
	}

	if p.Flat != nil {
		result.Flat = &dto.FlatDTO{
			FlatCategory: p.Flat.FlatCategory,
			Number:       int(p.Flat.Number),
			Floor:        int(p.Flat.Floor),
			RoomCount:    int(p.Flat.RoomCount),
		}
	}

	if p.House != nil {
		result.House = &dto.HouseDTO{}
	}

	if p.Company != nil {
		result.Company = &dto.UtilityCompanyCardDTO{
			ID:          int(p.Company.Id),
			CompanyName: p.Company.CompanyName,
			AvatarURL:   p.Company.AvatarUrl,
			Alias:       p.Company.Alias,
		}
	}

	result.Images = make([]dto.PhotoDTO, 0, len(p.Images))
	for _, img := range p.Images {
		if img == nil {
			continue
		}

		result.Images = append(result.Images, dto.PhotoDTO{
			ImgURL: img.ImgUrl,
			Order:  int(img.Order),
		})
	}

	result.Facilities = make([]dto.FacilityDTO, 0, len(p.Facilities))
	for _, f := range p.Facilities {
		if f == nil {
			continue
		}

		result.Facilities = append(result.Facilities, dto.FacilityDTO{
			Alias: f.Alias,
			Name:  f.Name,
		})
	}

	return result
}

func MapProtoMyPostersToDTO(items []*poster.MyPoster) []dto.MyPosterDTO {
	result := make([]dto.MyPosterDTO, 0, len(items))

	for _, p := range items {
		if p == nil {
			continue
		}

		result = append(result, dto.MyPosterDTO{
			ID:        int(p.Id),
			Alias:     p.Alias,
			Address:   p.Address,
			Area:      p.Area,
			Price:     p.Price,
			AvatarURl: p.AvatarUrl,
			Category: dto.CategoryDTO{
				Alias: p.Category.Alias,
				Name:  p.Category.Name,
			},
		})
	}

	return result
}

func MapMapBoundsDTOToProto(bounds dto.MapBounds) *poster.MapBounds {
	return &poster.MapBounds{
		Bbox: &poster.BBox{
			SouthWest: &poster.Geography{
				Lat: bounds.BBox.SouthWest.Lat,
				Lon: bounds.BBox.SouthWest.Lon,
			},
			NorthEast: &poster.Geography{
				Lat: bounds.BBox.NorthEast.Lat,
				Lon: bounds.BBox.NorthEast.Lon,
			},
		},
		Zoom: int32(bounds.Zoom),
	}
}

func MapGeoJSONResponseToDTO(resp *poster.GetPostersByCoordsResponse) *dto.GeoJSONFeatureResponse {
	if resp == nil {
		return &dto.GeoJSONFeatureResponse{
			Posters: make([]dto.GeoJSONFeature, 0),
		}
	}

	result := &dto.GeoJSONFeatureResponse{
		Len:     int(resp.Len),
		Posters: make([]dto.GeoJSONFeature, 0, len(resp.Features)),
	}

	for _, f := range resp.Features {
		if f == nil {
			continue
		}

		properties := make(map[string]any, len(f.Properties))
		for key, value := range f.Properties {
			if value == nil {
				continue
			}

			properties[key] = value.AsInterface()
		}

		geometry := dto.Geometry{}
		if f.Geometry != nil {
			geometry = dto.Geometry{
				Type:        f.Geometry.Type,
				Coordinates: f.Geometry.Coordinates,
			}
		}

		result.Posters = append(result.Posters, dto.GeoJSONFeature{
			Type:       f.Type,
			Properties: properties,
			Geometry:   geometry,
		})
	}

	return result
}

func MapProtoPriceHistoryToDTO(items []*poster.PriceHistory) []dto.PriceHistoryDTO {
	result := make([]dto.PriceHistoryDTO, 0, len(items))

	for _, item := range items {
		if item == nil {
			continue
		}

		result = append(result, dto.PriceHistoryDTO{
			Date:  item.Date,
			Price: item.Price,
		})
	}

	return result
}

func SendFlatPosterMeta(
	stream poster.PosterService_CreateFlatPosterClient,
	req *dto.PosterInputFlatDTO,
) error {
	return stream.Send(&poster.CreateFlatPosterRequest{
		Payload: &poster.CreateFlatPosterRequest_PosterMeta{
			PosterMeta: &poster.FlatPosterMeta{
				UserId:         int64(req.UserID),
				Price:          req.Price,
				Description:    req.Description,
				CategoryAlias:  req.CategoryAlias,
				Area:           req.Area,
				GeoLat:         req.GeoLat,
				GeoLon:         req.GeoLon,
				FlatCategoryId: int64(req.FlatCategoryID),
				FlatNumber:     intPtrToInt32Ptr(req.FlatNumber),
				FlatFloor:      int32(req.FlatFloor),
				Address:        req.Address,
				City:           req.City,
				District:       req.District,
				FloorCount:     int32(req.FloorCount),
				CompanyId:      intPtrToInt64Ptr(req.CompanyID),
				Features:       req.Features,
			},
		},
	})
}

func SendFlatPosterPhotos(stream poster.PosterService_CreateFlatPosterClient, photos []dto.PhotoInputDTO) error {
	const chunkSize = 32 * 1024

	buf := make([]byte, chunkSize)

	for _, photo := range photos {
		if photo.URL != nil {
			err := stream.Send(&poster.CreateFlatPosterRequest{
				Payload: &poster.CreateFlatPosterRequest_PhotoMeta{
					PhotoMeta: &poster.FlatPosterPhotoMeta{
						Order: int32(photo.Order),
						Source: &poster.FlatPosterPhotoMeta_Url{
							Url: *photo.URL,
						},
					},
				},
			})
			if err != nil {
				return err
			}

			continue
		}

		if photo.FileHeader == nil {
			return fmt.Errorf("photo %d: file or url is required", photo.Order)
		}

		err := stream.Send(&poster.CreateFlatPosterRequest{
			Payload: &poster.CreateFlatPosterRequest_PhotoMeta{
				PhotoMeta: &poster.FlatPosterPhotoMeta{
					Order: int32(photo.Order),
					Source: &poster.FlatPosterPhotoMeta_File{
						File: &poster.FileInput{
							Filename:    photo.FileHeader.Filename,
							Size:        photo.FileHeader.Size,
							ContentType: photo.FileHeader.ContentType,
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		for {
			n, readErr := photo.FileHeader.File.Read(buf)

			if n > 0 {
				err = stream.Send(&poster.CreateFlatPosterRequest{
					Payload: &poster.CreateFlatPosterRequest_PhotoChunk{
						PhotoChunk: &poster.PhotoChunk{
							Order: int32(photo.Order),
							Data:  buf[:n],
						},
					},
				})
				if err != nil {
					return err
				}
			}

			if readErr == io.EOF {
				break
			}

			if readErr != nil {
				return readErr
			}
		}
	}

	return nil
}

func SendUpdateFlatPosterMeta(stream poster.PosterService_UpdateFlatPosterClient, alias string, req *dto.PosterInputFlatDTO,
) error {
	return stream.Send(&poster.UpdateFlatPosterRequest{
		Payload: &poster.UpdateFlatPosterRequest_PosterMeta{
			PosterMeta: &poster.UpdateFlatPosterMeta{
				Alias: alias,
				Poster: &poster.FlatPosterMeta{
					UserId:         int64(req.UserID),
					Price:          req.Price,
					Description:    req.Description,
					CategoryAlias:  req.CategoryAlias,
					Area:           req.Area,
					GeoLat:         req.GeoLat,
					GeoLon:         req.GeoLon,
					FlatCategoryId: int64(req.FlatCategoryID),
					FlatNumber:     intPtrToInt32Ptr(req.FlatNumber),
					FlatFloor:      int32(req.FlatFloor),
					Address:        req.Address,
					City:           req.City,
					District:       req.District,
					FloorCount:     int32(req.FloorCount),
					CompanyId:      intPtrToInt64Ptr(req.CompanyID),
					Features:       req.Features,
				},
			},
		},
	})
}

func SendUpdateFlatPosterPhotos(stream poster.PosterService_UpdateFlatPosterClient, photos []dto.PhotoInputDTO) error {
	const chunkSize = 32 * 1024

	buf := make([]byte, chunkSize)

	for _, photo := range photos {
		if photo.URL != nil {
			err := stream.Send(&poster.UpdateFlatPosterRequest{
				Payload: &poster.UpdateFlatPosterRequest_PhotoMeta{
					PhotoMeta: &poster.FlatPosterPhotoMeta{
						Order: int32(photo.Order),
						Source: &poster.FlatPosterPhotoMeta_Url{
							Url: *photo.URL,
						},
					},
				},
			})
			if err != nil {
				return err
			}
			continue
		}

		if photo.FileHeader == nil {
			return fmt.Errorf("photo %d: file or url is required", photo.Order)
		}

		err := stream.Send(&poster.UpdateFlatPosterRequest{
			Payload: &poster.UpdateFlatPosterRequest_PhotoMeta{
				PhotoMeta: &poster.FlatPosterPhotoMeta{
					Order: int32(photo.Order),
					Source: &poster.FlatPosterPhotoMeta_File{
						File: &poster.FileInput{
							Filename:    photo.FileHeader.Filename,
							Size:        photo.FileHeader.Size,
							ContentType: photo.FileHeader.ContentType,
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		for {
			n, readErr := photo.FileHeader.File.Read(buf)

			if n > 0 {
				err = stream.Send(&poster.UpdateFlatPosterRequest{
					Payload: &poster.UpdateFlatPosterRequest_PhotoChunk{
						PhotoChunk: &poster.PhotoChunk{
							Order: int32(photo.Order),
							Data:  buf[:n],
						},
					},
				})
				if err != nil {
					return err
				}
			}

			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				return readErr
			}
		}
	}

	return nil
}
