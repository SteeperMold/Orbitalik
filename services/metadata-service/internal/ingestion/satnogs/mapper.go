package satnogs

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(r ingestion.Row) (json.RawMessage, error) {
	norad, err := strconv.Atoi(r["norad_id"])
	if err != nil {
		return nil, err
	}

	now := time.Now()

	meta := &models.SatelliteMetadataPartial{
		NoradID: norad,

		ObjectType: getPtr(models.ObjectTypePayload),

		MissionType: getPtr(detectMissionType(r["mode"])),

		OrbitRegime: getPtr(models.OrbitRegimeUnspecified),

		OperationalStatus: getPtr(models.OperationalStatusUnknown),
		FetchedAt:         now,
	}

	freqs := buildFrequencies(r)
	meta.Frequencies = freqs

	return json.Marshal(meta)
}

func buildFrequencies(r ingestion.Row) []models.Frequency {
	out := make([]models.Frequency, 0)

	down := parseHz(r["downlink_low"])
	up := parseHz(r["uplink_low"])

	if down > 0 {
		out = append(out, models.Frequency{
			Direction:    models.FrequencyDirectionDownlink,
			FrequencyMHz: down,
			Mode:         r["mode"],
			Modulation:   r["mode"],
		})
	}

	if up > 0 {
		out = append(out, models.Frequency{
			Direction:    models.FrequencyDirectionUplink,
			FrequencyMHz: up,
			Mode:         r["mode"],
			Modulation:   r["mode"],
		})
	}

	return out
}

func parseHz(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v == 0 {
		return 0
	}
	return v / 1e6
}

func detectMissionType(mode string) models.MissionType {
	switch mode {
	case "AFSK", "FM", "BPSK":
		return models.MissionTypeAmateur
	case "APT", "LRPT":
		return models.MissionTypeWeather
	default:
		return models.MissionTypeUnspecified
	}
}

func getPtr[T ~string](s T) *T {
	if s == "" {
		return nil
	}
	return &s
}
