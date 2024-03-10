package xmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
)

func TestChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		middlewares      []func(http.Handler) http.Handler
		requestMethod    string
		expectStatusCode int
	}{
		{
			name:             "No middlewares",
			middlewares:      []func(http.Handler) http.Handler{},
			requestMethod:    http.MethodGet,
			expectStatusCode: http.StatusOK,
		},
		{
			name: "Multiple middlewares",
			middlewares: []func(http.Handler) http.Handler{
				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Middleware", "1")
						next.ServeHTTP(w, r)
					})
				},

				func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Middleware", "2")
						next.ServeHTTP(w, r)
					})
				},
			},
			requestMethod:    http.MethodGet,
			expectStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(tt.requestMethod, "/", http.NoBody)
			)

			var (
				handlerToChain = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				})
				chain = xmiddleware.Chain(handlerToChain, tt.middlewares...)
			)

			chain.ServeHTTP(w, r)

			if w.Code != tt.expectStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectStatusCode, w.Code)
			}
		})
	}
}
