package ucs

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
	noradID, err := strconv.Atoi(clean(r["norad_id"]))
	if err != nil {
		return nil, err
	}

	now := time.Now()

	meta := &models.SatelliteMetadataPartial{
		NoradID: noradID,
		Name:    getPtr(clean(r["name"])),

		ObjectType:  getPtr(models.ObjectTypePayload),
		MissionType: getPtr(mapMissionType(r["users"], r["purpose"])),
		OrbitRegime: getPtr(mapOrbit(r["orbit_class"])),

		OperationalStatus: getPtr(models.OperationalStatusUnknown),
		FetchedAt:         now,
	}

	meta.Operator = getPtr(r["operator"])
	meta.Owner = getPtr(r["owner"])
	meta.LaunchSite = getPtr(r["launch_site"])
	meta.LaunchVehicle = getPtr(r["launch_vehicle"])
	meta.CosparID = getPtr(r["cospar"])

	meta.LaunchDate = parseDate(r["launch_date"])

	meta.Aliases = extractAliases(r["aliases"])

	return json.Marshal(meta)
}

func getPtr[T ~string](s T) *T {
	if s == "" {
		return nil
	}
	return &s
}

func parseDate(s string) *time.Time {
	s = clean(s)
	if s == "" {
		return nil
	}

	layouts := []string{
		"1/2/06",
		"1/2/2006",
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return &t
		}
	}

	return nil
}

func extractAliases(name string) []string {
	name = clean(name)

	start := strings.Index(name, "(")
	end := strings.LastIndex(name, ")")

	if start == -1 || end == -1 || end <= start {
		return []string{}
	}

	raw := name[start+1 : end]
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")

	aliases := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			aliases = append(aliases, p)
		}
	}

	return aliases
}

func mapMissionType(users, purpose string) models.MissionType {
	s := strings.ToLower(clean(users) + " " + clean(purpose))

	switch {
	case strings.Contains(s, "communication"):
		return models.MissionTypeCommunications
	case strings.Contains(s, "earth observation"):
		return models.MissionTypeEarthObservation
	case strings.Contains(s, "navigation"):
		return models.MissionTypeNavigation
	case strings.Contains(s, "science"):
		return models.MissionTypeScience
	case strings.Contains(s, "weather"):
		return models.MissionTypeWeather
	case strings.Contains(s, "technology"):
		return models.MissionTypeTechDemo
	case strings.Contains(s, "amateur"):
		return models.MissionTypeAmateur
	default:
		return models.MissionTypeUnspecified
	}
}

func mapOrbit(o string) models.OrbitRegime {
	switch strings.ToUpper(clean(o)) {
	case "LEO":
		return models.OrbitRegimeLEO
	case "MEO":
		return models.OrbitRegimeMEO
	case "GEO":
		return models.OrbitRegimeGEO
	case "HEO":
		return models.OrbitRegimeHEO
	default:
		return models.OrbitRegimeUnspecified
	}
}
