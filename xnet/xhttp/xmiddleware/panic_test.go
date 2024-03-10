package xmiddleware_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
)

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "with panic",
			handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				panic("forced panic")
			}),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{
  "message": "Internal server error. Please try again later.",
  "code": 500
}`,
		},
		{
			name: "without panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}),
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				buf     bytes.Buffer
				handler = slog.NewJSONHandler(&buf, nil)
				logger  = slog.New(handler)
			)

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			middleware := xmiddleware.PanicRecovery(logger, tt.handler)
			middleware.ServeHTTP(w, r)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
			}

			// Validate logger behavior for panic case
			if tt.name != "with panic" {
				return
			}

			var logs []map[string]any

			for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
				if len(line) == 0 {
					continue
				}

				var logEntry map[string]any
				if err := json.Unmarshal(line, &logEntry); err != nil {
					t.Fatal(err)
				}

				logs = append(logs, logEntry)
			}

			if len(logs) != 2 {
				t.Errorf("expected 1 log entry, got %d", len(logs))

				return
			}

			if logs[0]["error"] != "forced panic" {
				t.Errorf("log entry did not match expected values")
			}
		})
	}
}
