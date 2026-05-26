package poster

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/alias"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/cache"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

const (
	MaxPostersLimit = 12
	PropertyFlat    = "flat"
	PropertyHouse   = "house"
	MetroRadius     = 2500
	DefaultCacheTTL = 5 * time.Minute
)

type PosterUseCase struct {
	uow    usecase.UnitOfWork
	file   usecase.FileRepo
	search usecase.SearchRepo
	agent  usecase.LLMAgent
	cache  usecase.Cache
	maps   usecase.StreetMapProvider
}

func NewPosterUseCase(uow usecase.UnitOfWork, file usecase.FileRepo, search usecase.SearchRepo, agent usecase.LLMAgent, cache usecase.Cache, maps usecase.StreetMapProvider) *PosterUseCase {
	return &PosterUseCase{
		uow:    uow,
		file:   file,
		search: search,
		agent:  agent,
		cache:  cache,
		maps:   maps,
	}
}

func (uc *PosterUseCase) GetPostersUseCase(ctx context.Context, filters dto.PostersFiltersDTO) (*dto.PostersResponse, error) {
	if filters.Limit <= 0 || filters.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if filters.Limit > MaxPostersLimit {
		filters.Limit = MaxPostersLimit
	}

	posters, err := uc.uow.Posters().GetFlatsAll(ctx, filters)
	if err != nil {
		log.Printf("uc.repo.GetFlatsAll: %s", err)
		return nil, err
	}
	len, err := uc.uow.Posters().CountPosters(ctx)
	if err != nil {
		log.Printf("uc.repo.CountPosters: %s", err)
		return nil, err
	}

	response := dto.PostersResponse{
		Posters: dto.PostersToPostersDTO(posters),
		Len:     len,
	}
	return &response, nil
}

func (uc *PosterUseCase) SearchPostersUseCase(ctx context.Context, filters dto.PostersFiltersDTO) (*dto.PostersResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "SearchPostersUseCase")
	if filters.Limit <= 0 || filters.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if filters.Limit > MaxPostersLimit {
		filters.Limit = MaxPostersLimit
	}
	hash, err := cache.GenerateCacheKey(filters)
	if err != nil {
		return nil, fmt.Errorf("cache.GenerateCacheKey: %w", err)
	}

	cacheKey := fmt.Sprintf("posters:filters:%s", hash)
	data, err := uc.cache.Get(ctx, cacheKey)
	if err == nil {
		var cachedPoster dto.PostersResponse
		if err := json.Unmarshal(data, &cachedPoster); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(data, posterDTO): %w", err)
		}
		log.Info("cache hit")
		return &cachedPoster, nil
	}
	log.Info("cache miss")
	log.Println(filters)

	response, err := uc.search.SearchPosters(ctx, filters)
	if err != nil {
		log.Printf("uc.search.SearchPosters: %s", err)
		return nil, err
	}
	data, err = json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(posterDTO): %w", err)
	}
	err = uc.cache.Set(ctx, cacheKey, data, DefaultCacheTTL)
	if err != nil {
		return nil, fmt.Errorf("uc.cache.Set: %w", err)
	}
	// ids := make([]int, 0, response.Len)
	// for _, p := range response.Posters {
	// 	ids = append(ids, p.ID)
	// }
	// posters, err := uc.uow.Posters().GetFlatsByIDs(ctx, ids)
	// if err != nil {
	// 	log.Printf("uc.uow.Posters().GetFlatsByIDs: %s", err)
	// 	return nil, err
	// }
	// response.Posters = dto.PostersToPostersDTO(posters)
	return response, nil
}

