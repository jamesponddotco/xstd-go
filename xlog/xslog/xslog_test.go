package xslog_test

import (
	"context"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xlog/xslog"
)

func TestLogEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		logger    *slog.Logger
		level     slog.Level
		msg       string
		attrs     []slog.Attr
		wantLevel slog.Level
		wantMsg   string
	}{
		{
			name:   "Logger is nil",
			logger: nil,
			level:  slog.LevelInfo,
			msg:    "test message",
			attrs: []slog.Attr{
				slog.String("key", "value"),
			},
			wantLevel: slog.LevelInfo,
			wantMsg:   "",
		},
		{
			name:   "Logger is not nil",
			logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			level:  slog.LevelInfo,
			msg:    "test message",
			attrs: []slog.Attr{
				slog.String("key", "value"),
			},
			wantLevel: slog.LevelInfo,
			wantMsg:   `level=INFO msg="test message" key=value`,
		},
		{
			name:   "Different log level",
			logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			level:  slog.LevelError,
			msg:    "error message",
			attrs: []slog.Attr{
				slog.Int("code", 500),
			},
			wantLevel: slog.LevelError,
			wantMsg:   `level=ERROR msg="error message" code=500`,
		},
		{
			name:      "No attributes",
			logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
			level:     slog.LevelWarn,
			msg:       "warning message",
			attrs:     []slog.Attr{},
			wantLevel: slog.LevelWarn,
			wantMsg:   `level=WARN msg="warning message"`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var sb strings.Builder
			if tt.logger != nil {
				tt.logger = slog.New(slog.NewTextHandler(&sb, nil))
			}

			xslog.LogEvent(context.Background(), tt.logger, tt.level, tt.msg, tt.attrs...)

			if tt.logger == nil {
				if sb.Len() != 0 {
					t.Errorf("LogEvent() logged message when logger is nil")
				}

				return
			}

			got := strings.TrimSpace(sb.String())
			got = regexp.MustCompile(`^time=\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}(?:Z|[+-]\d{2}:\d{2}) `).ReplaceAllString(got, "")

			if got != tt.wantMsg {
				t.Errorf("LogEvent() logged message %q, want %q", got, tt.wantMsg)
			}
		})
	}
}
