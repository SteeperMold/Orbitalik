package models

import "time"

type AccessToken string

type RefreshToken struct {
	Value     string
	Hash      string
	ExpiresAt time.Time
}

type StoredRefreshToken struct {
	ID        int
	UserID    int
	Hash      string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
