package repository

import (
	"context"
	"errors"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type TLERepository struct {
	db domain.DBConn
}

func NewTLERepository(db domain.DBConn) *TLERepository {
	return &TLERepository{
		db: db,
	}
}

func (r *TLERepository) SaveBatch(ctx context.Context, tles []*models.TLE) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	const tempTableQ = `
		CREATE TEMP TABLE tle_stage
		(
			norad_id INT,
			satellite_name TEXT,
			line1 TEXT,
			line2 TEXT,
			epoch TIMESTAMPTZ
		) ON COMMIT DROP 
	`

	_, err = tx.Exec(ctx, tempTableQ)
	if err != nil {
		return err
	}

	rows := make([][]any, len(tles))
	for i, tle := range tles {
		rows[i] = []any{
			tle.NoradID,
			tle.SatelliteName,
			tle.Line1,
			tle.Line2,
			tle.Epoch,
		}
	}

	_, err = tx.CopyFrom(ctx,
		pgx.Identifier{"tle_stage"},
		[]string{"norad_id", "satellite_name", "line1", "line2", "epoch"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	const insertQ = `
		INSERT INTO tle(norad_id, satellite_name, line1, line2, epoch)
		SELECT norad_id, satellite_name, line1, line2, epoch
		FROM tle_stage
		ON CONFLICT (norad_id, epoch) DO NOTHING 
	`

	_, err = tx.Exec(ctx, insertQ)
	if err != nil {
		return err
	}

	return nil
}

func (r *TLERepository) GetAllTLEs(ctx context.Context) ([]*models.TLE, error) {
	const q = `
		SELECT id, norad_id, satellite_name, line1, line2, epoch, fetched_at
		FROM tle
		ORDER BY norad_id
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tles []*models.TLE
	for rows.Next() {
		var tle models.TLE

		err := rows.Scan(&tle.ID, &tle.NoradID, &tle.SatelliteName, &tle.Line1, &tle.Line2, &tle.Epoch, &tle.FetchedAt)
		if err != nil {
			return nil, err
		}

		tles = append(tles, &tle)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tles, nil
}

func (r *TLERepository) GetTLEByNoradID(ctx context.Context, noradID int) (*models.TLE, error) {
	const q = `
		SELECT id, norad_id, satellite_name, line1, line2, epoch, fetched_at
		FROM tle
		WHERE norad_id = $1
		ORDER BY epoch DESC 
		LIMIT 1
	`

	row := r.db.QueryRow(ctx, q, noradID)

	var tle models.TLE
	err := row.Scan(&tle.ID, &tle.NoradID, &tle.SatelliteName, &tle.Line1, &tle.Line2, &tle.Epoch, &tle.FetchedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tle, nil
}

func (r *TLERepository) GetTLEBySatelliteName(ctx context.Context, name string) (*models.TLE, error) {
	const q = `
		SELECT id, norad_id, satellite_name, line1, line2, epoch, fetched_at
		FROM tle
		WHERE satellite_name = $1
		ORDER BY epoch DESC 
		LIMIT 1
	`

	row := r.db.QueryRow(ctx, q, name)

	var tle models.TLE
	err := row.Scan(&tle.ID, &tle.NoradID, &tle.SatelliteName, &tle.Line1, &tle.Line2, &tle.Epoch, &tle.FetchedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tle, nil
}
