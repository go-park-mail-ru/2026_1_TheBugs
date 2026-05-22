package entity

type PromotionStatus string

var (
	ActiveStatus   PromotionStatus = "active"
	PendingStatus  PromotionStatus = "pending"
	CanceledStatus PromotionStatus = "cancelled"
	ExpieredStatus PromotionStatus = "expiered"
)

type Promotion struct {
	ID           int     `db:"id"`
	Code         string  `db:"code"`
	Price        float32 `db:"price"`
	DurationDays int     `db:"duration_days"`
	Name         string  `db:"name"`
	Description  string  `db:"description"`
}

type PosterPromotion struct {
	ID         int    `db:"id"`
	PosterID   int    `db:"poster_id"`
	UserID     int    `db:"user_id"`
	PromotinID int    `db:"promotion_id"`
	Status     string `db:"status"`
	PaymentID  string `db:"payment_id"`
}
