package xslog_test

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xlog/xslog"
)

func TestSlogAdapter_Write(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		level  slog.Level
		input  string
		output string
	}{
		{
			name:   "Test Error Level",
			level:  slog.LevelError,
			input:  "This is an error",
			output: "This is an error",
		},
		{
			name:   "Test Info Level",
			level:  slog.LevelInfo,
			input:  "This is info",
			output: "This is info",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				buf         bytes.Buffer
				slogHandler = slog.NewJSONHandler(&buf, nil)
				slogger     = slog.New(slogHandler)
			)

			adapter := &xslog.SlogAdapter{
				Logger: slogger,
				Level:  tt.level,
			}

			logger := log.New(adapter, "", 0)
			logger.Print(tt.input)

			got := buf.Bytes()

			var outputMap map[string]any
			if err := json.Unmarshal(got, &outputMap); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if outputMap["msg"] != tt.output {
				t.Errorf("got %q, want %q", outputMap["msg"], tt.output)
			}
		})
	}
}
