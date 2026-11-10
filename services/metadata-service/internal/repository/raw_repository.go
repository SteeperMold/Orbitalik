package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/jackc/pgx/v5"
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

func (r *RawMetadataRepository) SaveRawBatch(ctx context.Context, data []*models.SatelliteIngestRecord) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
	}()

	pgxTx, ok := tx.(postgres.CopyFromTx)
	if !ok {
		return errors.New("CopyFrom not supported")
	}

	_, err = pgxTx.CopyFrom(
		ctx,
		pgx.Identifier{"satellite_metadata_raw"},
		[]string{"norad_id", "cospar_id", "source", "payload", "fetched_at"},
		newRawMetadataCopySource(data),
	)

	return err
}
