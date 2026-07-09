package model

import "time"

type UserSessions struct {
	ID             string
	UserID         string
	RefreshTokenID string
	DeviceID       string
	Platform       string
	IPAdress       string
	UserAgent      string
	LastActiveAt   time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
