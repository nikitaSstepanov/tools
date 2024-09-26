package pg

import (
	"github.com/jackc/pgx/v5"
)

type Row = pgx.Row

var (
	ErrTxCommitRollback = pgx.ErrTxCommitRollback
	ErrTooManyRows      = pgx.ErrTooManyRows
	ErrTxClosed         = pgx.ErrTxClosed
	ErrNoRows           = pgx.ErrNoRows
)
