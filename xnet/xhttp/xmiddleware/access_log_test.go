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

func TestAccessLog(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		path           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			path:           "/",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			path:           "/create",
			statusCode:     http.StatusCreated,
			expectedStatus: http.StatusCreated,
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
				r = httptest.NewRequest(tt.method, tt.path, http.NoBody)
			)

			middleware := xmiddleware.AccessLog(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			middleware.ServeHTTP(w, r)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Validate logger behavior
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

			if len(logs) != 1 {
				t.Errorf("expected 1 log entry, got %d", len(logs))

				return
			}

			if logs[0]["status"] != float64(tt.expectedStatus) {
				t.Errorf("log entry did not match expected values")
			}
		})
	}
}
