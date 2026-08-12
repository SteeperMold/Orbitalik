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

type CredentialsRepository struct {
	db db.Conn
	sb sq.StatementBuilderType
}

func NewCredentialsRepository(db db.Conn) *CredentialsRepository {
	return &CredentialsRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *CredentialsRepository) Create(ctx context.Context, cred *models.Credentials) error {
	query, args, err := r.sb.
		Insert("credentials").
		Columns("user_id", "password_hash").
		Values(cred.UserID, cred.PasswordHash).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *CredentialsRepository) GetByUserID(ctx context.Context, userID int) (*models.Credentials, error) {
	query, args, err := r.sb.
		Select(
			"user_id",
			"password_hash",
			"created_at",
			"updated_at",
		).
		From("credentials").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var cred models.Credentials

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&cred.UserID,
		&cred.PasswordHash,
		&cred.CreatedAt,
		&cred.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &cred, nil
}

func (r *CredentialsRepository) UpdatePasswordHash(ctx context.Context, userID int, hash string) error {
	query, args, err := r.sb.
		Update("credentials").
		Set("password_hash", hash).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
