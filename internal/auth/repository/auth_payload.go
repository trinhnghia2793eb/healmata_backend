package repository

import "time"

type CreateUserPayload struct {
	FullName     string
	Email        *string
	Phone        *string
	PasswordHash string
}

type CreateRefreshTokenPayload struct {
	UserID    string
	TokenHash string
	DeviceID  string
	ExpiresAt time.Time
}

type CreateUserSessionPayload struct {
	UserID         string
	RefreshTokenID string
	DeviceID       string
	Platform       string
	IPAddress      string
	UserAgent      string
}
