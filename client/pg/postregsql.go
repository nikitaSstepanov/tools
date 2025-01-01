package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is interface for communicate with postgres database.
type Client interface {
	Pool

	RegisterTypes(types []string) error

	ToPgx() (*pgxpool.Pool, error)
}

// Config is type for database connection.
type Config struct {
	Host     string `yaml:"host"     env:"PG_HOST"    env-default:"localhost"`
	Port     int    `yaml:"port"     env:"PG_PORT"    env-default:"5432"`
	DBName   string `yaml:"dbname"   env:"PG_NAME"    env-default:"postgres"`
	Username string `yaml:"username" env:"PG_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	SSLMode  string `yaml:"sslmode"  env:"PG_SSLMODE" env-default:"disabled"`
}

type pgclient struct {
	afterConnectFuncs []func(ctx context.Context, conn *Conn) error
	Pool
}

// ConnectToDb returns a pointer to a pgxpool.Pool representing the database connection pool.
func New(ctx context.Context, cfg *Config) (Client, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	))
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

	client := &pgclient{
		Pool:              db,
		afterConnectFuncs: make([]func(ctx context.Context, conn *pgx.Conn) error, 0),
	}

	return client, nil
}

func NewWithPool(ctx context.Context, pool Pool) (Client, error) {
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	client := &pgclient{
		Pool:              pool,
		afterConnectFuncs: make([]func(ctx context.Context, conn *Conn) error, 0),
	}

	return client, nil
}

func (pc *pgclient) RegisterTypes(types []string) error {
	function := func(ctx context.Context, conn *pgx.Conn) error {
		for _, typeName := range types {
			t, err := conn.LoadType(ctx, typeName)
			if err != nil {
				return err
			}

			conn.TypeMap().RegisterType(t)
		}

		return nil
	}

	pc.afterConnectFuncs = append(pc.afterConnectFuncs, function)

	cfg := pc.Pool.Config()

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, f := range pc.afterConnectFuncs {
			if err := f(ctx, conn); err != nil {
				return err
			}
		}

		return nil
	}

	ctx := context.Background()

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	if err := db.Ping(ctx); err != nil {
		return err
	}

	pc.Pool = db

	return nil
}

func (pc *pgclient) ToPgx() (*pgxpool.Pool, error) {
	cfg := pc.Pool.Config()
	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
