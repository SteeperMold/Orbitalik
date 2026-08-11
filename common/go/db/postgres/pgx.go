package postgres

import (
	"context"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CopyFromTx interface {
	db.Tx
	CopyFrom(ctx context.Context, table pgx.Identifier, columns []string, source pgx.CopyFromSource) (int64, error)
}

type DB struct {
	pool *pgxpool.Pool
}

func (p *DB) Query(ctx context.Context, sql string, args ...any) (db.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows: rows}, nil
}

func (p *DB) QueryRow(ctx context.Context, sql string, args ...any) db.Row {
	row := p.pool.QueryRow(ctx, sql, args...)
	return &Row{row: row}
}

func (p *DB) Exec(ctx context.Context, sql string, args ...any) (db.Result, error) {
	return p.pool.Exec(ctx, sql, args...)
}

func (p *DB) Begin(ctx context.Context) (db.Tx, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

func (p *DB) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *DB) Close() {
	p.pool.Close()
}

type Rows struct {
	rows pgx.Rows
}

func (r *Rows) Next() bool {
	return r.rows.Next()
}

func (r *Rows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *Rows) Close() error {
	r.rows.Close()
	return nil
}

func (r *Rows) Err() error {
	return r.rows.Err()
}

type Row struct {
	row pgx.Row
}

func (r *Row) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

type Tx struct {
	tx pgx.Tx
}

func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *Tx) Query(ctx context.Context, sql string, args ...any) (db.Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{rows: rows}, nil
}

func (t *Tx) QueryRow(ctx context.Context, sql string, args ...any) db.Row {
	return &Row{row: t.tx.QueryRow(ctx, sql, args...)}
}

func (t *Tx) Exec(ctx context.Context, sql string, args ...any) (db.Result, error) {
	return t.tx.Exec(ctx, sql, args...)
}

func (t *Tx) CopyFrom(
	ctx context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	return t.tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
}
