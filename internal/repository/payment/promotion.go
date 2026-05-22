package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type YookassaRepo struct {
	client         *yookassa.Client
	paymentHandler *yookassa.PaymentHandler
}

func NewYookassaRepo() *YookassaRepo {
	client := yookassa.NewClient(
		config.Config.Yookassa.AccountID,
		config.Config.Yookassa.SecretKey,
	)
	return &YookassaRepo{
		client:         client,
		paymentHandler: yookassa.NewPaymentHandler(client),
	}
}

func (r *YookassaRepo) CreatePayment(ctx context.Context, price float32, description string) (*dto.PaymentDTO, error) {
	payment, err := r.paymentHandler.CreatePayment(ctx, &yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    fmt.Sprintf("%f", price),
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: config.Config.Yookassa.RedirectURL,
		},
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("paymentHandler.CreatePayment: %w", err)
	}

	confirm := payment.Confirmation.(map[string]any)
	raw, ok := confirm["confirmation_url"]
	if !ok {
		return nil, fmt.Errorf("Not Found confirmation_url : %w", entity.ServiceError)
	}
	url, ok := raw.(string)
	if !ok {
		return nil, fmt.Errorf("Bad type url: %w", entity.ServiceError)
	}
	return &dto.PaymentDTO{ConfirmationUrl: url, PaymentID: payment.ID}, nil

}

func (r *YookassaRepo) ConfirmPayment(ctx context.Context, payment yoopayment.Payment) error {
	captured, err := r.paymentHandler.CapturePayment(ctx, &payment)
	if err != nil {
		return fmt.Errorf("paymentHandler.CapturePayment: %w", err)
	}
	log.Printf("✅ Captured %s: %s", captured.ID, captured.Status)
	return nil
}

func (r *YookassaRepo) CancelPayment(ctx context.Context, paymentID string) error {
	captured, err := r.paymentHandler.CancelPayment(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("paymentHandler.CapturePayment: %w", err)
	}
	log.Printf("[X] Canseled %s: %s", captured.ID, captured.Status)
	return nil
}

func (r *YookassaRepo) GetPayment(ctx context.Context, paymentID string) (*yoopayment.Payment, error) {
	p, err := r.paymentHandler.FindPayment(context.Background(), paymentID)
	if err != nil {
		return nil, fmt.Errorf("paymentHandler.CapturePayment: %w", err)
	}
	return p, nil
}
