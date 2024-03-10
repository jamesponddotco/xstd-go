// Package xslog provides extensions to the standard structured logging
// module of the Go standard library.
package xslog

import (
	"context"
	"log/slog"
)

// LogEvent is a wrapper around slog.Logger.LogAttrs that only logs if logger
// isn't nil. It exists primarily to reduce the amount of ifs in the code.
func LogEvent(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, attrs ...slog.Attr) {
	if logger != nil {
		logger.LogAttrs(ctx, level, msg, attrs...)
	}
}
