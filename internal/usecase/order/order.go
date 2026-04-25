package order

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type OrderUseCase struct {
	uow    usecase.UnitOfWork
	file   usecase.FileRepo
	search usecase.SearchRepo
}

func NewOrderUseCase(uow usecase.UnitOfWork, file usecase.FileRepo, search usecase.SearchRepo) *OrderUseCase {
	return &OrderUseCase{
		uow:    uow,
		file:   file,
		search: search,
	}
}

func (uc *OrderUseCase) CreateFlatPoster(ctx context.Context, order *dto.OrderDTO) error {
	/*err := validator.ValidatePosterBase(poster)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePosterBase: %w", err)
	}*/

	err := validator.ValidatePhotos(order.Images)
	if err != nil {
		return fmt.Errorf("validator.ValidatePhotos: %w", err)
	}
	//validator.SanitizePosterInput(order)

	createOrder := dto.OrderDTOtoOrder(order)

	dto.MakePhotoPathsForOrder(createOrder)

	keys := make([]string, 0, len(createOrder.Images))

	err = uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {

		posterID, err := r.Posters().Create(ctx, createOrder)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.Create: %w", err)
		}

		for _, photoPoster := range createOrder.Images {
			key, err := uc.uploadPhoto(ctx, photoPoster)
			if err != nil {
				return fmt.Errorf("uc.uploadPhoto: %w", err)
			}
			keys = append(keys, key)
		}

		err = r.Posters().InsertPhotos(ctx, posterID, createOrder.Images)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.InsertPhotos: %w", err)
		}

		if len(createOrder.Images) > 0 {
			err = r.Posters().InsertMainPhoto(ctx, posterID, post.Images[0].Path)
			if err != nil {
				return fmt.Errorf("r.Posters().InsertMainPhoto: %w", err)
			}
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

func (uc *OrderUseCase) cleanUploadedFiles(ctx context.Context, keys []string) error {
	var resultErr error
	for _, key := range keys {
		err := uc.file.Delete(ctx, key)
		if err != nil {
			resultErr = fmt.Errorf("uc.file.Delete: %w", err)
		}
	}

	return resultErr
}

func (uc *OrderUseCase) uploadPhoto(ctx context.Context, photoPoster dto.PhotoInput) (string, error) {
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