func (uc *PosterUseCase) GetPosterByAliasUseCase(ctx context.Context, posterAlias string, userID *int) (*dto.PosterDTO, error) {
	var posterDTO *dto.PosterDTO
	log := ctxLogger.GetLogger(ctx).WithField("op", "GetPosterByAliasUseCase")
	cacheKey := fmt.Sprintf("posters:%s", posterAlias)
	data, err := uc.cache.Get(ctx, cacheKey)
	if err == nil {
		var cachedPoster dto.PosterDTO
		if err := json.Unmarshal(data, &cachedPoster); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(data, posterDTO): %w", err)
		}
		log.Info("cache hit")
		return &cachedPoster, nil
	}
	log.Info("cache miss")

	poster, err := uc.uow.Posters().GetByAlias(ctx, posterAlias, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.PosterRepo.GetByAlias: %w", err)
	}
	posterDTO = dto.PosterToPosterDTO(poster)

	switch poster.CategoryAlias {
	case PropertyFlat:
		flat, err := uc.uow.Posters().GetFlatByPropetyID(ctx, poster.PropertyID)
		if err != nil {
			log.Printf("uc.uow.Posters().GetFlatsAll: %s", err)
			return nil, err
		}

		flatDTO := dto.FlatToFlatFlatDTO(flat)
		posterDTO.Flat = flatDTO

	case PropertyHouse:
	default:
		return nil, fmt.Errorf("no such category %w", entity.ServiceError)
	}

	dto.MakeUrlsFromPaths(posterDTO, config.Config.PublicHost, config.Config.Bucket)

	data, err = json.Marshal(posterDTO)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(posterDTO): %w", err)
	}
	err = uc.cache.Set(ctx, cacheKey, data, DefaultCacheTTL)
	if err != nil {
		return nil, fmt.Errorf("uc.cache.Set: %w", err)
	}

	return posterDTO, nil
}

func (uc *PosterUseCase) GetPosterByUserID(ctx context.Context, userID int) ([]dto.MyPosterDTO, error) {
	posters, err := uc.uow.Posters().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.MyPosterToMyPosterDTO(posters), nil
}

func (uc *PosterUseCase) GetMetroStationsByRadius(ctx context.Context, coords dto.GeographyDTO) ([]dto.MetroStationDTO, error) {
	stations, err := uc.uow.Posters().GetMetroStationByRadius(ctx, coords, MetroRadius)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetMetroStationByRadius: %s", err)
	}
	dtos := make([]dto.MetroStationDTO, 0, len(stations))
	for _, d := range stations {
		dtos = append(dtos, dto.MetroToMetroStationDTO(d))
	}
	return dtos, nil

}

func (uc *PosterUseCase) CreateFlatPoster(ctx context.Context, poster *dto.PosterInputFlatDTO) (*dto.CreatedPoster, error) {
	var createdPoster *dto.CreatedPoster

	// err := validator.ValidatePosterBase(poster)
	// if err != nil {
	// 	return nil, fmt.Errorf("validator.ValidatePosterBase: %w", err)
	// }

	// err = validator.ValidatePhotos(poster.Images)
	// if err != nil {
	// 	return nil, fmt.Errorf("validator.ValidatePhotos: %w", err)
	// }
	validator.SanitizePosterInput(poster)

	post := dto.PosterInputFlatDTOtoPosterInput(poster)

	city, err := uc.uow.Posters().GetCityByName(ctx, poster.City)
	if err != nil {
		if errors.Is(err, entity.NotFoundError) {
			city, err = uc.uow.Posters().CreateCity(ctx, poster.City)
			if err != nil {
				return nil, fmt.Errorf("uc.PosterRepo.CreateCity: %w", err)
			}
		} else {
			return nil, fmt.Errorf("uc.PosterRepo.GetCityByName: %w", err)
		}
	}

	log.Printf("city: %s", city.Name)

	post.CityID = city.ID

	stations, err := uc.uow.Posters().GetMetroStationByRadius(ctx, dto.GeographyDTO(post.Geo), MetroRadius)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetMetroStationByRadius: %s", err)
	}

	if len(stations) > 0 {
		station := stations[0].ID
		post.MetroStationID = &station
	} else {
		stations, err = uc.maps.GetMetroStationByRadius(ctx, dto.GeographyDTO(post.Geo), MetroRadius)
		if err != nil {
			return nil, fmt.Errorf("uc.maps.GetMetroStationByRadius: %s", err)
		}
		if len(stations) > 0 {
			newSt, err := uc.uow.Posters().CreateMetroStation(ctx, stations[0].StationName, dto.GeographyDTO(stations[0].StationGEO))
			if err != nil {
				return nil, fmt.Errorf("uc.uow.Posters().CreateMetroStation: %s", err)
			}
			if newSt != nil {
				station := newSt.ID
				post.MetroStationID = &station
			}
		}
	}

	flat := dto.PosterInputFlatDTOtoFlatInput(poster)

	post.Alias = alias.GenerateAlias(post)

	p, err := uc.uow.Posters().GetByAlias(ctx, post.Alias, nil)
	if err == nil {
		return &dto.CreatedPoster{ID: p.ID, Alias: p.Alias}, nil
	}

	dto.MakePhotoPathsForPoster(post)

	keys := make([]string, 0, len(post.Images))
	var posterID int

	err = uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {

		buildingID, err := r.Posters().CreateBuilding(ctx, post)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.CreateBuilding: %w", err)
		}

		propertyID, err := r.Posters().CreateProperty(ctx, post, buildingID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.CreateProperty: %w", err)
		}

		err = r.Posters().InsertFacilities(ctx, propertyID, post.Features)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.InsertFacilities: %w", err)
		}

		posterID, err = r.Posters().Create(ctx, post, propertyID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.Create: %w", err)
		}

		err = r.Posters().AddPriceHistory(ctx, posterID, post.Price)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.AddPriceHistory: %w", err)
		}

		flat.PropertyID = propertyID
		err = r.Posters().InsertFlat(ctx, flat)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.InsertFlat: %w", err)
		}

		for _, photoPoster := range post.Images {
			key, err := uc.uploadPhoto(ctx, photoPoster)
			if err != nil {
				return fmt.Errorf("uc.uploadPhoto: %w", err)
			}
			keys = append(keys, key)
		}

		err = r.Posters().InsertPhotos(ctx, posterID, post.Images)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.InsertPhotos: %w", err)
		}

		// if len(post.Images) > 0 {

		// 	err = r.Posters().InsertMainPhoto(ctx, posterID, post.Images[0].Path)
		// 	if err != nil {
		// 		return fmt.Errorf("r.Posters().InsertMainPhoto: %w", err)
		// 	}
		// }

		createdPoster = &dto.CreatedPoster{
			ID:    posterID,
			Alias: post.Alias,
		}

		return nil
	})
	if len(keys) > 0 && len(post.Images) > 0 {
		go uc.scaleAndCropPhoto(context.Background(), keys[0], post.Images[0].Path, posterID)
	}

	if err != nil {
		cleanErr := uc.cleanUploadedFiles(ctx, keys)
		if cleanErr != nil {
			return nil, fmt.Errorf("uc.cleanUploadedFiles: %w", cleanErr)
		}

		return nil, fmt.Errorf("uc.uow.Do: %w", err)
	}

	return createdPoster, nil
}

