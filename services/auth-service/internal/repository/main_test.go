//go:build integration
// +build integration

package repository_test

import (
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"os"
	"testing"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	commonPgx "github.com/SteeperMold/Orbitalik/common/go/db/postgres"
	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testPool *pgxpool.Pool
	testDB   *sql.DB
	testConn db.Conn
)

// testUUID mirrors the deterministic legacy BIGINT -> UUIDv7 mapping used by
// the database migrations.
func testUUID(n uint64) uuid.UUID {
	var raw [16]byte
	binary.BigEndian.PutUint64(raw[7:15], n)
	raw[6] = 0x70
	raw[8] = (raw[8] & 0x3f) | 0x80

	return uuid.UUID(raw)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = pgContainer.Terminate(ctx)
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer testPool.Close()

	testDB, err = sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = testDB.Close()
	}()

	runMigrations(testDB)

	testConn = commonPgx.NewConn(testPool)

	code := m.Run()
	os.Exit(code)
}

func runMigrations(db *sql.DB) {
	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../../../../db/migrations/auth", "postgres", driver)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}
