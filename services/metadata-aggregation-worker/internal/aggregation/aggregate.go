package aggregation

import (
	"encoding/json"
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
)

type item struct {
	meta *models.SatelliteMetadataPartial
	rec  models.SatelliteIngestRecord
}

func Aggregate(records []models.SatelliteIngestRecord) *models.SatelliteMetadata {
	if len(records) == 0 {
		return nil
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
		return nil
	}

	return mergeItems(items)
}

func decodePayload(r models.SatelliteIngestRecord) (*models.SatelliteMetadataPartial, error) {
	var m models.SatelliteMetadataPartial
	if err := json.Unmarshal(r.Payload, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// higher number means more priority
var fieldPriority = map[string]map[models.Source]int{
	"name": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"cospar_id": {
		models.SourceUCS:       1, // freshness wins
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
	"launch_date": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"launch_site": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"launch_vehicle": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"owner": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
	"constellation": {
		models.SourceUCS:       2,
		models.SourceCelestrak: 1,
	},
}

func pickBest[T any](
	items []item,
	field string,
	extract func(*models.SatelliteMetadataPartial) (T, bool),
) (bestVal T, bestItem *item, found bool) {

	var (
		bestPri = -1
		bestTs  time.Time
	)

	for i := range items {
		it := &items[i]

		val, ok := extract(it.meta)
		if !ok {
			continue
		}

		pri := fieldPriority[field][it.rec.Source]

		if !found || pri > bestPri || (pri == bestPri && it.rec.FetchedAt.After(bestTs)) {
			bestVal = val
			bestItem = it
			bestPri = pri
			bestTs = it.rec.FetchedAt
			found = true
		}
	}

	return bestVal, bestItem, found
}

func mergeItems(items []item) *models.SatelliteMetadata {
	out := &models.SatelliteMetadata{
		NoradID: items[0].meta.NoradID,
	}

	mergeName(out, items)
	mergeCosparID(out, items)
	mergeOperator(out, items)
	mergeOrbitRegime(out, items)
	mergeObjectType(out, items)
	mergeOperationalStatus(out, items)
	mergeLaunchDate(out, items)
	mergeLaunchSite(out, items)
	mergeLaunchVehicle(out, items)
	mergeOwner(out, items)
	mergeConstellation(out, items)

	for _, it := range items {
		out.Aliases = uniqueAppend(out.Aliases, it.meta.Aliases)
		out.Frequencies = mergeFrequencies(out.Frequencies, it.meta.Frequencies)
	}

	if sources := sourcesForField(items, func(m *models.SatelliteMetadataPartial) bool {
		return len(m.Aliases) > 0
	}); len(sources) > 0 {
		out.Sources = append(out.Sources, models.FieldSource{
			Field:   "aliases",
			Sources: sources,
		})
	}

	if sources := sourcesForField(items, func(m *models.SatelliteMetadataPartial) bool {
		return len(m.Frequencies) > 0
	}); len(sources) > 0 {
		out.Sources = append(out.Sources, models.FieldSource{
			Field:   "frequencies",
			Sources: sources,
		})
	}

	out.UpdatedAt = time.Now()

	return out
}

func uniqueAppend(dst, src []string) []string {
	set := make(map[string]struct{}, len(dst))

	for _, v := range dst {
		set[v] = struct{}{}
	}

	for _, v := range src {
		if _, ok := set[v]; ok {
			continue
		}

		dst = append(dst, v)
		set[v] = struct{}{}
	}

	return dst
}

// TODO: merge frequencies properly
func mergeFrequencies(dst, src []models.Frequency) []models.Frequency {
	return append(dst, src...)
}

func sourcesForField(items []item, hasValue func(*models.SatelliteMetadataPartial) bool) []models.Source {
	seen := make(map[models.Source]struct{})
	var sources []models.Source

	for _, it := range items {
		if !hasValue(it.meta) {
			continue
		}

		if _, ok := seen[it.rec.Source]; ok {
			continue
		}

		seen[it.rec.Source] = struct{}{}
		sources = append(sources, it.rec.Source)
	}

	return sources
}
