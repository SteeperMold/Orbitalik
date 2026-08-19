package aggregation

import (
	"time"

	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
)

func appendFieldSource(
	out *models.SatelliteMetadata,
	field string,
	it *item,
) {
	out.Sources = append(out.Sources, models.FieldSource{
		Field:     field,
		Sources:   []models.Source{it.rec.Source},
		FetchedAt: it.rec.FetchedAt,
	})
}

func mergeName(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "name",
		func(m *models.SatelliteMetadataPartial) (string, bool) {
			if m.Name == nil || *m.Name == "" {
				return "", false
			}

			return *m.Name, true
		},
	); ok {
		out.Name = v
		appendFieldSource(out, "name", it)
	}
}

func mergeCosparID(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "cospar_id",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.CosparID == nil || *m.CosparID == "" {
				return nil, false
			}

			return m.CosparID, true
		},
	); ok {
		out.CosparID = v
		appendFieldSource(out, "cospar_id", it)
	}
}

func mergeOperator(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "operator",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.Operator == nil {
				return nil, false
			}

			return m.Operator, true
		},
	); ok {
		out.Operator = v
		appendFieldSource(out, "operator", it)
	}
}

func mergeOrbitRegime(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "orbit_regime",
		func(m *models.SatelliteMetadataPartial) (models.OrbitRegime, bool) {
			if m.OrbitRegime == nil {
				return models.OrbitRegimeUnspecified, false
			}

			return *m.OrbitRegime, true
		},
	); ok {
		out.OrbitRegime = v
		appendFieldSource(out, "orbit_regime", it)
	}
}

func mergeObjectType(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "object_type",
		func(m *models.SatelliteMetadataPartial) (models.ObjectType, bool) {
			if m.ObjectType == nil {
				return models.ObjectTypeUnspecified, false
			}

			return *m.ObjectType, true
		},
	); ok {
		out.ObjectType = v
		appendFieldSource(out, "object_type", it)
	}
}

func mergeOperationalStatus(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "operational_status",
		func(m *models.SatelliteMetadataPartial) (models.OperationalStatus, bool) {
			if m.OperationalStatus == nil {
				return models.OperationalStatusUnspecified, false
			}

			return *m.OperationalStatus, true
		},
	); ok {
		out.OperationalStatus = v
		appendFieldSource(out, "operational_status", it)
	}
}

func mergeLaunchDate(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "launch_date",
		func(m *models.SatelliteMetadataPartial) (*time.Time, bool) {
			if m.LaunchDate == nil {
				return nil, false
			}

			return m.LaunchDate, true
		},
	); ok {
		out.LaunchDate = v
		appendFieldSource(out, "launch_date", it)
	}
}

func mergeLaunchSite(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "launch_site",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.LaunchSite == nil {
				return nil, false
			}

			return m.LaunchSite, true
		},
	); ok {
		out.LaunchSite = v
		appendFieldSource(out, "launch_site", it)
	}
}

func mergeLaunchVehicle(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "launch_vehicle",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.LaunchVehicle == nil {
				return nil, false
			}

			return m.LaunchVehicle, true
		},
	); ok {
		out.LaunchVehicle = v
		appendFieldSource(out, "launch_vehicle", it)
	}
}

func mergeOwner(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "owner",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.Owner == nil {
				return nil, false
			}

			return m.Owner, true
		},
	); ok {
		out.Owner = v
		appendFieldSource(out, "owner", it)
	}
}

func mergeConstellation(out *models.SatelliteMetadata, items []item) {
	if v, it, ok := pickBest(items, "constellation",
		func(m *models.SatelliteMetadataPartial) (*string, bool) {
			if m.Constellation == nil {
				return nil, false
			}

			return m.Constellation, true
		},
	); ok {
		out.Constellation = v
		appendFieldSource(out, "constellation", it)
	}
}
