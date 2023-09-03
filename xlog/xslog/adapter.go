package xslog

import (
	"context"
	"log/slog"
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

// SlogAdapter adapts an slog.Logger instance to the standard log.Logger
// interface.
type SlogAdapter struct {
	// Logger is the underlying slog.Logger instance.
	Logger *slog.Logger

	// Level is the minimum logging level for the adapter to emit.
	Level slog.Level
}

// Write implements the io.Writer interface for SlogAdapter.
func (sa *SlogAdapter) Write(p []byte) (n int, _ error) { //nolint:unparam // Write is an interface, so we can't remove the error return.
	message := strings.TrimSpace(xunsafe.BytesToString(p))

	sa.Logger.LogAttrs(
		context.Background(),
		sa.Level,
		message,
	)

	return len(p), nil
}
