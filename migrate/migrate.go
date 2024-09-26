package migrate

import (
	"database/sql"

	"github.com/pressly/goose"
)

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
