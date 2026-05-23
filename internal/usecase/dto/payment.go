package dto

type PaymentDTO struct {
	PaymentID       string `json:"payment_id"`
	ConfirmationUrl string `json:"confirmation_url"`
}
