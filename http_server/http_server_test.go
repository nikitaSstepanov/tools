package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	e "github.com/nikitaSstepanov/tools/error"
	"github.com/nikitaSstepanov/tools/sl"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	// Define a simple handler for testing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := &Config{
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		Url:             ":8080",
		ShutdownTimeout: 30 * time.Second,
	}

	server := New(handler, cfg)

	// Assertions using testify
	assert.NotNil(t, server, "Expected non-nil Server")
	assert.NotNil(t, server.server, "Expected non-nil http.Server")

	assert.Equal(t, cfg.ReadTimeout, server.server.ReadTimeout, "ReadTimeout mismatch")
	assert.Equal(t, cfg.WriteTimeout, server.server.WriteTimeout, "WriteTimeout mismatch")
	assert.Equal(t, cfg.Url, server.server.Addr, "Addr mismatch")

	assert.NotNil(t, server.notify, "Expected non-nil notify channel")

	assert.Equal(t, cfg.ShutdownTimeout, server.shutdownTimeout, "shutdownTimeout mismatch")
}

func TestServer_Shutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := &Config{
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		Url:             ":8080",
		ShutdownTimeout: 30 * time.Second,
	}

	server := New(handler, cfg)

	go server.Start()

	time.Sleep(100 * time.Millisecond)

	go func() {
		time.Sleep(3 * time.Second)
		server.notify <- e.New("some msg", e.Internal)
	}()

	mockLogger := sl.New(&sl.Config{Type: "discard"})
	ctx := sl.ContextWithLogger(context.Background(), mockLogger)

	err := server.Shutdown(ctx)

	assert.NoError(t, err)
}

func TestServer_Notify(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := &Config{
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		Url:             ":8080",
		ShutdownTimeout: 30 * time.Second,
	}

	server := New(handler, cfg)

	// Start the server
	server.Start()

	// Wait for the notification channel to be ready
	time.Sleep(100 * time.Millisecond)

	notifyChan := server.Notify()

	assert.NotNil(t, notifyChan, "Expected non-nil notification channel")
}
