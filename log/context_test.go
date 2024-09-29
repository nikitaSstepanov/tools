package log

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWithLogger(t *testing.T) {
	ctx := context.Background()
	logger := New(&Config{})
	ctxWithLogger := ContextWithLogger(ctx, logger)

	// Extract logger from new context
	extractedLogger, ok := ctxWithLogger.Value(ctxLogger{}).(*slog.Logger)
	if !ok || extractedLogger != logger {
		t.Errorf("Logger was not properly added to context")
	}
}

func TestLoggerFromContext_WithLogger(t *testing.T) {
	ctx := context.Background()
	logger := New(&Config{})
	ctxWithLogger := ContextWithLogger(ctx, logger)

	// Extract logger using loggerFromContext
	extractedLogger := LoggerFromContext(ctxWithLogger)
	if extractedLogger != logger {
		t.Errorf("Did not retrieve correct logger from context")
	}
}

func TestLoggerFromContext_NoLogger(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() { LoggerFromContext(ctx) })
}
