package handlers

import (
	"io"
	"log/slog"

	"github.com/golang-cz/devslog"
)

func NewDevSlog(out io.Writer, opts *slog.HandlerOptions) slog.Handler {
	devOpts := &devslog.Options{
		HandlerOptions:    opts,
		MaxSlicePrintSize: 5,
		SortKeys:          true,
		NewLineAfterLog:   true,
		StringerFormatter: true,
	}

	return devslog.NewHandler(out, devOpts)
}
