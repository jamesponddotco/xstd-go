package xhttp_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

func TestRetryRoundTripper_RoundTrip(t *testing.T) { //nolint:gocognit // I see no way of making this less complex.
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           io.Reader
		responses      []int
		wantAttempts   int
		wantFinal      int
		retries        int
		expectErr      bool
		delay          time.Duration
		timeoutContext time.Duration
	}{
		{
			name:   "GET Successful on first attempt",
			method: http.MethodGet,
			body:   http.NoBody,
			responses: []int{
				http.StatusOK,
			},
			wantAttempts:   1,
			wantFinal:      http.StatusOK,
			retries:        3,
			expectErr:      false,
			delay:          0,
			timeoutContext: 5 * time.Second,
		},
		{
			name:   "GET Successful on third attempt",
			method: http.MethodGet,
			body:   http.NoBody,
			responses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusOK,
			},
			wantAttempts:   3,
			wantFinal:      http.StatusOK,
			retries:        3,
			expectErr:      false,
			delay:          0,
			timeoutContext: 5 * time.Second,
		},
		{
			name:   "GET Exceeds max retries",
			method: http.MethodGet,
			body:   http.NoBody,
			responses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
			wantAttempts:   3,
			wantFinal:      http.StatusServiceUnavailable,
			retries:        2,
			expectErr:      true,
			delay:          0,
			timeoutContext: 5 * time.Second,
		},
		{
			name:   "GET Context canceled before max retries",
			method: http.MethodGet,
			body:   http.NoBody,
			responses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
			wantAttempts:   0,
			wantFinal:      http.StatusInternalServerError,
			retries:        3,
			expectErr:      true,
			delay:          2 * time.Second,
			timeoutContext: 1 * time.Second,
		},
		{
			name:   "POST Successful on first attempt",
			method: http.MethodPost,
			body:   strings.NewReader(`{"key":"value"}`),
			responses: []int{
				http.StatusCreated,
			},
			wantAttempts:   1,
			wantFinal:      http.StatusCreated,
			retries:        3,
			expectErr:      false,
			delay:          0,
			timeoutContext: 5 * time.Second,
		},
		{
			name:   "POST Successful on third attempt",
			method: http.MethodPost,
			body:   strings.NewReader(`{"key":"value"}`),
			responses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusCreated,
			},
			wantAttempts:   3,
			wantFinal:      http.StatusCreated,
			retries:        3,
			expectErr:      false,
			delay:          0,
			timeoutContext: 5 * time.Second,
		},
		{
			name:   "POST Exceeds max retries",
			method: http.MethodPost,
			body:   strings.NewReader(`{"key":"value"}`),
			responses: []int{
				http.StatusInternalServerError,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
			wantAttempts:   3,
			wantFinal:      http.StatusServiceUnavailable,
			retries:        2,
			expectErr:      true,
			delay:          0,
			timeoutContext: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				mu       sync.Mutex
				requests = 0
				srv      = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != tt.method {
						t.Fatalf("Expected request method %q, got %q", tt.method, r.Method)
					}

					if tt.method == http.MethodPost {
						body, err := io.ReadAll(r.Body)
						if err != nil {
							t.Fatalf("Expected no error, got %v", err)
						}
						defer r.Body.Close()

						if string(body) != `{"key":"value"}` {
							t.Errorf("Expected request body %q, got %q", `{"key":"value"}`, string(body))
						}
					}

					time.Sleep(tt.delay)

					respIdx := requests
					if respIdx >= len(tt.responses) {
						respIdx = len(tt.responses) - 1
					}

					w.WriteHeader(tt.responses[respIdx])

					mu.Lock()
					requests++
					mu.Unlock()
				}))
			)

			defer srv.Close()

			policy := &xhttp.RetryPolicy{
				MaxRetries:    tt.retries,
				MinRetryDelay: 1 * time.Millisecond,
				MaxRetryDelay: 1 * time.Millisecond,
				IsRetryable: func(resp *http.Response, _ error) bool {
					return resp.StatusCode >= 500
				},
			}

			rrt := xhttp.NewRetryRoundTripper(policy, slog.New(slog.NewTextHandler(io.Discard, nil)))

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeoutContext)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, tt.method, srv.URL, tt.body)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			resp, err := rrt.RoundTrip(req)
			if err != nil {
				if !tt.expectErr {
					t.Errorf("Unexpected error: %v", err)
				}

				return
			}
			defer resp.Body.Close()

			if tt.expectErr {
				t.Error("Expected error, got nil")
			}

			if resp.StatusCode != tt.wantFinal {
				t.Errorf("Expected final status code %d, got %d", tt.wantFinal, resp.StatusCode)
			}

			mu.Lock()
			got := requests
			mu.Unlock()

			if got != tt.wantAttempts {
				t.Errorf("Expected %d attempts, got %d", tt.wantAttempts, got)
			}
		})
	}
}
