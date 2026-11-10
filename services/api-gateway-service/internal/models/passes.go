package models

import "time"

type PassesResult struct {
	ComputationMetadata PassComputationMetadata `json:"computation_metadata"`
	Passes              []Pass                  `json:"passes"`
	Range               TimeRange               `json:"range"`
}

type NextPassesResult struct {
	ComputationMetadata PassComputationMetadata `json:"computation_metadata"`
	Passes              []Pass                  `json:"passes"`
}

type PassComputationMetadata struct {
	PropagationModel    string       `json:"propagation_model"`
	ComputationTime     time.Time    `json:"computation_time"`
	NoradIDs            []int64      `json:"norad_ids"`
	SatelliteNames      []string     `json:"satellite_names"`
	TleEpoch            time.Time    `json:"tle_epoch"`
	Units               UnitSettings `json:"units"`
	SatellitesEvaluated int          `json:"satellites_evaluated"`
	PassesFound         int          `json:"passes_found"`
}

type SatelliteIdentifier struct {
	NoradID *int64  `json:"norad_id,omitempty"`
	Name    *string `json:"name,omitempty"`
}

type Pass struct {
	Satellite           SatelliteIdentifier `json:"satellite"`
	AOS                 time.Time           `json:"aos"`
	AOSAzimuth          float64             `json:"aos_azimuth"`
	MaxElevationTime    time.Time           `json:"max_elevation_time"`
	MaxElevation        float64             `json:"max_elevation"`
	MaxElevationAzimuth float64             `json:"max_elevation_azimuth"`
	LOS                 time.Time           `json:"los"`
	LOSAzimuth          float64             `json:"los_azimuth"`
	DurationSeconds     int                 `json:"duration_seconds"`
}
