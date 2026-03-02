package entity

import "time"

const AccessTokenType = "access"
const RefreshTokenType = "refresh"

type RefreshToken struct {
	ID        int
	TokenID   string
	UserID    int
	ExpiresAt time.Time
}