func (uc *PosterUseCase) cleanUploadedFiles(ctx context.Context, keys []string) error {
	var resultErr error
	for _, key := range keys {
		err := uc.file.Delete(ctx, key)
		if err != nil {
			resultErr = fmt.Errorf("uc.file.Delete: %w", err)
		}
	}

	return resultErr
}

func (uc *PosterUseCase) uploadPhoto(ctx context.Context, photoPoster dto.PhotoInput) (string, error) {
	file := photoPoster.FileHeader.File
	defer file.Close()

	if !validator.ValidatePhoto(photoPoster.FileHeader) {
		return "", entity.NewValidationError("photo")
	}

	key := photo.GetKeyFromPath(photoPoster.Path)
	size := photoPoster.FileHeader.Size
	contentType := photoPoster.FileHeader.ContentType

	log.Println("key ", key)

	if err := uc.file.Upload(ctx, key, file, size, contentType); err != nil {
		return "", fmt.Errorf("uc.file.Upload: %w", err)
	}

	return key, nil
}

func (uc *PosterUseCase) UpdateFlatPoster(ctx context.Context, alias string, poster *dto.PosterInputFlatDTO) (*dto.CreatedPoster, error) {
	err := validator.ValidatePosterBase(poster)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePosterInputFlat: %w", err)
	}

	err = validator.ValidatePhotos(poster.Images)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePhotos: %w", err)
	}
	validator.SanitizePosterInput(poster)
	ids, err := uc.uow.Posters().GetUpdateIDsByAlias(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("uc.PosterRepo.GetUpdateIDsByPosterID: %w", err)
	}

	if ids.UserID != poster.UserID {
		return nil, entity.NotFoundError
	}

	post := dto.PosterInputFlatDTOtoPosterInput(poster)
	flat := dto.PosterInputFlatDTOtoFlatInput(poster)

	post.Alias = alias

	dto.MakePhotoPathsForPoster(post)

	log.Print(post.Images)

	city, err := uc.uow.Posters().GetCityByName(ctx, poster.City)
	if err != nil {
		if errors.Is(err, entity.NotFoundError) {
			city, err = uc.uow.Posters().CreateCity(ctx, poster.City)
			if err != nil {
				return nil, fmt.Errorf("uc.PosterRepo.CreateCity: %w", err)
			}
		} else {
			return nil, fmt.Errorf("uc.PosterRepo.GetCityByName: %w", err)
		}
	}
	post.CityID = city.ID

	stations, err := uc.uow.Posters().GetMetroStationByRadius(ctx, dto.GeographyDTO(post.Geo), MetroRadius)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetMetroStationByRadius: %s", err)
	}
	if len(stations) > 0 {
		station := stations[0].ID
		post.MetroStationID = &station
	} else {
		stations, err = uc.maps.GetMetroStationByRadius(ctx, dto.GeographyDTO(post.Geo), MetroRadius)
		if err != nil {
			return nil, fmt.Errorf("uc.maps.GetMetroStationByRadius: %s", err)
		}
		if len(stations) > 0 {
			newSt, err := uc.uow.Posters().CreateMetroStation(ctx, stations[0].StationName, dto.GeographyDTO(stations[0].StationGEO))
			if err != nil {
				return nil, fmt.Errorf("uc.uow.Posters().CreateMetroStation: %s", err)
			}
			if newSt != nil {
				station := newSt.ID
				post.MetroStationID = &station
			}
		}
	}

	flat.PropertyID = ids.PropertyID

	var oldKeys []string
	var newKeys []string

	if len(post.Images) > 0 {
		oldPaths, err := uc.uow.Posters().GetPhotoPathsByPosterID(ctx, ids.PosterID)
		if err != nil {
			return nil, fmt.Errorf("uc.PosterRepo.GetPhotoPathsByPosterID: %w", err)
		}

		oldKeys = make([]string, 0, len(oldPaths))
		for _, path := range oldPaths {
			oldKeys = append(oldKeys, photo.GetKeyFromPath(path))
		}
		log.Printf("oldKeys: %v", oldKeys)

		newKeys = make([]string, 0, len(post.Images))
		for _, photoPoster := range post.Images {
			var key string
			if photoPoster.FileHeader == nil {
				key = photo.GetKeyFromPath(photoPoster.Path)
			} else {
				key, err = uc.uploadPhoto(ctx, photoPoster)
				if err != nil {
					_ = uc.cleanUploadedFiles(ctx, newKeys) // FIXME: а если старый путь и новый совпадают?
					return nil, fmt.Errorf("uc.uploadPhoto: %w", err)
				}
			}

			newKeys = append(newKeys, key)
		}
	}

	var createdPoster *dto.CreatedPoster

	err = uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {
		lastHistory, err := r.Posters().GetLastPriceHistoryByPosterID(ctx, ids.PosterID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.GetLastPriceHistoryByPosterID: %w", err)
		}

		if lastHistory.Price != post.Price {
			err = r.Posters().AddPriceHistory(ctx, ids.PosterID, post.Price)
			if err != nil {
				return fmt.Errorf("uc.PosterRepo.AddPriceHistory: %w", err)
			}
		}

		err = r.Posters().Update(ctx, ids.PosterID, post)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.Update: %w", err)
		}

		err = r.Posters().UpdateProperty(ctx, ids.PropertyID, post)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.UpdateProperty: %w", err)
		}

		err = r.Posters().UpdateBuilding(ctx, ids.BuildingID, post)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.UpdateBuilding: %w", err)
		}

		err = r.Posters().DeleteFacilitiesByPropertyID(ctx, ids.PropertyID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.DeleteFacilitiesByPropertyID: %w", err)
		}

		if len(post.Features) > 0 {
			err = r.Posters().InsertFacilities(ctx, ids.PropertyID, post.Features)
			if err != nil {
				return fmt.Errorf("uc.PosterRepo.InsertFacilities: %w", err)
			}
		}

		err = r.Posters().UpdateFlat(ctx, flat)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.UpdateFlat: %w", err)
		}

		if len(post.Images) > 0 {
			err = r.Posters().DeletePhotosByPosterID(ctx, ids.PosterID)
			if err != nil {
				return fmt.Errorf("uc.PosterRepo.DeletePhotosByPosterID: %w", err)
			}

			err = r.Posters().InsertPhotos(ctx, ids.PosterID, post.Images)
			if err != nil {
				return fmt.Errorf("uc.PosterRepo.InsertPhotos: %w", err)
			}
			// err = r.Posters().InsertMainPhoto(ctx, ids.PosterID, post.Images[0].Path)
			// if err != nil {
			// 	return fmt.Errorf("r.Posters().InsertMainPhoto: %w", err)
			// }
		}

		createdPoster = &dto.CreatedPoster{
			ID:    ids.PosterID,
			Alias: alias,
		}

		return nil
	})
	if err != nil {
		if len(newKeys) > 0 {
			_ = uc.cleanUploadedFiles(ctx, newKeys) // FIXME: а если старый путь и новый совпадают?
		}
		return nil, fmt.Errorf("uc.uow.Do: %w", err)
	}
	if len(newKeys) > 0 && len(post.Images) > 0 {
		go uc.scaleAndCropPhoto(context.Background(), newKeys[0], post.Images[0].Path, ids.PosterID)
	}

	for _, key := range oldKeys {
		if !slices.Contains(newKeys, key) {
			err = uc.file.Delete(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("uc.file.Delete: %w", err)
			}
		}

	}
	cacheKey := fmt.Sprintf("posters:%s", alias)
	uc.cache.Delete(ctx, cacheKey)

	return createdPoster, nil
}

