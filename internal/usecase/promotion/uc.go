package promotion

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	yoopay "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/payment"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	yoowebhook "github.com/rvinnie/yookassa-sdk-go/yookassa/webhook"
)

type PromotionUseCase struct {
	uow     usecase.UnitOfWork
	payment *yoopay.YookassaRepo
}

func NewPromotionUseCase(uow usecase.UnitOfWork) *PromotionUseCase {
	return &PromotionUseCase{
		uow:     uow,
		payment: yoopay.NewYookassaRepo(),
	}
}

func (uc *PromotionUseCase) CreatePaymentOrder(ctx context.Context, promotionCode string, posterID int, userID int) (*dto.PaymentDTO, error) { // TOFO: сделать валидацию на пользователя
	promo, err := uc.uow.Promotion().GetByCode(ctx, promotionCode)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Promotion().GetByCode: %w", err)
	}
	_, err = uc.uow.Promotion().GetActiveByPosterID(ctx, posterID)
	if err != nil {
		if !errors.Is(err, entity.NotFoundError) {
			return nil, fmt.Errorf("uc.uow.Promotion().GetByCode: %w", err)
		}
	} else {
		return nil, fmt.Errorf("uc.uow.Promotion().GetActiveByPosterID: %w", entity.AlredyExitError)
	}
	payment, err := uc.payment.CreatePayment(ctx, promo.Price, "Оплата продвижения объявления")
	if err != nil {
		return nil, fmt.Errorf("uc.payment.CreatePayment: %w", err)
	}
	endsAt := time.Now().Add(time.Hour * 24 * time.Duration(promo.DurationDays))
	_, err = uc.uow.Promotion().Create(ctx, dto.CreatePromotionDTO{
		PosterID:      posterID,
		PromotionCode: promotionCode,
		PaymentID:     &payment.PaymentID,
		EndsAt:        endsAt,
		AmountPaid:    promo.Price,
		UserID:        userID,
	})
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Promotion().Create: %w", err)
	}
	return payment, nil
}

func (uc *PromotionUseCase) ActivatePromotion(ctx context.Context, data yoowebhook.WebhookEvent[yoopayment.Payment]) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "ActivatePromotion")
	log.Info(data)
	switch data.Event {
	case yoowebhook.EventPaymentWaitingForCapture:
		log.Infof("Process payment: %s", data.Object.ID)
		prom, err := uc.uow.Promotion().GetByPaymentID(ctx, data.Object.ID)
		if err != nil {
			return fmt.Errorf("uc.uow.Promotion().GetByPaymentID: %w", err)
		}
		if prom.Status != string(entity.ActiveStatus) {
			err = uc.uow.Promotion().Activate(ctx, data.Object.ID, time.Now())
			if err != nil {
				return fmt.Errorf("uc.uow.Promotion().Start: %w", err)
			}
		}
		go func(paymentData yoopayment.Payment) {
			err := uc.payment.ConfirmPayment(context.Background(), paymentData)
			if err != nil {
				log.Errorf("uc.payment.ConfirmPayment: %s", err)
			}
		}(data.Object)
		return nil
	case yoowebhook.EventPaymentCanceled:
		err := uc.uow.Promotion().UpdateStatus(ctx, data.Object.ID, string(entity.CanceledStatus))
		if err != nil {
			return fmt.Errorf("uc.uow.Promotion().UpdateStatus: %w", err)
		}
		go func(paymentData yoopayment.Payment) {
			err := uc.payment.CancelPayment(context.Background(), paymentData.ID)
			if err != nil {
				log.Errorf("uc.payment.ConfirmPayment: %s", err)
			}
		}(data.Object)
		return nil
	default:
		return nil
	}

}

func (uc *PromotionUseCase) CheckPaymentStatus(ctx context.Context, userID int, paymentID string) (yoopayment.Status, error) {
	promo, err := uc.uow.Promotion().GetByPaymentID(ctx, paymentID)
	if err != nil {
		return "", fmt.Errorf("uc.uow.Promotion().GetByPaymentID: %w", err)
	}
	if promo.UserID != userID {
		return "", entity.NotFoundError
	}
	p, err := uc.payment.GetPayment(ctx, paymentID)
	if err != nil {
		return "", fmt.Errorf("uc.payment.GetPayment: %w", err)
	}
	return p.Status, nil

}
