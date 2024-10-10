package sl

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	ctx := context.Background()

	logger := New(&Config{})

	if logger == nil {
		t.Errorf("logger should not be nil")
	}

	logger = New(&Config{IsJSON: true})

	_ = New(&Config{IsJSON: false})
	if slog.Default() == logger {
		t.Errorf("logger should NOT be default logger")
	}

	logger = New(&Config{Level: "info"})
	if logger.Handler().Enabled(ctx, LevelDebug) {
		t.Errorf("logger should NOT be enabled for debug level")
	}

	logger = New(&Config{Level: "warn"})
	if logger.Handler().Enabled(ctx, LevelDebug) {
		t.Errorf("logger should NOT be enabled for debug level")
	}

	logger = New(&Config{Level: "debug"})
	enabled := []Level{LevelDebug, LevelInfo, LevelWarn, LevelError}
	for _, level := range enabled {
		if !logger.Handler().Enabled(ctx, level) {
			t.Errorf("logger should be enabled for all levels")
		}
	}

	logger = New(&Config{Level: "abcdef"})
	if !logger.Handler().Enabled(ctx, LevelInfo) {
		t.Errorf("logger should be enabled for info level")
	}

	if logger.Handler().Enabled(ctx, LevelDebug) {
		t.Errorf("logger should NOT be enabled for info level")
	}

	logger = New(&Config{SetDefault: false})
	if slog.Default() == logger {
		t.Errorf("logger should NOT be default logger")
	}

	if logger == Default() {
		t.Errorf("logger should NOT be default logger")
	}

	assert.Equal(t, slog.Default(), L(ctx))

	logger = New(&Config{AddSource: false, IsJSON: false})
	ctx = ContextWithLogger(ctx, logger)

	if loggerFromContext(ctx) != logger {
		t.Errorf("logger should be from context logger")
	}
}
