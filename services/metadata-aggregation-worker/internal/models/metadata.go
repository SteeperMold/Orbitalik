package models

import (
	"encoding/json"
	"time"
)

type SatelliteIngestRecord struct {
	ID int64

	NoradID  int
	CosparID *string

	Source Source

	Payload json.RawMessage

	FetchedAt time.Time
	StoredAt  time.Time
}

type SatelliteMetadata struct {
	NoradID  int
	CosparID *string

	Name    string
	Aliases []string

	ObjectType  ObjectType
	MissionType MissionType
	OrbitRegime OrbitRegime

	Operator      *string
	Owner         *string
	Constellation *string

	LaunchDate    *time.Time
	LaunchSite    *string
	LaunchVehicle *string

	OperationalStatus OperationalStatus

	Frequencies []Frequency

	Sources []FieldSource

	UpdatedAt time.Time
}

type SatelliteMetadataPartial struct {
	NoradID  int     `json:"norad_id"`
	CosparID *string `json:"cospar_id"`

	Name    *string  `json:"name"`
	Aliases []string `json:"aliases"`

	ObjectType  *ObjectType  `json:"object_type"`
	MissionType *MissionType `json:"mission_type"`
	OrbitRegime *OrbitRegime `json:"orbit_regime"`

	Operator      *string `json:"operator"`
	Owner         *string `json:"owner"`
	Constellation *string `json:"constellation"`

	LaunchDate    *time.Time `json:"launch_date"`
	LaunchSite    *string    `json:"launch_site"`
	LaunchVehicle *string    `json:"launch_vehicle"`

	OperationalStatus *OperationalStatus `json:"operational_status"`

	Frequencies []Frequency `json:"frequencies"`

	Sources   []FieldSource `json:"sources"`
	FetchedAt time.Time     `json:"fetched_at"`
}

type Frequency struct {
	Direction    FrequencyDirection
	FrequencyMHz float64
	BandwidthKHz *float64
	Modulation   string // FM, CW, SSB, BPSK, etc
	Mode         string // Beacon, Transponder, etc
}

type FieldSource struct {
	Field     string
	Sources   []Source
	FetchedAt time.Time
}

type (
	ObjectType         string
	MissionType        string
	OrbitRegime        string
	OperationalStatus  string
	FrequencyDirection string
	Source             string
)

const (
	ObjectTypeUnspecified ObjectType = "unspecified"
	ObjectTypePayload     ObjectType = "payload"
	ObjectTypeRocketBody  ObjectType = "rocket_body"
	ObjectTypeDebris      ObjectType = "debris"
	ObjectTypeUnknown     ObjectType = "unknown"

	MissionTypeUnspecified      MissionType = "unspecified"
	MissionTypeCommunications   MissionType = "communications"
	MissionTypeEarthObservation MissionType = "earth_observation"
	MissionTypeNavigation       MissionType = "navigation"
	MissionTypeScience          MissionType = "science"
	MissionTypeWeather          MissionType = "weather"
	MissionTypeAmateur          MissionType = "amateur"
	MissionTypeTechDemo         MissionType = "tech_demo"

	OrbitRegimeUnspecified OrbitRegime = "unspecified"
	OrbitRegimeLEO         OrbitRegime = "leo"
	OrbitRegimeMEO         OrbitRegime = "meo"
	OrbitRegimeGEO         OrbitRegime = "geo"
	OrbitRegimeHEO         OrbitRegime = "heo"

	OperationalStatusUnspecified OperationalStatus = "unspecified"
	OperationalStatusActive      OperationalStatus = "active"
	OperationalStatusInactive    OperationalStatus = "inactive"
	OperationalStatusDecayed     OperationalStatus = "decayed"
	OperationalStatusUnknown     OperationalStatus = "unknown"

	FrequencyDirectionUnspecified FrequencyDirection = "unspecified"
	FrequencyDirectionUplink      FrequencyDirection = "uplink"
	FrequencyDirectionDownlink    FrequencyDirection = "downlink"

	SourceUnspecified Source = "unspecified"
	SourceAMSAT       Source = "amsat"
	SourceCelestrak   Source = "celestrak"
	SourceSatNOGS     Source = "satnogs"
	SourceUCS         Source = "ucs"
	SourceManual      Source = "manual"
)
