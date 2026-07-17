package model

import "time"

type SocialAccounts struct {
	ID             string
	UserID         string
	Provider       string
	ProviderUserID string
	ProviderEmail  string
	CreatedAt      time.Time
}
