package repository

import (
	"context"
	"errors"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/user-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/user-service/internal/models"
)

type UserRepository struct {
	db     db.Conn
	sb     sq.StatementBuilderType
	logger applog.Logger
}

func NewUserRepository(db db.Conn, logger applog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		sb:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, params *domain.CreateUserParams) (*models.User, error) {
	query, args, err := r.sb.
		Insert("users").
		Columns("email", "username", "password_hash").
		Values(params.Email, params.Username, params.PasswordHash).
		Suffix("RETURNING id, email, username, password_hash, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)

	var u models.User

	err = row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query, args, err := r.sb.
		Select("id", "email", "username", "password_hash", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)

	var u models.User

	err = row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := r.sb.
		Select("id", "email", "username", "password_hash", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()

	row := r.db.QueryRow(ctx, query, args...)

	var u models.User

	err = row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, params *domain.UpdateUserParams) (*models.User, error) {
	builder := r.sb.Update("users")

	if params.Email != "" {
		builder = builder.Set("email", params.Email)
	}

	if params.Username != "" {
		builder = builder.Set("username", params.Username)
	}

	if params.PasswordHash != "" {
		builder = builder.Set("password_hash", params.PasswordHash)
	}

	query, args, err := builder.
		Where(sq.Eq{"id": params.ID}).
		Suffix("RETURNING id, email, username, password_hash, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)

	var u models.User

	err = row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query, args, err := r.sb.
		Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	return nil
}