func (uc *PosterUseCase) DeleteFlatPoster(ctx context.Context, alias string, userID int) (*dto.CreatedPoster, error) {
	ids, err := uc.uow.Posters().GetUpdateIDsByAlias(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("get ids by alias: %w", err)
	}
	if ids.UserID != userID {
		return nil, entity.NotFoundError
	}

	oldPaths, err := uc.uow.Posters().GetPhotoPathsByPosterID(ctx, ids.PosterID)
	if err != nil {
		return nil, fmt.Errorf("get photo paths: %w", err)
	}

	oldKeys := make([]string, 0, len(oldPaths))
	for _, path := range oldPaths {
		oldKeys = append(oldKeys, photo.GetKeyFromPath(path))
	}

	var deletedPoster *dto.CreatedPoster

	err = uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {
		if len(oldPaths) > 0 {
			err := r.Posters().DeletePhotosByPosterID(ctx, ids.PosterID)
			if err != nil {
				return fmt.Errorf("delete photos: %w", err)
			}
		}
		err = r.Posters().DeleteFacilitiesByPropertyID(ctx, ids.PropertyID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.DeleteFacilitiesByPropertyID: %w", err)
		}

		err := r.Posters().Delete(ctx, ids.PosterID)
		if err != nil {
			return fmt.Errorf("delete poster: %w", err)
		}
		err = r.Posters().DeleteFlat(ctx, ids.PropertyID)
		if err != nil {
			return fmt.Errorf("delete property: %w", err)
		}
		err = r.Posters().DeleteProperty(ctx, ids.PropertyID)
		if err != nil {
			return fmt.Errorf("delete property: %w", err)
		}
		err = r.Posters().DeleteBuilding(ctx, ids.BuildingID)
		if err != nil {
			return fmt.Errorf("delete building: %w", err)
		}
		err = uc.search.DeletePoster(ctx, ids.PosterID)
		if err != nil {
			return fmt.Errorf("uc.search.DeletePoster: %w", err)
		}

		deletedPoster = &dto.CreatedPoster{
			ID:    ids.PosterID,
			Alias: alias,
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("uow transaction: %w", err)
	}

	if len(oldKeys) > 0 {
		_ = uc.cleanUploadedFiles(ctx, oldKeys)
	}
	cacheKey := fmt.Sprintf("posters:%s", alias)
	uc.cache.Delete(ctx, cacheKey)

	return deletedPoster, nil
}

func (uc *PosterUseCase) GetPostersByCoords(ctx context.Context, bounds dto.MapBounds, filters dto.PostersFiltersDTO) (*dto.GeoJSONFeatureResponse, error) {
	filters.Offset = 0
	filters.Limit = 100

	if bounds.Zoom < 13 {
		clusters, err := uc.search.GetClustersByMapBounds(ctx, bounds, filters)
		if err != nil {
			return nil, fmt.Errorf("uc.uow.Posters().GetClustersByCoords: %w", err)
		}

		return &dto.GeoJSONFeatureResponse{
			Posters: dto.ClustersToGEOJsons(clusters),
			Len:     len(clusters),
		}, nil
	}

	posters, err := uc.search.GetPostersByMapBounds(ctx, bounds, filters)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetPostersByCoords: %w", err)
	}

	return &dto.GeoJSONFeatureResponse{
		Posters: dto.PostersToGEOJsons(posters),
		Len:     len(posters),
	}, nil
}

