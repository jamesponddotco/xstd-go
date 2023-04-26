package xlog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// DefaultZeroLogger is a logger that writes to os.Stderr.
var DefaultZeroLogger = &ZeroLogger{Output: os.Stderr} //nolint:gochecknoglobals // global by desin

// ZeroLogger is a small logger that doesn't allocate any memory but implements
// only a subset of the Logger interface.
type ZeroLogger struct {
	Output io.Writer
	mu     sync.Mutex
}

// Printf calls l.Output to print to the logger. Arguments are handled in the
// manner of fmt.Printf.
func (l *ZeroLogger) Printf(format string, v ...any) {
	l.mu.Lock()
	fmt.Fprintf(l.Output, format, v...)
	l.mu.Unlock()
}

// Printf calls Output to print to the standard zero logger. Arguments are
// handled in the manner of fmt.Printf.
func Printf(format string, v ...any) {
	DefaultZeroLogger.Printf(format, v...)
}
