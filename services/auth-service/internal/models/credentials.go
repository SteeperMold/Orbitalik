package models

import "time"

type Credentials struct {
	UserID       int
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
}