func (uc *PosterUseCase) GetPostersByRadius(ctx context.Context, point dto.GeographyDTO) ([]dto.MyPosterDTO, error) {
	posters, err := uc.uow.Posters().GetPostersByRadius(ctx, point, 10)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetPostersByRadius: %w", err)
	}

	return dto.MyPosterToMyPosterDTO(posters), nil
}

func (uc *PosterUseCase) AddViewPoster(ctx context.Context, alias string, userID int) error {
	poster, err := uc.uow.Posters().GetByAlias(ctx, alias, nil)
	if err != nil {
		return fmt.Errorf("uc.PosterRepo.GetByAlias: %w", err)
	}

	uc.uow.Posters().AddView(ctx, userID, poster.ID)

	return nil
}

func (uc *PosterUseCase) AddFavoritePoster(ctx context.Context, alias string, userID int) error {
	poster, err := uc.uow.Posters().GetByAlias(ctx, alias, nil)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().GetByAlias: %w", err)
	}

	err = uc.uow.Posters().AddFavorite(ctx, userID, poster.ID)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().AddFavorite: %w", err)
	}

	return nil
}

func (uc *PosterUseCase) GetFavoritesPoster(ctx context.Context, userID int) (*dto.PostersResponse, error) {
	posters, err := uc.uow.Posters().GetFavoritesFlatsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetFavoritesFlatsByUserID: %w", err)
	}

	len, err := uc.uow.Posters().CountFavoritesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().СountFavoritesByUserID: %w", err)
	}

	response := dto.PostersResponse{
		Posters: dto.PostersToPostersDTO(posters),
		Len:     len,
	}

	return &response, nil
}

