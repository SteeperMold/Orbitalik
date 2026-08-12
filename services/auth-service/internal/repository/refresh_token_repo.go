package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/jackc/pgx/v5"
)

type RefreshTokenRepository struct {
	db db.Conn
	sb sq.StatementBuilderType
}

func NewRefreshTokenRepository(db db.Conn) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token *models.StoredRefreshToken) error {
	query, args, err := r.sb.
		Insert("refresh_tokens").
		Columns("user_id", "hash", "expires_at").
		Values(token.UserID, token.Hash, token.ExpiresAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *RefreshTokenRepository) GetByHash(ctx context.Context, hash string) (*models.StoredRefreshToken, error) {
	query, args, err := r.sb.
		Select(
			"id",
			"user_id",
			"hash",
			"expires_at",
			"created_at",
			"revoked_at",
		).
		From("refresh_tokens").
		Where(sq.Eq{"hash": hash}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var token models.StoredRefreshToken

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&token.ID,
		&token.UserID,
		&token.Hash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTokenInvalid
		}
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id int) error {
	query, args, err := r.sb.
		Update("refresh_tokens").
		Set("revoked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).
		Where("revoked_at IS NULL").
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return domain.ErrTokenInvalid
	}

	return nil
}

func (r *RefreshTokenRepository) Rotate(ctx context.Context, oldID int, newToken *models.StoredRefreshToken) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx db.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	insertQuery, insertArgs, err := r.sb.
		Insert("refresh_tokens").
		Columns("user_id", "hash", "expires_at").
		Values(newToken.UserID, newToken.Hash, newToken.ExpiresAt).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, insertQuery, insertArgs...)
	if err != nil {
		return err
	}

	updateQuery, updateArgs, err := r.sb.
		Update("refresh_tokens").
		Set("revoked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": oldID}).
		Where("revoked_at IS NULL").
		ToSql()
	if err != nil {
		return err
	}

	res, err := tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return domain.ErrTokenInvalid
	}

	return tx.Commit(ctx)
}
