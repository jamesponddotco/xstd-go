package xhttp_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
)

func TestNewRetryPolicy(t *testing.T) {
	t.Parallel()

	policy := xhttp.NewRetryPolicy()

	if policy.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", policy.MaxRetries)
	}

	if policy.MinRetryDelay != 1*time.Second {
		t.Errorf("Expected MinRetryDelay to be 1 second, got %s", policy.MinRetryDelay)
	}

	if policy.MaxRetryDelay != 30*time.Second {
		t.Errorf("Expected MaxRetryDelay to be 30 seconds, got %s", policy.MaxRetryDelay)
	}

	// Check IsRetryable function with a retryable status code
	resp := &http.Response{
		StatusCode: http.StatusRequestTimeout,
	}

	if !policy.IsRetryable(resp, nil) {
		t.Errorf("Expected IsRetryable to return true for status code %d", resp.StatusCode)
	}

	// Check IsRetryable function with a non-retryable status code
	resp = &http.Response{
		StatusCode: http.StatusOK,
	}

	if policy.IsRetryable(resp, nil) {
		t.Errorf("Expected IsRetryable to return false for status code %d", resp.StatusCode)
	}

	if !policy.IsRetryable(nil, errors.New("network error")) {
		t.Error("Expected IsRetryable to return true for an error")
	}
}

func TestRetryPolicy_Wait(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		maxRetries      int
		minRetryDelay   time.Duration
		maxRetryDelay   time.Duration
		attempt         int
		contextFunc     func() (context.Context, context.CancelFunc)
		expectedWaitMax time.Duration
		expectedErr     error
	}{
		{
			name:          "Context Canceled",
			maxRetries:    3,
			minRetryDelay: 1 * time.Second,
			maxRetryDelay: 30 * time.Second,
			attempt:       1,
			contextFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())

				cancel() // Cancel the context immediately.

				return ctx, cancel
			},
			expectedWaitMax: 1 * time.Millisecond,
			expectedErr:     context.Canceled,
		},
		{
			name:          "Maximum Retry Delay Exceeded",
			maxRetries:    3,
			minRetryDelay: 1 * time.Second,
			maxRetryDelay: 2 * time.Second,
			attempt:       4, // This should result in a calculated delay > maxRetryDelay
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			expectedWaitMax: 3 * time.Second,
			expectedErr:     nil,
		},
		{
			name:          "Minimum Retry Delay",
			maxRetries:    3,
			minRetryDelay: 1 * time.Second,
			maxRetryDelay: 30 * time.Second,
			attempt:       0, // First attempt
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 10*time.Second)
			},
			expectedWaitMax: 2 * time.Second,
			expectedErr:     nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			policy := xhttp.RetryPolicy{
				MaxRetries:    tt.maxRetries,
				MinRetryDelay: tt.minRetryDelay,
				MaxRetryDelay: tt.maxRetryDelay,
			}

			ctx, cancel := tt.contextFunc()
			defer cancel()

			var (
				start    = time.Now()
				err      = policy.Wait(ctx, tt.attempt)
				duration = time.Since(start)
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("Wait() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if duration > tt.expectedWaitMax {
				t.Errorf("Wait() duration = %v, expected maximum %v", duration, tt.expectedWaitMax)
			}
		})
	}
}