func (uc *PosterUseCase) DeleteFavoritePoster(ctx context.Context, alias string, userID int) error {
	poster, err := uc.uow.Posters().GetByAlias(ctx, alias, nil)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().GetByAlias: %w", err)
	}

	err = uc.uow.Posters().DeleteFavorite(ctx, userID, poster.ID)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().DeleteFavorite: %w", err)
	}

	return nil
}

func (uc *PosterUseCase) GetFavoritesCountPoster(ctx context.Context, posterAlias string, userID *int) (int, bool, error) {
	var isFavorite bool
	count, err := uc.uow.Posters().GetFavoritesCountByAlias(ctx, posterAlias)
	if err != nil {
		return 0, false, fmt.Errorf("uc.PosterRepo.GetFavoritesCountByAlias: %w", err)
	}
	if userID != nil {
		isFavorite, err = uc.uow.Posters().IsFavoriteByAliasAndUserID(ctx, posterAlias, *userID)
		if err != nil {
			return 0, false, fmt.Errorf("uc.uow.Posters().IsFavorite: %w", err)
		}
	}

	return count, isFavorite, nil
}

func (uc *PosterUseCase) GetViewsPoster(ctx context.Context, alias string) (int, error) {
	poster, err := uc.uow.Posters().GetByAlias(ctx, alias, nil)
	if err != nil {
		return 0, fmt.Errorf("uc.PosterRepo.GetByAlias: %w", err)
	}

	views, err := uc.uow.Posters().GetViewsCount(ctx, poster.ID)
	if err != nil {
		return 0, fmt.Errorf("uc.PosterRepo.GetViewsCount: %w", err)
	}

	return views, nil
}

