package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host              string
	Port              string
	Name              string
	User              string
	Password          string
	ConnectionTimeout time.Duration
}

func OpenConn(ctx context.Context, cfg *Config) (db.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
	defer cancel()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres conn: %w", err)
	}

	return &DB{pool: pool}, nil
}
