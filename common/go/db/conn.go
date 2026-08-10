package db

import "context"

// Result is the return type of Exec queries
type Result interface {
	RowsAffected() int64
}

// Rows abstracts a result set returned by Query
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// Row abstracts a single-row result from QueryRow
type Row interface {
	Scan(dest ...any) error
}

// Tx abstracts a transaction
type Tx interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Conn abstracts a database connection (pool or single connection)
type Conn interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Begin(ctx context.Context) (Tx, error)
	Ping(ctx context.Context) error
}
