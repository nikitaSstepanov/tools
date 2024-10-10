package sl

import (
	"context"
	"log/slog"
)

type ctxLogger struct{}

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}


func loggerFromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*Logger); ok {
		return l
	}

	return slog.Default()
}
