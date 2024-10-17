package migrate

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// MigratePg applies all pending postgres migrations.
func MigratePg(pool *pgxpool.Pool, path string) error {
	db := stdlib.OpenDBFromPool(pool)

	return Migrate(db, Postgres, path)
}

// Down rolls back a single postgres migration from the current version.
func DownPg(pool *pgxpool.Pool, path string) error {
	db := stdlib.OpenDBFromPool(pool)

	return Down(db, Postgres, path)
}

// DownTo rolls back postgres migrations to a specific version.
func DownPgTo(pool *pgxpool.Pool, path string, version int64) error {
	db := stdlib.OpenDBFromPool(pool)

	return DownTo(db, Postgres, path, version)
}