func (uc *PosterUseCase) GenerateDescription(ctx context.Context, input dto.GenerateDescriptionDTO) (string, error) {
	const systemPrompt = `
	Ты — система генерации описаний недвижимости.

	ОБЯЗАТЕЛЬНЫЕ ПРАВИЛА:
	1. Пиши ТОЛЬКО на русском языке.
	2. Запрещено использовать любые другие языки (включая английский, китайский и др.).
	3. НЕ добавляй информацию, которой нет во входных данных.
	4. НЕ используй форматирование (Markdown, списки, заголовки).
	5. Текст должен быть одним абзацем.
	6. Максимум 2000 символов.
	7. Стиль: профессиональный, продающий, но без фантазий.

	Если данных недостаточно — просто опиши только то, что есть, без догадок.
	`

	const prompt = `
	ДАННЫЕ ОБЪЕКТА:
	Тип: %s
	Комнаты: %s
	Площадь: %.2f кв.м
	Особенности: %s

	Сгенерируй описание объявления.
	`

	res, err := uc.agent.Chat(ctx, systemPrompt, fmt.Sprintf(prompt, input.Category, input.FlatCategory, input.Area, input.Features))
	if err != nil {
		return "", fmt.Errorf("uc.agent.Chat: %w", err)
	}
	return res.Content, nil
}

func (uc *PosterUseCase) GetPriceHistoryPoster(ctx context.Context, posterAlias string) ([]dto.PriceHistoryDTO, error) {
	history, err := uc.uow.Posters().GetPriceHistoryByAlias(ctx, posterAlias)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetPriceHistoryByAlias: %w", err)
	}

	return dto.PriceHistoryToPriceHistoryDTO(history), nil
}

func (uc *PosterUseCase) GetPosterRoommates(ctx context.Context, alias string) ([]dto.RoommateUserDTO, error) {
	users, err := uc.uow.Posters().GetPosterRoommates(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Posters().GetPosterRoommates: %w", err)
	}
	for i, user := range users {
		if user.AvatarURL != nil {
			avatar := photo.MakeUrlFromPath(*user.AvatarURL, config.Config.PublicHost, config.Config.Bucket)
			users[i].AvatarURL = &avatar
		}
	}

	return users, nil
}

func (uc *PosterUseCase) AddPosterRoommate(ctx context.Context, alias string, userID int) error {
	hasForm, err := uc.uow.Posters().HasRoommateForm(ctx, userID)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().HasRoommateForm: %w", err)
	}

	if !hasForm {
		return entity.InvalidInput
	}

	err = uc.uow.Posters().AddPosterRoommate(ctx, alias, userID)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().AddPosterRoommate: %w", err)
	}

	return nil
}

func (uc *PosterUseCase) DeletePosterRoommate(ctx context.Context, alias string, userID int) error {
	err := uc.uow.Posters().DeletePosterRoommate(ctx, alias, userID)
	if err != nil {
		return fmt.Errorf("uc.uow.Posters().DeletePosterRoommate: %w", err)
	}

	return nil
}

func (uc *PosterUseCase) scaleAndCropPhoto(ctx context.Context, srcKey string, srcPath string, posterID int) error {
	file, err := uc.file.Get(context.Background(), srcKey)
	if err != nil {
		log.Printf("Failed to get original from MinIO: %v", err)
		return err
	}
	defer file.Close()
	buffer, err := io.ReadAll(file)
	if err != nil {
		log.Printf("io.ReadAll: %s", err)
		return err
	}
	log.Printf("Buffer size: %d", len(buffer))
	if len(buffer) > 20 {
		log.Printf("First bytes (hex): %x", buffer[:min(20, len(buffer))])
	}
	mime := http.DetectContentType(buffer)
	log.Printf("Detected MIME: %s", mime)
	if !strings.HasPrefix(mime, "image/") {
		return fmt.Errorf("file is not an image, MIME: %s", mime)
	}
	newImage, err := photo.ResizeAndCropJPEG(buffer, 600, 340, 80)
	if err != nil {
		log.Printf("resizeAndCropJPEG: %s", err)
		return fmt.Errorf("resizeAndCropJPEG: %w", err)
	}
	size := int64(len(newImage))
	contentType := "image/jpeg"
	path := strings.Replace(srcPath, "poster/img/", "poster/img/scaled/", 1)
	key := photo.GetKeyFromPath(path)

	if err := uc.file.Upload(context.Background(), key, bytes.NewReader(newImage), size, contentType); err != nil {
		log.Printf("uc.file.Upload: %s", err)
		return fmt.Errorf("uc.file.Upload: %w", err)
	}

	err = uc.uow.Posters().InsertMainPhoto(context.Background(), posterID, path)
	if err != nil {
		log.Printf("uc.uow.Posters().InsertMainPhoto: %s", err)
		return fmt.Errorf("r.Posters().InsertMainPhoto: %w", err)
	}
	return nil
}
