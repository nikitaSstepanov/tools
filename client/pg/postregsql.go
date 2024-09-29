package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is interface for communicate with postgres database.
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Close()
}

// Config is type for database connection
type Config struct {
	Host     string `yaml:"host"     env:"PG_HOST"    env-default:"localhost"`
	Port     int    `yaml:"port"     env:"PG_PORT"    env-default:"5432"`
	DBName   string `yaml:"dbname"   env:"PG_NAME"    env-default:"postgres"`
	Username string `yaml:"username" env:"PG_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	SSLMode  string `yaml:"sslmode"  env:"PG_SSLMODE" env-default:"disabled"`
}

func getConfig(cfg *Config) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	))
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ConnectToDb creates connection
func ConnectToDb(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	config, err := getConfig(cfg)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
