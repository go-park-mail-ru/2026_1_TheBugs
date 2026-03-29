package poster

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/alias"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

const (
	MaxPostersLimit = 12
	PropertyFlat    = "Квартиры" // TODO: делать по id а не по имени а то оно меняется либо поменять в
	PropertyHouse   = "house"
)

type PosterUseCase struct {
	uow  usecase.UnitOfWork
	file usecase.FileRepo
}

func NewPosterUseCase(uow usecase.UnitOfWork, file usecase.FileRepo) *PosterUseCase {
	return &PosterUseCase{
		uow:  uow,
		file: file,
	}
}

func (uc *PosterUseCase) GetPostersUseCase(ctx context.Context, filters dto.PostersFiltersDTO) ([]dto.PosterCardDTO, error) {
	if filters.Limit <= 0 || filters.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if filters.Limit > MaxPostersLimit {
		filters.Limit = MaxPostersLimit
	}

	posters, err := uc.uow.Posters().GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("uc.PosterRepo.GetAll: %w", err)
	}

	return dto.PostersToPostersDTO(posters), nil
}

func (uc *PosterUseCase) GetPosterByAliasUseCase(ctx context.Context, posterAlias string) (*dto.PosterDTO, error) {
	var posterDTO *dto.PosterDTO

	poster, err := uc.uow.Posters().GetByAlias(ctx, posterAlias)
	if err != nil {
		return nil, fmt.Errorf("uc.PosterRepo.GetByAlias: %w", err)
	}

	posterDTO = dto.PosterToPosterDTO(poster)
	log.Println(poster, posterDTO)

	switch poster.Category {
	case PropertyFlat:
		flat, err := uc.uow.Posters().GetFlatByPropetyID(ctx, poster.PropertyID)
		if err != nil {
			return nil, fmt.Errorf("uc.PosterRepo.GetFlatByPropetyID: %w", err)
		}

		flatDTO := dto.FlatToFlatFlatDTO(flat)
		posterDTO.Flat = flatDTO

	case PropertyHouse:
	default:
		return nil, fmt.Errorf("no such category %w", entity.ServiceError)
	}

	photo.MakeUrlsFromPaths(posterDTO, config.Config.PublicHost, config.Config.Bucket)

	return posterDTO, nil
}

func (uc *PosterUseCase) CreateFlatPoster(ctx context.Context, poster *dto.PosterInputFlatDTO) (*dto.CreatedPoster, error) {
	var createdPoster *dto.CreatedPoster
	post := dto.PosterInputFlatDTOtoPosterInput(poster)
	post.Alias = alias.GenerateAlias(post)

	createFlat := dto.PosterInputFlatDTOtoFlatInput(poster)

	photo.MakePhotoPathsForPoster(post)

	keys := make([]string, 0, len(post.Images))

	err := uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {
		buildingID, err := r.Posters().CreateBuilding(ctx, post)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.CreateBuilding: %w", err)
		}

		propertyID, err := r.Posters().CreateProperty(ctx, post, buildingID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.CreateProperty: %w", err)
		}

		posterID, err := r.Posters().Create(ctx, post, propertyID)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.Create: %w", err)
		}

		createFlat.PropertyID = propertyID
		err = r.Posters().InsertFlat(ctx, createFlat)
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

func (uc *PosterUseCase) uploadPhoto(ctx context.Context, photoPoster entity.PhotoInput) (string, error) {
	file, err := photoPoster.FileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("photoPoster.FileHeader.Open: %w", err)
	}
	defer file.Close()

	if !validator.ValidatePhoto(photoPoster.FileHeader) {
		return "", entity.NewValidationError("photo")
	}

	key := photo.GetKeyFromPath(photoPoster.Path)
	size := photoPoster.FileHeader.Size
	contentType := photoPoster.FileHeader.Header.Get("Content-Type")

	if err := uc.file.Upload(ctx, key, file, size, contentType); err != nil {
		return "", fmt.Errorf("uc.file.Upload: %w", err)
	}

	return key, nil
}
