package models

import "time"

type TLE struct {
	ID            int
	NoradID       int
	SatelliteName string
	Line1         string
	Line2         string
	Epoch         time.Time
	FetchedAt     time.Time
}
