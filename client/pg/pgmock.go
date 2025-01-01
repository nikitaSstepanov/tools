package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgmock struct {
	afterConnectFuncs []func(ctx context.Context, conn *Conn) error
	Pool
}

func NewWithMock(pool Pool) Client {
	return &pgmock{
		Pool:              pool,
		afterConnectFuncs: make([]func(ctx context.Context, conn *Conn) error, 0),
	}
}

func (pc *pgmock) RegisterTypes(types []string) error {
	return nil
}

func (pc *pgmock) ToPgx() *pgxpool.Pool {
	return nil
}
