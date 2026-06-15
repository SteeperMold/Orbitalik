package repository

import (
	"context"
	"encoding/json"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
)

type MetadataRepository struct {
	db   db.Conn
	psql sq.StatementBuilderType
}

func NewMetadataRepository(db db.Conn) *MetadataRepository {
	return &MetadataRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *MetadataRepository) Upsert(ctx context.Context, meta *models.SatelliteMetadata) error {
	now := time.Now()

	if meta.Aliases == nil {
		meta.Aliases = []string{}
	}
	if meta.Frequencies == nil {
		meta.Frequencies = []models.Frequency{}
	}
	if meta.Sources == nil {
		meta.Sources = []models.SourceAttribution{}
	}

	freqJSON, err := json.Marshal(meta.Frequencies)
	if err != nil {
		return err
	}

	srcJSON, err := json.Marshal(meta.Sources)
	if err != nil {
		return err
	}

	query := r.psql.
		Insert("satellite_metadata").
		Columns(
			"norad_id", "cospar_id", "name", "aliases", "object_type", "mission_type", "orbit_regime",
			"operator", "owner", "constellation", "launch_date", "launch_site", "launch_vehicle", "operational_status",
			"frequencies", "sources", "updated_at",
		).
		Values(
			meta.NoradID, meta.CosparID, meta.Name, meta.Aliases, meta.ObjectType, meta.MissionType,
			meta.OrbitRegime, meta.Operator, meta.Owner, meta.Constellation, meta.LaunchDate, meta.LaunchSite,
			meta.LaunchVehicle, meta.OperationalStatus, freqJSON, srcJSON, now,
		).
		Suffix(`
			ON CONFLICT (norad_id)
			DO UPDATE SET
				cospar_id = EXCLUDED.cospar_id,
				name = EXCLUDED.name,
				aliases = EXCLUDED.aliases,
				object_type = EXCLUDED.object_type,
				mission_type = EXCLUDED.mission_type,
				orbit_regime = EXCLUDED.orbit_regime,
				operator = EXCLUDED.operator,
				owner = EXCLUDED.owner,
				constellation = EXCLUDED.constellation,
				launch_date = EXCLUDED.launch_date,
				launch_site = EXCLUDED.launch_site,
				launch_vehicle = EXCLUDED.launch_vehicle,
				operational_status = EXCLUDED.operational_status,
				frequencies = EXCLUDED.frequencies,
				sources = EXCLUDED.sources,
				updated_at = EXCLUDED.updated_at
		`)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
