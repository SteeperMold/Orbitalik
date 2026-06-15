package aggregation

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
)

type item struct {
	meta *models.SatelliteMetadataPartial
	rec  models.SatelliteIngestRecord
}

func Aggregate(records []models.SatelliteIngestRecord) (*models.SatelliteMetadata, error) {
	if len(records) == 0 {
		return nil, nil
	}

	var items []item

	for _, r := range records {
		m, err := decodePayload(r)
		if err != nil {
			continue
		}
		items = append(items, item{m, r})
	}

	if len(items) == 0 {
		return nil, nil
	}

	return mergeItems(items), nil
}

func decodePayload(r models.SatelliteIngestRecord) (*models.SatelliteMetadataPartial, error) {
	var m models.SatelliteMetadataPartial
	if err := json.Unmarshal(r.Payload, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// the less the more priority
var fieldPriority = map[string]map[models.Source]int{
	"name": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"operator": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"orbit_regime": {
		models.SourceCelestrak: 2,
		models.SourceUCS:       1,
	},
	"object_type": {
		models.SourceCelestrak: 2,
		models.SourceUCS:       1,
	},
	"operational_status": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
}

func pickBest[T any](
	items []item,
	field string,
	extract func(*models.SatelliteMetadataPartial) (T, bool),
) (T, bool) {

	var (
		bestVal T
		found   bool
		bestPri = -1
		bestTs  time.Time
	)

	for _, it := range items {
		val, ok := extract(it.meta)
		if !ok {
			continue
		}

		pri := fieldPriority[field][it.rec.Source]

		if !found || pri > bestPri || (pri == bestPri && it.rec.FetchedAt.After(bestTs)) {
			bestVal = val
			bestPri = pri
			bestTs = it.rec.FetchedAt
			found = true
		}
	}

	return bestVal, found
}

func mergeItems(items []item) *models.SatelliteMetadata {
	out := &models.SatelliteMetadata{
		NoradID:  items[0].meta.NoradID,
		CosparID: items[0].meta.CosparID,
	}

	if v, ok := pickBest(items, "name",
		func(m *models.SatelliteMetadataPartial) (string, bool) {
			if m.Name != nil && *m.Name != "" {
				return *m.Name, true
			}
			return "", false
		},
	); ok {
		out.Name = v
	}

	if v, ok := pickBest(items, "operator",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.Operator != nil {
				return m.Operator, true
			}
			return nil, false
		},
	); ok {
		out.Operator = v
	}

	if v, ok := pickBest(items, "orbit_regime",
		func(m *models.SatelliteMetadataPartial) (models.OrbitRegime, bool) {
			if m.OrbitRegime != nil {
				return *m.OrbitRegime, true
			}
			return models.OrbitRegimeUnspecified, false
		},
	); ok {
		out.OrbitRegime = v
	}

	if v, ok := pickBest(items, "object_type",
		func(m *models.SatelliteMetadataPartial) (models.ObjectType, bool) {
			if m.ObjectType != nil {
				return *m.ObjectType, true
			}
			return models.ObjectTypeUnspecified, false
		},
	); ok {
		out.ObjectType = v
	}

	if v, ok := pickBest(items, "operational_status",
		func(m *models.SatelliteMetadataPartial) (models.OperationalStatus, bool) {
			if m.OperationalStatus != nil {
				return *m.OperationalStatus, true
			}
			return models.OperationalStatusUnspecified, false
		},
	); ok {
		out.OperationalStatus = v
	}

	sortByFreshness(items)

	for _, it := range items {
		m := it.meta

		if out.LaunchDate == nil && m.LaunchDate != nil {
			out.LaunchDate = m.LaunchDate
		}

		if out.LaunchSite == nil && m.LaunchSite != nil {
			out.LaunchSite = m.LaunchSite
		}

		if out.LaunchVehicle == nil && m.LaunchVehicle != nil {
			out.LaunchVehicle = m.LaunchVehicle
		}

		if out.Owner == nil && m.Owner != nil {
			out.Owner = m.Owner
		}

		if out.Constellation == nil && m.Constellation != nil {
			out.Constellation = m.Constellation
		}
	}

	for _, it := range items {
		out.Aliases = uniqueAppend(out.Aliases, it.meta.Aliases)
		out.Frequencies = mergeFrequencies(out.Frequencies, it.meta.Frequencies)

		out.Sources = append(out.Sources, models.SourceAttribution{
			Source:    it.rec.Source,
			FetchedAt: it.rec.FetchedAt,
		})
	}

	out.UpdatedAt = time.Now()

	return out
}

func sortByFreshness(items []item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].rec.FetchedAt.After(items[j].rec.FetchedAt)
	})
}

func uniqueAppend(dst, src []string) []string {
	set := make(map[string]struct{}, len(dst))

	for _, v := range dst {
		set[v] = struct{}{}
	}

	for _, v := range src {
		if _, ok := set[v]; !ok {
			dst = append(dst, v)
		}
	}

	return dst
}

func mergeFrequencies(dst, src []models.Frequency) []models.Frequency {
	return append(dst, src...)
}
