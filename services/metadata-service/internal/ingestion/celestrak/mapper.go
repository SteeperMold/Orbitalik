package celestrak

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/filesource"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(r filesource.Row) (json.RawMessage, error) {
	now := time.Now()

	noradID, err := strconv.Atoi(r["norad_id"])
	if err != nil {
		return nil, err
	}

	meta := &models.SatelliteMetadata{
		NoradID: noradID,

		UpdatedAt: now,

		Sources: []models.SourceAttribution{
			{
				Source:         models.SourceCelestrak,
				SourceRecordID: r["norad_id"],
				FetchedAt:      now,
			},
		},
	}

	if v := mapObjectType(r["object_type"]); v != models.ObjectTypeUnspecified {
		meta.ObjectType = v
	}

	if v := mapStatus(r["status"]); v != models.OperationalStatusUnspecified {
		meta.OperationalStatus = v
	}

	return json.Marshal(meta)
}

func mapStatus(s string) models.OperationalStatus {
	switch strings.TrimSpace(s) {
	case "+":
		return models.OperationalStatusActive
	case "-":
		return models.OperationalStatusDecayed
	default:
		return models.OperationalStatusUnknown
	}
}

func mapObjectType(s string) models.ObjectType {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "PAY":
		return models.ObjectTypePayload
	case "R/B":
		return models.ObjectTypeRocketBody
	case "DEB":
		return models.ObjectTypeDebris
	default:
		return models.ObjectTypeUnknown
	}
}
