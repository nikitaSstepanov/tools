package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Row            = pgx.Row
	Rows           = pgx.Rows
	Conn           = pgx.Conn
	Tx             = pgx.Tx
	TxOptions      = pgx.TxOptions
	Stat           = pgxpool.Stat
	Batch          = pgx.Batch
	BatchResults   = pgx.BatchResults
	Identifier     = pgx.Identifier
	CopyFromSource = pgx.CopyFromSource
)

var (
	ErrTxCommitRollback = pgx.ErrTxCommitRollback
	ErrTooManyRows      = pgx.ErrTooManyRows
	ErrTxClosed         = pgx.ErrTxClosed
	ErrNoRows           = pgx.ErrNoRows
)

type Pool interface {
	Acquire(ctx context.Context) (c *pgxpool.Conn, err error)
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
	Begin(ctx context.Context) (pgx.Tx, error)
	Stat() *pgxpool.Stat
	Config() *pgxpool.Config
	Reset()
	Close()
}
