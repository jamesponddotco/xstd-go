package xlog_test

import (
	"bytes"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xlog"
)

func TestZeroLogger_Printf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format string
		values []any
		want   string
	}{
		{
			name:   "Basic printf",
			format: "Hello, %s!\n",
			values: []any{"world"},
			want:   "Hello, world!\n",
		},
		{
			name:   "Multiple types",
			format: "This is %s, version %d.%d, release in %d.\n",
			values: []any{"Go", 1, 18, 2022},
			want:   "This is Go, version 1.18, release in 2022.\n",
		},
		{
			name:   "No values",
			format: "Just a message.\n",
			values: []any{},
			want:   "Just a message.\n",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := &xlog.ZeroLogger{Output: &buf}

			logger.Printf(tt.format, tt.values...)

			got := buf.String()
			if got != tt.want {
				t.Errorf("ZeroLogger.Printf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format string
		values []any
		want   string
	}{
		{
			name:   "Basic printf",
			format: "Hello, %s!\n",
			values: []any{"world"},
			want:   "Hello, world!\n",
		},
		{
			name:   "Multiple types",
			format: "This is %s, version %d.%d, release in %d.\n",
			values: []any{"Go", 1, 18, 2022},
			want:   "This is Go, version 1.18, release in 2022.\n",
		},
		{
			name:   "No values",
			format: "Just a message.\n",
			values: []any{},
			want:   "Just a message.\n",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			xlog.DefaultZeroLogger.Output = &buf

			xlog.Printf(tt.format, tt.values...)

			got := buf.String()
			if got != tt.want {
				t.Errorf("Printf() = %v, want %v", got, tt.want)
			}
		})
	}
}
