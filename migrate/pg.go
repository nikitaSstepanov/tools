package migrate

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func MigratePg(pool *pgxpool.Pool, path string) error {
	db := stdlib.OpenDBFromPool(pool)

	return Migrate(db, Postgres, path)
}

func DownPg(pool *pgxpool.Pool, path string) error {
	db := stdlib.OpenDBFromPool(pool)

	return Down(db, Postgres, path)
}

func DownPgTo(pool *pgxpool.Pool, path string, version int64) error {
	db := stdlib.OpenDBFromPool(pool)

	return DownTo(db, Postgres, path, version)
}
