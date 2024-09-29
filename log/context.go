package log

import "context"

type ctxLogger struct{}

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// loggerFromContext returns logger from context.
func LoggerFromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*Logger); ok {
		return l
	}

	panic("logger not found in context")
}
