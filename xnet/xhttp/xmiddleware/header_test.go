package xmiddleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcontext"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xmiddleware"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		userAgent      string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing User-Agent",
			userAgent:      "",
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
  "message": "User agent is missing. Please provide a valid user agent.",
  "code": 400
}`,
		},
		{
			name:           "valid User-Agent",
			userAgent:      "TestAgent",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
	}

	for _, tt := range tests {
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

			r.Header.Set("User-Agent", tt.userAgent)

			middleware := xmiddleware.UserAgent(logger, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK")) //nolint:errcheck // safe to ignore, we're testing
			}))
			middleware.ServeHTTP(w, r)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}

func TestPrivacyPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		uri         string
		expectedURI string
	}{
		{
			name:        "valid URI",
			uri:         "https://example.com/privacy",
			expectedURI: "https://example.com/privacy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			middleware := xmiddleware.PrivacyPolicy(tt.uri, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
			middleware.ServeHTTP(w, r)

			if w.Header().Get("Privacy-Policy") != tt.expectedURI {
				t.Errorf("expected header %s, got %s", tt.expectedURI, w.Header().Get("Privacy-Policy"))
			}
		})
	}
}

func TestTermsOfService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		uri         string
		expectedURI string
	}{
		{
			name:        "valid URI",
			uri:         "https://example.com/terms",
			expectedURI: "https://example.com/terms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			middleware := xmiddleware.TermsOfService(tt.uri, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
			middleware.ServeHTTP(w, r)

			if w.Header().Get("Terms-Of-Service") != tt.expectedURI {
				t.Errorf("expected header %s, got %s", tt.expectedURI, w.Header().Get("Terms-Of-Service"))
			}
		})
	}
}

func TestRequestID(t *testing.T) {
	t.Parallel()

	const contextKeyRequestID xcontext.ContextKey = "request_id"

	tests := []struct {
		name           string
		existingID     string
		validateOutput func(*testing.T, string)
	}{
		{
			name:       "with existing request ID",
			existingID: "test-request-id",
			validateOutput: func(t *testing.T, id string) {
				t.Helper()

				if id != "test-request-id" {
					t.Errorf("expected request ID %q, got %q", "test-request-id", id)
				}
			},
		},
		{
			name:       "without existing request ID",
			existingID: "",
			validateOutput: func(t *testing.T, id string) {
				t.Helper()

				if len(id) != 36 {
					t.Errorf("expected UUID length of 36, got %d", len(id))
				}

				matched := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id)
				if !matched {
					t.Errorf("invalid UUID format: %s", id)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			)

			if tt.existingID != "" {
				r.Header.Set(xhttp.RequestID, tt.existingID)
			}

			var contextID string

			middleware := xmiddleware.RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				id, ok := r.Context().Value(contextKeyRequestID).(string)
				if !ok {
					t.Error("request ID not found in context")
				}

				contextID = id
			}))

			middleware.ServeHTTP(w, r)

			headerID := w.Header().Get(xhttp.RequestID)

			tt.validateOutput(t, headerID)

			if headerID != contextID {
				t.Errorf("header ID %q does not match context ID %q", headerID, contextID)
			}
		})
	}
}
