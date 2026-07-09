package model

import "time"

type User struct {
	ID                  string
	FullName            string
	Email               string
	Phone               string
	PasswordHash        string
	Status              string
	FirstSetupCompleted bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
