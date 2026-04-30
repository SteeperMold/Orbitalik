package models

type ListFilter struct {
	ObjectType        *ObjectType
	MissionType       *MissionType
	OperationalStatus *OperationalStatus
	OrbitRegime       *OrbitRegime
	Constellation     *string

	PageSize  uint32
	PageToken string
}
