package poster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/alias"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

const (
	MaxPostersLimit = 12
	PropertyFlat    = "flat"
	PropertyHouse   = "house"
	MetroRadius     = 2500
)

type PosterUseCase struct {
	uow    usecase.UnitOfWork
	file   usecase.FileRepo
	search usecase.SearchRepo
}

func NewPosterUseCase(uow usecase.UnitOfWork, file usecase.FileRepo, search usecase.SearchRepo) *PosterUseCase {
	return &PosterUseCase{
		uow:    uow,
		file:   file,
		search: search,
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
	if filters.Limit <= 0 || filters.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if filters.Limit > MaxPostersLimit {
		filters.Limit = MaxPostersLimit
	}
	log.Println(filters)

	response, err := uc.search.SearchPosters(ctx, filters)
	if err != nil {
		log.Printf("uc.search.SearchPosters: %s", err)
		return nil, err
	}
	ids := make([]int, 0, response.Len)
	for _, p := range response.Posters {
		ids = append(ids, p.ID)
	}
	posters, err := uc.uow.Posters().GetFlatsByIDs(ctx, ids)
	if err != nil {
		log.Printf("uc.uow.Posters().GetFlatsByIDs: %s", err)
		return nil, err
	}
	response.Posters = dto.PostersToPostersDTO(posters)
	return response, nil
}

func (uc *PosterUseCase) GetPosterByAliasUseCase(ctx context.Context, posterAlias string, userID *int) (*dto.PosterDTO, error) {
	var posterDTO *dto.PosterDTO

	poster, err := uc.uow.Posters().GetByAlias(ctx, posterAlias, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.PosterRepo.GetByAlias: %w", err)
	}

	posterDTO = dto.PosterToPosterDTO(poster)
	log.Println(poster, posterDTO)

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

	err := validator.ValidatePosterBase(poster)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePosterBase: %w", err)
	}

	err = validator.ValidatePosterInputFlat(poster)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePosterInputFlat: %w", err)
	}

	err = validator.ValidatePhotos(poster.Images)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePhotos: %w", err)
	}

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
	}

	flat := dto.PosterInputFlatDTOtoFlatInput(poster)

	post.Alias = alias.GenerateAlias(post)

	p, err := uc.uow.Posters().GetByAlias(ctx, post.Alias, nil)
	if err == nil {
		return &dto.CreatedPoster{ID: p.ID, Alias: p.Alias}, nil
	}

	dto.MakePhotoPathsForPoster(post)

	keys := make([]string, 0, len(post.Images))

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

		posterID, err := r.Posters().Create(ctx, post, propertyID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.Create: %w", err)
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

		if len(post.Images) > 0 {
			err = r.Posters().InsertMainPhoto(ctx, posterID, post.Images[0].Path)
			if err != nil {
				return fmt.Errorf("r.Posters().InsertMainPhoto: %w", err)
			}
		}

		createdPoster = &dto.CreatedPoster{
			ID:    posterID,
			Alias: post.Alias,
		}

		return nil
	})

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
	err := validator.ValidatePosterInputFlat(poster)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePosterInputFlat: %w", err)
	}

	err = validator.ValidatePhotos(poster.Images)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePhotos: %w", err)
	}
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
			err = r.Posters().InsertMainPhoto(ctx, ids.PosterID, post.Images[0].Path)
			if err != nil {
				return fmt.Errorf("r.Posters().InsertMainPhoto: %w", err)
			}
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

	for _, key := range oldKeys {
		if !slices.Contains(newKeys, key) {
			err = uc.file.Delete(ctx, key)
			if err != nil {
				return nil, fmt.Errorf("uc.file.Delete: %w", err)
			}
		}

	}

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

	return deletedPoster, nil
}
