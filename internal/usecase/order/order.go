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
	uow   usecase.UnitOfWork
	file  usecase.FileRepo
	agent usecase.SupportAgent
}

func NewOrderUseCase(uow usecase.UnitOfWork, file usecase.FileRepo, agent usecase.SupportAgent) *OrderUseCase {
	return &OrderUseCase{
		uow:   uow,
		file:  file,
		agent: agent,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, order *dto.OrderDTO) error {
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

		orderID, err := r.Order().Create(ctx, createOrder)
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

		err = r.Order().InsertPhotos(ctx, orderID, createOrder.Images)
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

func (uc *OrderUseCase) GetSupportAgentResponse(ctx context.Context, userPrompt string) (*dto.ChatResult, error) {
	systemPrompt := `
	Ты — помощник службы поддержки ДомДели, сервиса аренды жилья.

	Твоя задача — помочь пользователю решить проблему быстро, кратко и по существу.

	Правила:
	- Отвечай на языке пользователя.
	- Если знаешь ответ — дай его сразу.
	- Если не знаешь — честно скажи, что не знаешь.
	- Не выдумывай факты, условия, цены, адреса, статусы брони и правила.
	- Не предлагай обратиться к человеку без причины.
	- Если для ответа не хватает данных, явно укажи, чего именно не хватает.
	- Если проблема требует ручной проверки, доступа к заказу, платежам, личным данным, спору или технической диагностики — пометь это как обращение к человеку.

	Формат ответа:
	Всегда возвращай только JSON без лишнего текста:

	{
	"status": "ok|error|human_required",
	"answer": "краткий ответ пользователю",
	"missing_info": ["что нужно уточнить или проверить"],
	"reason": "краткая причина, если status = error или human_required"
	}

	Правила для status:
	- ok — если вопрос решен.
	- error — если ответ неизвестен, данных недостаточно для безопасного ответа или есть противоречие.
	- human_required — если нужна ручная проверка или участие сотрудника.

	Дополнительные правила:
	- Если status = ok, missing_info должен быть пустым массивом.
	- Если status = error, кратко объясни, что не так.
	- Если status = human_required, объясни, почему нужен человек.
	- Если пользователь не указал важные данные, перечисли их в missing_info.
	- Не давай длинных объяснений.
	- Не используй markdown, только JSON.
		`

	response, err := uc.agent.Chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("uc.agent.Chat: %w", err)
	}

	return response, nil

}

func (uc *OrderUseCase) GetUserOrders(ctx context.Context, userID int) (*dto.OrdersResponse, error) {
	orders, err := uc.uow.Order().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	orderDTOs := dto.ToOrderPreview(orders)

	response := &dto.OrdersResponse{
		Len:    len(orderDTOs),
		Orders: orderDTOs,
	}

	return response, nil
}
