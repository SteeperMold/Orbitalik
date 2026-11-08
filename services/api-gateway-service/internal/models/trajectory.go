package models

import "time"

type TrajectoryResult struct {
	ComputationMetadata TrajectoryComputationMetadata `json:"computation_metadata"`
	States              []StateVector                 `json:"states"`
	Range               TimeRange                     `json:"range"`
	Sampling            SamplingOptions               `json:"sampling"`
}

type StateVector struct {
	Time     time.Time `json:"time"`
	Teme     *Vector3  `json:"teme,omitempty"`
	Ecef     *Vector3  `json:"ecef,omitempty"`
	Geodetic *Geodetic `json:"geodetic,omitempty"`
}

type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type SamplingOptions struct {
	StepSeconds int `json:"step_seconds"`
}

type ObserverTrajectoryPoint struct {
	Time      time.Time `json:"time"`
	Azimuth   *float64  `json:"azimuth,omitempty"`
	Elevation *float64  `json:"elevation,omitempty"`
	Range     *float64  `json:"range,omitempty"`
}

type ObserverTrajectoryResult struct {
	ComputationMetadata TrajectoryComputationMetadata `json:"computation_metadata"`
	Points              []ObserverTrajectoryPoint     `json:"points"`
	Range               TimeRange                     `json:"range"`
	Sampling            SamplingOptions               `json:"sampling"`
}
