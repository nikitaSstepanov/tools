package ctx

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	log := slog.Default()
	c := New(log)

	assert.Equal(t, c.Logger(), log)
	assert.Equal(t, c.SlHandler(), log.Handler())
}

func TestNewWithCtx(t *testing.T) {
	base := context.TODO()
	log := slog.Default()

	key := "key"
	value := "value"

	base = context.WithValue(base, key, value)
	base, cancel := context.WithTimeout(base, 5*time.Second)
	defer cancel()

	c := NewWithCtx(base, log)

	assert.Equal(t, c.Logger(), log)
	assert.Equal(t, c.SlHandler(), log.Handler())
	assert.Equal(t, c.Value(key), value)
	assert.Equal(t, c.Done(), base.Done())
	assert.Equal(t, c.Err(), base.Err())

	deadline1, ok1 := base.Deadline()
	deadline2, ok2 := c.Deadline()

	assert.Equal(t, deadline1, deadline2)
	assert.Equal(t, ok1, ok2)
}

func TestGetValue(t *testing.T) {
	c := New(slog.Default())

	key := "key"
	value := "value"

	c.AddValue(key, value, true)

	assert.Equal(t, c.GetValue(key).Val, value)
}

func TestGetValues(t *testing.T) {
	c := New(slog.Default())

	key := "key"
	value := "value"

	c.AddValue(key, value, true)

	vals := c.GetValues()

	assert.Equal(t, vals[key].Val, value)
}

func TestAddValue(t *testing.T) {
	c := New(slog.Default())

	key := "key"
	value := "value"

	c.AddValue(key, value, true)

	assert.Equal(t, c.GetValue(key).Val, value)
}

func TestLogger(t *testing.T) {
	log := slog.Default()
	c := New(log)

	assert.Equal(t, log, c.Logger())
}

func TestSlHandler(t *testing.T) {
	log := slog.Default()
	c := New(log)

	assert.Equal(t, log.Handler(), c.SlHandler())
}

func TestAddErr(t *testing.T) {
	c := New(slog.Default())

	err := errors.New("some error")
	c.AddErr(err)

	assert.Equal(t, err, c.GetErr())
}

func TestGetErr(t *testing.T) {
	c := New(slog.Default())

	err := errors.New("some error")
	c.AddErr(err)

	assert.Equal(t, err, c.GetErr())
}

func TestHasErr(t *testing.T) {
	c := New(slog.Default())

	err := errors.New("some error")
	c.AddErr(err)

	assert.Equal(t, true, c.HasErr())
}

func TestDeadline(t *testing.T) {
	base := context.TODO()
	log := slog.Default()

	base, cancel := context.WithTimeout(base, 5*time.Second)
	defer cancel()

	c := NewWithCtx(base, log)

	deadline1, ok1 := base.Deadline()
	deadline2, ok2 := c.Deadline()

	assert.Equal(t, deadline1, deadline2)
	assert.Equal(t, ok1, ok2)
}

func TestDone(t *testing.T) {
	base := context.TODO()
	log := slog.Default()

	base, cancel := context.WithTimeout(base, 5*time.Second)
	defer cancel()

	c := NewWithCtx(base, log)

	assert.Equal(t, c.Done(), base.Done())
}

func TestErr(t *testing.T) {
	base := context.TODO()
	log := slog.Default()

	base, cancel := context.WithTimeout(base, 5*time.Second)
	defer cancel()

	c := NewWithCtx(base, log)

	assert.Equal(t, c.Err(), base.Err())
}

func TestValue(t *testing.T) {
	base := context.TODO()
	log := slog.Default()

	key := "key"
	value := "value"

	base = context.WithValue(base, key, value)

	c := NewWithCtx(base, log)

	assert.Equal(t, c.Value(key), value)
}
