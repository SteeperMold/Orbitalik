package models

import (
	"time"

	"github.com/google/uuid"
)

type TLE struct {
	ID            uuid.UUID
	NoradID       int
	SatelliteName string
	Line1         string
	Line2         string
	Epoch         time.Time
	FetchedAt     time.Time
}
