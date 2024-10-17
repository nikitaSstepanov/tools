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
	// Begin starts a pseudo nested transaction implemented with a savepoint.
	Begin(ctx context.Context) (pgx.Tx, error)

	// BeginTx acquires a connection from the Pool and starts a transaction with pgx.TxOptions determining the transaction mode.
	// Unlike database/sql, the context only affects the begin command. i.e. there is no auto-rollback on context cancellation.
	// *pgxpool.Tx is returned, which implements the pgx.Tx interface.
	// Commit or Rollback must be called on the returned transaction to finalize the transaction block.
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)

	// Close closes all connections in the pool and rejects future Acquire calls. Blocks until all connections are returned
	// to pool and closed.
	Close()

	// Exec acquires a connection from the Pool and executes the given SQL.
	// SQL can be either a prepared statement name or an SQL string.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// The acquired connection is returned to the pool when the Exec function returns.
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)

	// Ping acquires a connection from the Pool and executes an empty sql statement against it.
	// If the sql returns without error, the database Ping is considered successful, otherwise, the error is returned.
	Ping(ctx context.Context) error

	// Query acquires a connection and executes a query that returns pgx.Rows.
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	// See pgx.Rows documentation to close the returned Rows and return the acquired connection to the Pool.
	//
	// If there is an error, the returned pgx.Rows will be returned in an error state.
	// If preferred, ignore the error returned from Query and handle errors using the returned pgx.Rows.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	// QueryRow acquires a connection and executes a query that is expected
	// to return at most one row (pgx.Row). Errors are deferred until pgx.Row's
	// Scan method is called. If the query selects no rows, pgx.Row's Scan will
	// return ErrNoRows. Otherwise, pgx.Row's Scan scans the first selected row
	// and discards the rest. The acquired connection is returned to the Pool when
	// pgx.Row's Scan method is called.
	//
	// Arguments should be referenced positionally from the SQL string as $1, $2, etc.
	//
	// For extra control over how the query is executed, the types QuerySimpleProtocol, QueryResultFormats, and
	// QueryResultFormatsByOID may be used as the first args to control exactly how the query is executed. This is rarely
	// needed. See the documentation for those types for details.
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// Reset closes all connections, but leaves the pool open. It is intended for use when an error is detected that would
	// disrupt all connections (such as a network interruption or a server state change).
	//
	// It is safe to reset a pool while connections are checked out. Those connections will be closed when they are returned
	// to the pool.
	Reset()
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

// ConnectToDb returns a pointer to a pgxpool.Pool representing the database connection pool.
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
