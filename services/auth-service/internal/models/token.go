package models

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken string

type RefreshToken struct {
	Value     string
	Hash      string
	ExpiresAt time.Time
}

type StoredRefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Hash      string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
