package sl

import (
	"log/slog"
	"os"

	"github.com/nikitaSstepanov/tools/sl/handlers"
)

// Config holds the configuration settings for the logger.
// It includes options for logging level, output format, and other behavior.
type Config struct {
	// Level specifies the logging level (e.g., info, debug, error).
	// It can be set via YAML configuration or environment variable.
	Level string `yaml:"level" env:"LOGGER_LEVEL" env-default:"info"`

	// AddSource indicates whether to include the source of the log message (e.g., file and line number).
	// This can be controlled through YAML or an environment variable.
	AddSource bool ` yaml:"add_source" env:"LOGGER_ADD_SOURCE" env-default:"true"`

	// IsJSON determines if the log output should be in JSON format.
	// If true, logs will be structured as JSON; otherwise, they will be plain text.
	IsJSON bool `yaml:"is_json" env:"LOGGER_IS_JSON" env-default:"true"`

	// Writer specifies where the logs should be written.
	// It can be a file path or a predefined output like "stdout".
	Writer string `yaml:"writer" env:"LOGGER_WRITER" env-default:"stderr"`

	// OutPath is the path to the output file if logging to a file.
	// If left empty, logs will go to the Writer specified.
	OutPath string `yaml:"out_path" env:"LOGGER_OUT_PATH" env-default:""`

	// SetDefault indicates whether to set default logger options.
	// This can be used to ensure that certain configurations are applied automatically.
	SetDefault bool `yaml:"set_default" env:"LOGGER_SET_DEFAULT" env-default:"true"`

	// Type defines the type of logger to use.
	// It can be one of the following:
	// - "pretty" or "dev": for human-readable logs with color and formatting.
	// - "discard": to ignore all log messages.
	// - "default": for standard logging behavior.
	Type string `yaml:"type" env:"LOGGER_TYPE" env-default:"default"`
}

// New returns the slog.Logger with the specified configuration.
func New(cfg *Config) *Logger {
	handler := setupHandler(cfg)

	logger := slog.New(handler)

	if cfg.SetDefault {
		SetDefault(logger)
	}

	return logger
}

func setupHandler(cfg *Config) Handler {
	level := setLoggerLevel(cfg.Level)

	opts := setHandlerOptions(level, cfg.AddSource)

	out := setOut(cfg)

	var handler Handler

	switch cfg.Type {

	case DevLogger:
		handler = handlers.NewDevSlog(out, opts)

	case PrettyLogger:
		handler = handlers.NewPretty(out, opts)

	case DiscardLogger:
		handler = handlers.NewDiscard()

	default:
		if cfg.IsJSON {
			handler = NewJSONHandler(out, opts)
		} else {
			handler = NewTextHandler(out, opts)
		}

	}

	return handler
}

func setLoggerLevel(lvl string) Level {
	var level Level

	switch lvl {

	case "debug":
		level = -4
	case "info":
		level = 0
	case "warn":
		level = 4
	case "error":
		level = 8
	default:
		level = 0

	}

	return level
}

func setHandlerOptions(level Level, AddSource bool) *HandlerOptions {
	return &HandlerOptions{AddSource: AddSource, Level: level}
}

func setOut(cfg *Config) *os.File {
	if cfg.Writer == FileOut {
		return getLogFile(cfg.OutPath)
	}

	return os.Stderr
}

func getLogFile(path string) *os.File {
	if path == "" {
		path = "logs"
	}

	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(path+"/all.log", os.O_CREATE|os.O_RDWR, 0o0644)
	if err != nil {
		panic(err)
	}

	return logFile
}
