package xhttp

import (
	"log/slog"
	"net/http"
	"time"
)

// DefaultClientTimeout is the default timeout for the http.Client.
const DefaultClientTimeout = 15 * time.Second

// NewClient returns a new http.Client given the provided timeout. Unlike Go's
// http.DefaultClient:
// - It is not shared, thus cannot be altered by other modules.
// - It doesn't follow redirects.
// - It doesn't accept cookies.
//
// A zero timeout means to use a default timeout.
func NewClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = DefaultClientTimeout
	}

	return &http.Client{
		Transport: NewTransport(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     nil,
		Timeout: timeout,
	}
}

// NewRetryingClient returns a new http.Client given the provided timeout and
// RetryPolicy. Unlike Go's http.DefaultClient:
// - It is not shared, thus cannot be altered by other modules.
// - It doesn't follow redirects.
// - It doesn't accept cookies.
// - It retries failed requests.
//
// A zero timeout means to use a default timeout. A nil policy means to use a
// default policy. A nil logger means to not log.
func NewRetryingClient(timeout time.Duration, policy *RetryPolicy, logger *slog.Logger) *http.Client {
	var (
		client    = NewClient(timeout)
		transport = NewRetryRoundTripper(policy, logger)
	)

	client.Transport = transport

	return client
}
