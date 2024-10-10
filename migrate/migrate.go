package migrate

import (
	"database/sql"

	"github.com/pressly/goose"
)

// Migrate applies all pending database migrations to the specified database.
func Migrate(db *sql.DB, dialect Dialect, path string) error {
	if err := db.Ping(); err != nil {
		return err
	}

	d := diaToStr(dialect)

	if err := goose.SetDialect(d); err != nil {
		return err
	}

	if err := goose.Up(db, path); err != nil {
		return err
	}

	return nil
}

// Down rolls back a single migration from the current version
func Down(db *sql.DB, dialect Dialect, path string) error {
	if err := db.Ping(); err != nil {
		return err
	}

	d := diaToStr(dialect)

	if err := goose.SetDialect(d); err != nil {
		return err
	}

	if err := goose.Down(db, path); err != nil {
		return err
	}

	return nil
}

// DownTo rolls back migrations to a specific version.
func DownTo(db *sql.DB, dialect Dialect, path string, version int64) error {
	if err := db.Ping(); err != nil {
		return err
	}

	d := diaToStr(dialect)

	if err := goose.SetDialect(d); err != nil {
		return err
	}

	if err := goose.DownTo(db, path, version); err != nil {
		return err
	}

	return nil
}
