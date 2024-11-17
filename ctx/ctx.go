package ctx

import (
	"context"
	"log/slog"
	"time"
)

type Context interface {
	Add(key string, val interface{}, forLog bool)
	Get(key string) *Value
	Log() *slog.Logger

	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

type ctx struct {
	log  *slog.Logger
	data map[string]Value
	base context.Context
}

type Value struct {
	val   interface{}
	share bool
}

func New(log *slog.Logger) Context {
	return &ctx{
		log:  slog.New(log.Handler()),
		data: make(map[string]Value),
		base: context.TODO(),
	}
}

func NewWithCtx(base context.Context, log *slog.Logger) Context {
	return &ctx{
		log:  slog.New(log.Handler()),
		data: make(map[string]Value),
		base: base,
	}
}

func (c *ctx) Get(key string) *Value {
	value := c.data[key]
	if value.val == nil {
		return nil
	}

	return &value
}

func (c *ctx) Add(key string, val interface{}, share bool) {
	c.data[key] = Value{
		val, share,
	}

	if share {
		(*c.log) = *c.log.With(key, val)
	}
}

func (c *ctx) Log() *slog.Logger {
	return c.log
}

func (c *ctx) Deadline() (deadline time.Time, ok bool) {
	return c.base.Deadline()
}

func (c *ctx) Done() <-chan struct{} {
	return c.base.Done()
}

func (c *ctx) Err() error {
	return c.base.Err()
}

func (c *ctx) Value(key any) any {
	return c.base.Value(key)
}
