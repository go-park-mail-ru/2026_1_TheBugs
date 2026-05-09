package dto

import "time"

type PromotionDTO struct {
	ID           int
	Code         string
	Price        float32
	DurationDays int
	Name         string
	Description  string
}

type CreatePromotionDTO struct {
	PosterID      int
	UserID        int
	PromotionCode string
	PaymentID     *string
	EndsAt        time.Time
	AmountPaid    float32
}
