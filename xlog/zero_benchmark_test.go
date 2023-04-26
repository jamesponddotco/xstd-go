package xlog_test

import (
	"bytes"
	"log"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xlog"
)

func BenchmarkZeroLogger_Printf(b *testing.B) {
	var buf bytes.Buffer
	logger := &xlog.ZeroLogger{Output: &buf}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Printf("This is a log message: %d\n", i)
	}
}

func BenchmarkStdLogger_Printf(b *testing.B) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Printf("This is a log message: %d\n", i)
	}
}
