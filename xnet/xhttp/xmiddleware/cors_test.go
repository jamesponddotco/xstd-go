package xmiddleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		config      *xmiddleware.CORSConfig
		request     *http.Request
		statusCode  int
		allowOrigin string
	}{
		{
			name: "Allowed origin",
			config: &xmiddleware.CORSConfig{
				AllowedOrigins: []string{"http://allowed.com"},
			},
			request:     httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			statusCode:  http.StatusOK,
			allowOrigin: "http://allowed.com",
		},
		{
			name: "Wildcard origin",
			config: &xmiddleware.CORSConfig{
				AllowedOrigins: []string{"*"},
			},
			request:     httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			statusCode:  http.StatusOK,
			allowOrigin: "*",
		},
		{
			name:        "Nil config",
			config:      nil,
			request:     httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			statusCode:  http.StatusOK,
			allowOrigin: "*",
		},
		{
			name: "Disallowed origin",
			config: &xmiddleware.CORSConfig{
				AllowedOrigins: []string{"http://allowed.com"},
			},
			request:     httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			statusCode:  http.StatusForbidden,
			allowOrigin: "",
		},
		{
			name: "AllowCredentials with specific origin",
			config: &xmiddleware.CORSConfig{
				AllowedOrigins:   []string{"http://allowed.com"},
				AllowCredentials: true,
			},
			request:     httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			statusCode:  http.StatusOK,
			allowOrigin: "http://allowed.com",
		},
		{
			name: "AllowCredentials with wildcard origin",
			config: &xmiddleware.CORSConfig{
				AllowedOrigins:   []string{"*"},
				AllowCredentials: true,
			},
			request:     httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			statusCode:  http.StatusOK,
			allowOrigin: "*",
		},
		{
			name: "OPTIONS request",
			config: &xmiddleware.CORSConfig{
				AllowedOrigins: []string{"*"},
			},
			request:     httptest.NewRequest(http.MethodOptions, "/", http.NoBody),
			statusCode:  http.StatusNoContent,
			allowOrigin: "*",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				buf           bytes.Buffer
				loggerHandler = slog.NewJSONHandler(&buf, nil)
				logger        = slog.New(loggerHandler)
			)

			tt.request.Header.Set(xhttp.Origin, tt.allowOrigin)
			recorder := httptest.NewRecorder()

			handler := xmiddleware.CORS(tt.config, logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(recorder, tt.request)

			if recorder.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, recorder.Code)
			}

			allowOrigin := recorder.Header().Get(xhttp.AccessControlAllowOrigin)
			if allowOrigin != tt.allowOrigin {
				t.Errorf("expected Allow-Origin header to be %s, got %s", tt.allowOrigin, allowOrigin)
			}
		})
	}
}
