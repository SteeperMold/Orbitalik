package models

type Satellite struct {
	NoradID           int64  `json:"norad_id,omitempty"`
	CosparID          string `json:"cospar_id,omitempty"`
	Name              string `json:"name,omitempty"`
	ObjectType        string `json:"object_type,omitempty"`
	MissionType       string `json:"mission_type,omitempty"`
	OrbitRegime       string `json:"orbit_regime,omitempty"`
	Operator          string `json:"operator,omitempty"`
	OperationalStatus string `json:"operational_status,omitempty"`
}
