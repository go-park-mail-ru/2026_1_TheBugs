package entity

import "time"

const AccessTokenType = "access"
const RefreshTokenType = "refresh"

type RefreshToken struct {
	ID        int       `db:"id"`
	TokenID   string    `db:"token_id"`
	UserID    int       `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}
