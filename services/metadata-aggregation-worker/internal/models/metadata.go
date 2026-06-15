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

	Sources []SourceAttribution

	UpdatedAt time.Time
}

type SatelliteMetadataPartial struct {
	NoradID  int
	CosparID *string

	Name    *string
	Aliases []string

	ObjectType  *ObjectType
	MissionType *MissionType
	OrbitRegime *OrbitRegime

	Operator      *string
	Owner         *string
	Constellation *string

	LaunchDate    *time.Time
	LaunchSite    *string
	LaunchVehicle *string

	OperationalStatus *OperationalStatus

	Frequencies []Frequency

	Source    SourceAttribution
	FetchedAt time.Time
}

type Frequency struct {
	Direction    FrequencyDirection
	FrequencyMHz float64
	BandwidthKHz *float64
	Modulation   string // FM, CW, SSB, BPSK, etc
	Mode         string // Beacon, Transponder, etc
}

type SourceAttribution struct {
	Source         Source
	SourceRecordID string
	FetchedAt      time.Time
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
