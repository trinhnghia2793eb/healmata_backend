package model

import "time"

type RefreshTokens struct {
	ID        string
	UserID    string
	TokenHash string
	DeviceID  string
	ExpiresAt time.Time
	RevokedAt time.Time
	CreatedAt time.Time
}
