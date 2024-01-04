package xhttp

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const (
	_backoffFactor float64 = 2.0
	_jitterFactor  float64 = 0.1
)

// RetryPolicy defines a policy for retrying HTTP requests.
type RetryPolicy struct {
	// IsRetryable determines whether a given response and error combination
	// should be retried.
	IsRetryable func(*http.Response, error) bool

	// MaxRetries is the maximum number of times a request will be retried.
	MaxRetries int

	// MinRetryDelay is the minimum duration to wait before retrying a request.
	MinRetryDelay time.Duration

	// MaxRetryDelay is the maximum duration to wait before retrying a request.
	MaxRetryDelay time.Duration
}

// NewRetryPolicy returns a new RetryPolicy with sane default values.
//
// The default values are:
// - IsRetryable: Retries if the status code is 408, 429, 502, 503, or 504.
// - MaxRetries: 3.
// - MinRetryDelay: 1 second.
// - MaxRetryDelay: 30 seconds.
func NewRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		IsRetryable:   defaultIsRetryable,
		MaxRetries:    3,
		MinRetryDelay: 1 * time.Second,
		MaxRetryDelay: 30 * time.Second,
	}
}

// Wait calculates the time to wait before the next retry attempt and blocks
// until it is time to retry the request or the context is canceled. It
// incorporates an exponential backoff strategy with jitter to prevent the
// "thundering herd" problem.
//
// If the context is canceled before the wait is over, Wait returns the
// context's error.
func (p *RetryPolicy) Wait(ctx context.Context, attempt int) error {
	waitTime := float64(p.MinRetryDelay) * math.Pow(_backoffFactor, float64(attempt))

	if time.Duration(waitTime) > p.MaxRetryDelay {
		waitTime = float64(p.MaxRetryDelay)
	}

	jitter := (rand.Float64()*2 - 1) * _jitterFactor * waitTime //nolint:gosec // we don't need cryptographic randomness

	waitTimeWithJitter := time.Duration(waitTime + jitter)

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", ctx.Err())
	case <-time.After(waitTimeWithJitter):
		return nil
	}
}

// defaultIsRetryable defines the default logic to determine if a request should
// be retried.
func defaultIsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	switch resp.StatusCode {
	case http.StatusRequestTimeout, // 408
		http.StatusTooManyRequests,    // 429
		http.StatusBadGateway,         // 502
		http.StatusServiceUnavailable, // 503
		http.StatusGatewayTimeout:     // 504
		return true
	}

	return false
}
