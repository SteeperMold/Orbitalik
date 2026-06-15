package celestrak

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(r ingestion.Row) (json.RawMessage, error) {
	now := time.Now()

	noradID, err := strconv.Atoi(r["norad_id"])
	if err != nil {
		return nil, err
	}

	meta := &models.SatelliteMetadataPartial{
		NoradID:   noradID,
		FetchedAt: now,
		Source: models.SourceAttribution{
			Source:         models.SourceCelestrak,
			SourceRecordID: r["norad_id"],
			FetchedAt:      now,
		},
	}

	if v := strings.TrimSpace(r["cospar_id"]); v != "" {
		meta.CosparID = &v
	}

	if v := strings.TrimSpace(r["name"]); v != "" {
		meta.Name = &v
	}

	if v := strings.TrimSpace(r["owner"]); v != "" {
		meta.Owner = &v
	}

	if t := parseDate(r["launch_date"]); t != nil {
		meta.LaunchDate = t
	}

	if v := strings.TrimSpace(r["launch_site"]); v != "" {
		meta.LaunchSite = &v
	}

	if v := deriveObjectType(r["name"]); v != models.ObjectTypeUnspecified {
		meta.ObjectType = &v
	}

	if v := deriveStatus(r["flags"], r["decay_date"]); v != models.OperationalStatusUnspecified {
		meta.OperationalStatus = &v
	}

	if v := deriveOrbit(r["apogee"]); v != models.OrbitRegimeUnspecified {
		meta.OrbitRegime = &v
	}

	return json.Marshal(meta)
}

func parseDate(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}

	return &t
}

func deriveStatus(flags, decay string) models.OperationalStatus {
	flags = strings.TrimSpace(flags)

	if strings.Contains(flags, "D") || strings.TrimSpace(decay) != "" {
		return models.OperationalStatusDecayed
	}

	return models.OperationalStatusUnknown
}

func deriveObjectType(name string) models.ObjectType {
	n := strings.ToUpper(name)

	switch {
	case strings.Contains(n, "R/B"):
		return models.ObjectTypeRocketBody
	case strings.Contains(n, "DEB"):
		return models.ObjectTypeDebris
	default:
		return models.ObjectTypePayload
	}
}

func deriveOrbit(apogeeStr string) models.OrbitRegime {
	apogeeStr = strings.TrimSpace(apogeeStr)
	if apogeeStr == "" {
		return models.OrbitRegimeUnspecified
	}

	apogee, err := strconv.Atoi(apogeeStr)
	if err != nil {
		return models.OrbitRegimeUnspecified
	}

	switch {
	case apogee < 2000:
		return models.OrbitRegimeLEO
	case apogee < 35786:
		return models.OrbitRegimeMEO
	case apogee <= 36000:
		return models.OrbitRegimeGEO
	default:
		return models.OrbitRegimeHEO
	}
}
