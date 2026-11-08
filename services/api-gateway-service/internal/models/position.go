package models

import "time"

type PositionResult struct {
	ComputationMetadata TrajectoryComputationMetadata `json:"computation_metadata"`
	Time                time.Time                     `json:"time"`
	Teme                *Vector3                      `json:"teme,omitempty"`
	Ecef                *Vector3                      `json:"ecef,omitempty"`
	Geodetic            *Geodetic                     `json:"geodetic,omitempty"`
}

type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Geodetic struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
	Alt float64 `json:"altitude"`
}

type LookAnglesResult struct {
	ComputationMetadata TrajectoryComputationMetadata `json:"computation_metadata"`
	Time                time.Time                     `json:"time"`
	Azimuth             *float64                      `json:"azimuth,omitempty"`
	Elevation           *float64                      `json:"elevation,omitempty"`
	Range               *float64                      `json:"range,omitempty"`
}

type TrajectoryComputationMetadata struct {
	PropagationModel string       `json:"propagation_model"`
	ComputationTime  time.Time    `json:"computation_time"`
	NoradID          int64        `json:"norad_id"`
	SatelliteName    string       `json:"satellite_name"`
	TleEpoch         time.Time    `json:"tle_epoch"`
	Units            UnitSettings `json:"units"`
}

type UnitSettings struct {
	DistanceUnit string `json:"distance_unit"`
	AngleUnit    string `json:"angle_unit"`
}
