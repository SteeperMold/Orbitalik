package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/metadata-aggregation-worker/internal/models"
)

type RawMetadataRepository struct {
	db   db.Conn
	psql sq.StatementBuilderType
}

func NewRawMetadataRepository(db db.Conn) *RawMetadataRepository {
	return &RawMetadataRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RawMetadataRepository) GetByNoradID(ctx context.Context, noradID int) ([]models.SatelliteIngestRecord, error) {
	query := r.psql.
		Select(
			"id",
			"norad_id",
			"cospar_id",
			"source",
			"payload",
			"fetched_at",
			"stored_at",
		).
		From("satellite_metadata_raw").
		Where(sq.Eq{"norad_id": noradID}).
		OrderBy("fetched_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows db.Rows) {
		_ = rows.Close()
	}(rows)

	var result []models.SatelliteIngestRecord

	for rows.Next() {
		var rcd models.SatelliteIngestRecord
		var cosparID *string

		var fetchedAt, storedAt time.Time

		err := rows.Scan(
			&rcd.ID,
			&rcd.NoradID,
			&cosparID,
			&rcd.Source,
			&rcd.Payload,
			&fetchedAt,
			&storedAt,
		)
		if err != nil {
			return nil, err
		}

		rcd.CosparID = cosparID
		rcd.FetchedAt = fetchedAt
		rcd.StoredAt = storedAt

		result = append(result, rcd)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
