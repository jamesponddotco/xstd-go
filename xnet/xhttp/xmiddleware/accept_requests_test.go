package xmiddleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
)

func TestAcceptRequests(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		methods            []string
		requestMethod      string
		expectStatusCode   int
		expectAllowMethods string
	}{
		{
			name:             "GET allowed",
			methods:          []string{http.MethodGet},
			requestMethod:    http.MethodGet,
			expectStatusCode: http.StatusOK,
		},
		{
			name:               "POST disallowed",
			methods:            []string{http.MethodGet},
			requestMethod:      http.MethodPost,
			expectStatusCode:   http.StatusMethodNotAllowed,
			expectAllowMethods: "GET",
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
				r = httptest.NewRequest(tt.requestMethod, "/", http.NoBody)
			)

			middleware := xmiddleware.AcceptRequests(tt.methods, logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			middleware.ServeHTTP(w, r)

			if w.Code != tt.expectStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectStatusCode, w.Code)
			}

			if tt.expectAllowMethods != "" {
				allowHeader := w.Header().Get("Allow")
				if allowHeader != tt.expectAllowMethods {
					t.Errorf("Expected Allow header to be %s, got %s", tt.expectAllowMethods, allowHeader)
				}
			}
		})
	}
}
