package models

import "time"

type Tle struct {
	NoradID       uint32    `json:"norad_id"`
	SatelliteName string    `json:"satellite_name"`
	Line1         string    `json:"line1"`
	Line2         string    `json:"line2"`
	Epoch         time.Time `json:"epoch"`
	FetchedAt     time.Time `json:"fetched_at"`
}
