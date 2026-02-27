package entity

import "time"

type RefreshToken struct {
	ID        int
	TokenID   string
	UserID    int
	ExpiresAt time.Time
}
