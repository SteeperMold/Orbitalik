package models

import "time"

type SatelliteMetadata struct {
	NoradID  int64   `json:"norad_id"`
	CosparID *string `json:"cospar_id,omitempty"`

	Name    string   `json:"name"`
	Aliases []string `json:"aliases,omitempty"`

	ObjectType  string `json:"object_type"`
	MissionType string `json:"mission_type"`
	OrbitRegime string `json:"orbit_regime"`

	Operator      *string `json:"operator,omitempty"`
	Owner         *string `json:"owner,omitempty"`
	Constellation *string `json:"constellation,omitempty"`

	LaunchDate    *time.Time `json:"launch_date,omitempty"`
	LaunchSite    *string    `json:"launch_site,omitempty"`
	LaunchVehicle *string    `json:"launch_vehicle,omitempty"`

	OperationalStatus string      `json:"operational_status"`
	Frequencies       Frequencies `json:"frequencies,omitempty"`

	Sources Sources

	UpdatedAt time.Time `json:"updated_at"`
}

type Frequencies []Frequency

type Frequency struct {
	Direction    string   `json:"direction"`
	FrequencyMHz float64  `json:"frequency_mhz"`
	BandwidthKHz *float64 `json:"bandwidth_khz,omitempty"`
	Modulation   string   `json:"modulation"`
	Mode         string   `json:"mode"`
}

type Sources []FieldSource

type FieldSource struct {
	Field     string    `json:"field"`
	Sources   []string  `json:"sources"`
	FetchedAt time.Time `json:"fetched_at"`
}
