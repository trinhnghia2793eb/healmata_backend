package model

import "time"

type OtpRequest struct {
	ID             string     `json:"id"`
	Identifier     string     `json:"identifier"`
	OtpHash        string     `json:"otpHash"`
	Purpose        string     `json:"purpose"`
	Attempts       int        `json:"attempts"`
	ExpiresAt      time.Time  `json:"expiresAt"`
	VerifiedAt     *time.Time `json:"verifiedAt,omitempty"`
	ResetTokenHash *string    `json:"resetTokenHash,omitempty"`
	TokenExpiresAt *time.Time `json:"tokenExpiresAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
}
