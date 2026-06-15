package repository

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/jackc/pgx/v5"
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

func (r *MetadataRepository) GetMetadataByNoradID(ctx context.Context, noradID int) (*models.SatelliteMetadata, error) {
	q := r.psql.
		Select(metadataSelectColumns()...).
		From("satellite_metadata").
		Where(sq.Eq{"norad_id": noradID})

	sqlQ, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlQ, args...)
	return scanSatelliteMetadataRow(row)
}

func (r *MetadataRepository) GetMetadataByName(ctx context.Context, name string) (*models.SatelliteMetadata, error) {
	q := r.psql.
		Select(metadataSelectColumns()...).
		From("satellite_metadata").
		Where(sq.Eq{"name": name})

	sqlQ, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlQ, args...)
	return scanSatelliteMetadataRow(row)
}

func (r *MetadataRepository) ListSatellites(
	ctx context.Context,
	filter *models.ListFilter,
) (items []*models.SatelliteMetadata, nextPageToken string, total uint32, err error) {

	if filter == nil {
		filter = &models.ListFilter{}
	}

	offset, err := parsePageToken(filter.PageToken)
	if err != nil {
		return nil, "", 0, fmt.Errorf("invalid page token: %w", err)
	}

	where := buildListWhere(filter)

	total64, err := r.countList(ctx, where)
	if err != nil {
		return nil, "", 0, err
	}
	total = uint32(total64)

	q := r.psql.
		Select(metadataSelectColumns()...).
		From("satellite_metadata").
		OrderBy("norad_id ASC").
		Limit(uint64(filter.PageSize)).
		Offset(offset)

	for _, w := range where {
		q = q.Where(w)
	}

	sqlQ, args, err := q.ToSql()
	if err != nil {
		return nil, "", 0, err
	}

	rows, err := r.db.Query(ctx, sqlQ, args...)
	if err != nil {
		return nil, "", 0, err
	}
	defer rows.Close()

	items = make([]*models.SatelliteMetadata, 0, filter.PageSize)
	for rows.Next() {
		item, err := scanSatelliteMetadataRow(rows)
		if err != nil {
			return nil, "", 0, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, "", 0, err
	}

	nextOffset := offset + uint64(len(items))
	nextPageToken = ""
	if nextOffset < uint64(total) {
		nextPageToken = strconv.FormatUint(nextOffset, 10)
	}

	return items, nextPageToken, total, nil
}

func (r *MetadataRepository) countList(ctx context.Context, where []sq.Sqlizer) (uint64, error) {
	q := r.psql.
		Select("COUNT(*)").
		From("satellite_metadata")

	for _, w := range where {
		q = q.Where(w)
	}

	sqlQ, args, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	var total uint64
	if err := r.db.QueryRow(ctx, sqlQ, args...).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func buildListWhere(filter *models.ListFilter) []sq.Sqlizer {
	where := make([]sq.Sqlizer, 0, 5)

	if filter.ObjectType != nil {
		where = append(where, sq.Eq{"object_type": *filter.ObjectType})
	}
	if filter.MissionType != nil {
		where = append(where, sq.Eq{"mission_type": *filter.MissionType})
	}
	if filter.OperationalStatus != nil {
		where = append(where, sq.Eq{"operational_status": *filter.OperationalStatus})
	}
	if filter.OrbitRegime != nil {
		where = append(where, sq.Eq{"orbit_regime": *filter.OrbitRegime})
	}
	if filter.Constellation != nil && *filter.Constellation != "" {
		where = append(where, sq.Eq{"constellation": *filter.Constellation})
	}

	return where
}

func parsePageToken(token string) (uint64, error) {
	if token == "" {
		return 0, nil
	}

	n, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func scanSatelliteMetadataRow(row pgx.Row) (*models.SatelliteMetadata, error) {
	var m models.SatelliteMetadata

	var (
		cosparID      *string
		operator      *string
		owner         *string
		constellation *string
		launchSite    *string
		launchVehicle *string
	)

	err := row.Scan(
		&m.NoradID,
		&cosparID,
		&m.Name,
		&m.Aliases,
		&m.ObjectType,
		&m.MissionType,
		&m.OrbitRegime,
		&operator,
		&owner,
		&constellation,
		&m.LaunchDate,
		&launchSite,
		&launchVehicle,
		&m.OperationalStatus,
		&m.Frequencies,
		&m.Sources,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	m.CosparID = cosparID
	m.Operator = operator
	m.Owner = owner
	m.Constellation = constellation
	m.LaunchSite = launchSite
	m.LaunchVehicle = launchVehicle

	return &m, nil
}

func metadataSelectColumns() []string {
	return []string{
		"norad_id",
		"cospar_id",
		"name",
		"aliases",
		"object_type",
		"mission_type",
		"orbit_regime",
		"operator",
		"owner",
		"constellation",
		"launch_date",
		"launch_site",
		"launch_vehicle",
		"operational_status",
		"frequencies",
		"sources",
		"updated_at",
	}
}
