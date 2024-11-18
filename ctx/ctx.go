package ctx

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Context interface {
	Logger() *slog.Logger
	SlHandler() slog.Handler
	AddValue(key string, val interface{}, forLog bool)
	GetValue(key string) *Value
	GetValues() map[string]Value
	AddErr(err error)
	GetErr() error
	HasErr() bool

	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

type ctx struct {
	log       *slog.Logger
	slHandler slog.Handler
	data      map[string]Value
	errors    []error
	mu        sync.Mutex
	base      context.Context
}

type Value struct {
	Val   interface{}
	Share bool
}

func New(log *slog.Logger) Context {
	return &ctx{
		log:       slog.New(log.Handler()),
		slHandler: log.Handler(),
		data:      make(map[string]Value),
		base:      context.TODO(),
	}
}

func NewWithCtx(base context.Context, log *slog.Logger) Context {
	return &ctx{
		log:       slog.New(log.Handler()),
		slHandler: log.Handler(),
		data:      make(map[string]Value),
		base:      base,
	}
}

func (c *ctx) GetValue(key string) *Value {
	value := c.data[key]
	if value.Val == nil {
		return nil
	}

	return &value
}

func (c *ctx) GetValues() map[string]Value {
	return c.data
}

func (c *ctx) AddValue(key string, val interface{}, share bool) {
	c.data[key] = Value{
		val, share,
	}

	if share {
		(*c.log) = *c.log.With(key, val)
	}
}

func (c *ctx) Logger() *slog.Logger {
	return c.log
}

func (c *ctx) SlHandler() slog.Handler {
	return c.slHandler
}

func (c *ctx) AddErr(err error) {
	c.mu.Lock()
	c.errors = append(c.errors, err)
	c.mu.Unlock()
}

func (c *ctx) GetErr() error {
	c.mu.Lock()

	if len(c.errors) == 0 {
		return nil
	}

	err := c.errors[0]

	if len(c.errors) > 1 {
		c.errors = c.errors[1:]
	}

	c.mu.Unlock()

	return err
}

func (c *ctx) HasErr() bool {
	c.mu.Lock()
	errCount := len(c.errors)
	c.mu.Unlock()

	return errCount > 0
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
