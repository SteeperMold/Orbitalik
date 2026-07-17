package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSQLDatabase(dbConfig *DBConfig) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.ConnectionTimeout)
	defer cancel()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	log.Printf("db is running on %v\n", connString)

	return conn
}

func CloseDBConnection(pool *pgxpool.Pool) {
	pool.Close()
}
