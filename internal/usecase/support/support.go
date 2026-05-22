package support

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type SupportUseCase struct {
	uow    usecase.UnitOfWork
	file   usecase.FileRepo
	sender usecase.MailSender
}

func NewSupportUseCase(uow usecase.UnitOfWork, file usecase.FileRepo) *SupportUseCase {
	return &SupportUseCase{
		uow:  uow,
		file: file,
	}
}

func (uc *SupportUseCase) CreateOrder(ctx context.Context, order *dto.OrderDTO) error {
	/*err := validator.ValidatePosterBase(poster)
	if err != nil {
		return nil, fmt.Errorf("validator.ValidatePosterBase: %w", err)
	}*/

	err := validator.ValidatePhotos(order.Images)
	if err != nil {
		return fmt.Errorf("validator.ValidatePhotos: %w", err)
	}
	//validator.SanitizePosterInput(order)
	user, err := uc.uow.Users().GetByEmail(ctx, order.Email)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().GetByEmail: %w", err)
	}
	//validator.SanitizePosterInput(order)

	createOrder := dto.OrderDTOtoOrder(order, user.ID)

	dto.MakePhotoPathsForOrder(createOrder)

	keys := make([]string, 0, len(createOrder.Images))

	err = uc.uow.Do(ctx, func(r usecase.UnitOfWork) error {

		orderID, err := r.Support().Create(ctx, createOrder)
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

		err = r.Support().InsertPhotos(ctx, orderID, createOrder.Images)
		if err != nil {
			return fmt.Errorf("uc.PosterRepo.InsertPhotos: %w", err)
		}

		return nil
	})

	if err != nil {
		cleanErr := uc.cleanUploadedFiles(ctx, keys)
		if cleanErr != nil {
			return fmt.Errorf("uc.cleanUploadedFiles: %w", cleanErr)
		}

		return fmt.Errorf("uc.uow.Do: %w", err)
	}

	return nil
}

func (uc *SupportUseCase) cleanUploadedFiles(ctx context.Context, keys []string) error {
	var resultErr error
	for _, key := range keys {
		err := uc.file.Delete(ctx, key)
		if err != nil {
			resultErr = fmt.Errorf("uc.file.Delete: %w", err)
		}
	}

	return resultErr
}

func (uc *SupportUseCase) uploadPhoto(ctx context.Context, photoPoster dto.PhotoInput) (string, error) {
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

func (uc *SupportUseCase) GetUserOrders(ctx context.Context, userID int) (*dto.OrdersResponse, error) {
	user, err := uc.uow.Users().GetAdminByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetAdminByID: %w", err)
	}

	var orders []entity.Order
	if user.IsAdmin {
		orders, err = uc.uow.Support().GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("uc.uow.Order().GetAll: %w", err)
		}
	} else {
		orders, err = uc.uow.Support().GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	orderDTOs := dto.ToOrderPreview(orders)

	response := &dto.OrdersResponse{
		Len:    len(orderDTOs),
		Orders: orderDTOs,
	}

	return response, nil
}

func (uc *SupportUseCase) GetOrderByID(ctx context.Context, userID int, orderID int) (*dto.OrderFullDTO, error) {
	user, err := uc.uow.Users().GetAdminByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetByEmail: %w", err)
	}

	order, err := uc.uow.Support().GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Order().GetByID: %w", err)
	}

	if !user.IsAdmin && order.UserID != userID {
		return nil, entity.ServiceError
	}

	photos, err := uc.uow.Support().GetOrderImages(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Support().GetPhotosByOrderID: %w", err)
	}

	order.Photos = photos

	orderDTO := dto.OrderToOrderFullDTO(order)
	dto.MakeOrderUrlsFromPaths(orderDTO, config.Config.PublicHost, config.Config.Bucket)

	return orderDTO, nil
}

func (uc *SupportUseCase) AnswerOrder(ctx context.Context, adminID int, orderID int, answer string) error {

	user, err := uc.uow.Users().GetAdminByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().GetAdminByEmail: %w", err)
	}

	if !user.IsAdmin {
		return entity.ServiceError
	}

	order, err := uc.uow.Support().GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("uc.uow.Support().GetByID: %w", err)
	}

	client, err := uc.uow.Users().GetByID(ctx, order.UserID)
	if err != nil {
		return fmt.Errorf("uc.uow.Users().GetByID: %w", err)
	}

	err = uc.sender.SendAnswer(ctx, client.Email, orderID, answer)
	if err != nil {
		return fmt.Errorf("uc.sender.SendAnswer: %w", err)
	}

	err = uc.uow.Support().FinishOrder(ctx, orderID, adminID)
	if err != nil {
		return fmt.Errorf("uc.uow.Order().FinishOrder: %w", err)
	}

	return nil
}
