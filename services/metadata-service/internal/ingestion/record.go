package ingestion

import (
	"encoding/json"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type SatelliteSourceRecord struct {
	Source    models.Source
	FetchedAt time.Time

	Raw json.RawMessage
}

type Row map[string]string
