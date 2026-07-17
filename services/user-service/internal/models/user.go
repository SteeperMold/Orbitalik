package models

import "time"

type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
