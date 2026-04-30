package repository

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type TLERepository struct {
	db   db.Conn
	psql sq.StatementBuilderType
}

func NewTLERepository(db db.Conn) *TLERepository {
	return &TLERepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TLERepository) SaveBatch(ctx context.Context, tles []*models.TLE) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	pgxTx, ok := tx.(postgres.CopyFromTx)
	if !ok {
		return errors.New("CopyFrom not supported")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
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
		) ON COMMIT DROP;
		CREATE INDEX ON tle_stage (norad_id, epoch);
	`

	_, err = pgxTx.Exec(ctx, tempTableQ)
	if err != nil {
		return err
	}

	_, err = pgxTx.CopyFrom(
		ctx,
		pgx.Identifier{"tle_stage"},
		[]string{"norad_id", "satellite_name", "line1", "line2", "epoch"},
		newTLECopySource(tles),
	)
	if err != nil {
		return err
	}

	insertQ, args, err := r.psql.
		Insert("tle").
		Columns("norad_id", "satellite_name", "line1", "line2", "epoch").
		Select(
			r.psql.
				Select("norad_id", "satellite_name", "line1", "line2", "epoch").
				From("tle_stage"),
		).
		Suffix("ON CONFLICT (norad_id, epoch) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	_, err = pgxTx.Exec(ctx, insertQ, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TLERepository) DeleteOlderThan(ctx context.Context, d time.Duration) error {
	q, args, err := r.psql.
		Delete("tle").
		Where(sq.Expr("epoch < NOW() - ?::interval", d.String())).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, q, args...)
	return err
}

func (r *TLERepository) GetAllTLEs(ctx context.Context) ([]*models.TLE, error) {
	q, args, err := r.psql.
		Select("id", "norad_id", "satellite_name", "line1", "line2", "epoch", "fetched_at").
		From("tle").
		OrderBy("norad_id").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, q, args...)
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
	q, args, err := r.psql.
		Select("id", "norad_id", "satellite_name", "line1", "line2", "epoch", "fetched_at").
		From("tle").
		Where(sq.Eq{"norad_id": noradID}).
		OrderBy("epoch DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, q, args...)

	var tle models.TLE
	err = row.Scan(&tle.ID, &tle.NoradID, &tle.SatelliteName, &tle.Line1, &tle.Line2, &tle.Epoch, &tle.FetchedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tle, nil
}

func (r *TLERepository) GetTLEBySatelliteName(ctx context.Context, name string) (*models.TLE, error) {
	q, args, err := r.psql.
		Select("id", "norad_id", "satellite_name", "line1", "line2", "epoch", "fetched_at").
		From("tle").
		Where(sq.Eq{"satellite_name": name}).
		OrderBy("epoch DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, q, args...)

	var tle models.TLE
	err = row.Scan(&tle.ID, &tle.NoradID, &tle.SatelliteName, &tle.Line1, &tle.Line2, &tle.Epoch, &tle.FetchedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &tle, nil
}
