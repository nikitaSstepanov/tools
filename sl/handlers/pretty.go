package handlers

import (
	"context"
	"encoding/json"
	"io"
	std "log"
	"log/slog"

	"github.com/fatih/color"
)

type Pretty struct {
	log   *std.Logger
	attrs []slog.Attr
	slog.Handler
}

func NewPretty(out io.Writer, opts *slog.HandlerOptions) *Pretty {
	return &Pretty{
		Handler: slog.NewJSONHandler(out, opts),
		log:     std.New(out, "", 0),
	}
}

func (p *Pretty) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range p.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	p.log.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (p *Pretty) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Pretty{
		Handler: p.Handler,
		log:     p.log,
		attrs:   attrs,
	}
}

func (p *Pretty) WithGroup(name string) slog.Handler {
	return &Pretty{
		Handler: p.Handler.WithGroup(name),
		log:     p.log,
	}
}
