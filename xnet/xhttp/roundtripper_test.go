package xhttp_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

func TestRetryRoundTripper_RoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		responses      []int
		wantAttempts   int
		wantFinal      int
		retries        int
		expectErr      bool
		delay          time.Duration
		timeoutContext time.Duration
	}{
		{
			name: "Successful on first attempt",
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
			name: "Successful on third attempt",
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
			name: "Exceeds max retries",
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
			name: "Context canceled before max retries",
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				mu       sync.Mutex
				requests = 0
				srv      = httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, _ *http.Request) {
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
				MaxRetries: tt.retries,
				IsRetryable: func(resp *http.Response, _ error) bool {
					return resp.StatusCode >= 500
				},
				MinRetryDelay: 1 * time.Millisecond,
				MaxRetryDelay: 1 * time.Millisecond,
			}

			rrt := xhttp.NewRetryRoundTripper(policy, slog.New(slog.NewTextHandler(io.Discard, nil)))

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeoutContext)
			defer cancel()

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, http.NoBody)

			resp, err := rrt.RoundTrip(req)
			if resp != nil {
				defer resp.Body.Close()
			}

			if tt.expectErr != (err != nil) {
				t.Errorf("expectErr=%v, got: %v", tt.expectErr, err)
			}

			if err == nil && resp.StatusCode != tt.wantFinal {
				t.Errorf("unexpected status code: got %v, want %v", resp.StatusCode, tt.wantFinal)
			}

			mu.Lock()
			got := requests
			mu.Unlock()

			if got != tt.wantAttempts {
				t.Errorf("unexpected number of attempts: got %d, want %d", got, tt.wantAttempts)
			}
		})
	}
}
