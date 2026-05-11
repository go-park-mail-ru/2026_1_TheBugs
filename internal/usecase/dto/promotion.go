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

type UserPromotionDTO struct {
	PromotionID int       `db:"promotion_id" json:"promotion_id"`
	EndsAt      time.Time `db:"ends_at" json:"ends_at"`
	PosterID    int       `db:"poster_id" json:"poster_id"`
	Status      string    `db:"status" json:"status"`
}
type UserPromotionsDTO struct {
	Lenght     int                `json:"lenght"`
	Promotions []UserPromotionDTO `json:"promotions"`
}

type CreatePromotionDTO struct {
	PosterID      int
	UserID        int
	PromotionCode string
	PaymentID     *string
	EndsAt        time.Time
	AmountPaid    float32
}
