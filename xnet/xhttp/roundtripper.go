package xhttp

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"git.sr.ht/~jamesponddotco/xstd-go/xio"
	"git.sr.ht/~jamesponddotco/xstd-go/xlog/xslog"
)

// ErrExceededMaxRetries is returned when the maximum number of retries has
// been exceeded.
const ErrExceededMaxRetries xerrors.Error = "exceeded max retries"

// RetryRoundTripper is an implementation of http.RoundTripper that retries
// requests.
type RetryRoundTripper struct {
	// Policy defines the retry behavior.
	Policy *RetryPolicy

	// Logger is used for structured logging. It should be provided if logging
	// of the retry process is required. If nil, logging is disabled.
	Logger *slog.Logger

	// transport is the underlying http.RoundTripper.
	transport http.RoundTripper
}

// NewRetryRoundTripper creates a new instance of RetryRoundTripper with the
// specified RetryPolicy and Logger.
//
// If policy is nil, NewRetryRoundTripper will use NewRetryPolicy. If logger is
// nil, NewRetryRoundTripper will not log.
func NewRetryRoundTripper(policy *RetryPolicy, logger *slog.Logger) *RetryRoundTripper {
	if policy == nil {
		policy = NewRetryPolicy()
	}

	return &RetryRoundTripper{
		Policy:    policy,
		Logger:    logger,
		transport: NewTransport(),
	}
}

// RoundTrip executes an HTTP transaction with retry logic.
func (rt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyReader io.ReadSeeker

	if req.Body != nil {
		var err error

		bodyReader, err = xio.ReaderToReadSeeker(req.Body)
		if err != nil {
			xslog.LogEvent(
				req.Context(),
				rt.Logger,
				slog.LevelError,
				"failed to convert request body to ReadSeeker",
				slog.Any("error", err),
			)

			return nil, fmt.Errorf("%w", err)
		}
	}

	originalBody := req.Body

	defer func() {
		if originalBody != nil {
			originalBody.Close()
		}
	}()

	for attempt := 0; attempt <= rt.Policy.MaxRetries; attempt++ {
		xslog.LogEvent(
			req.Context(),
			rt.Logger,
			slog.LevelDebug,
			"attempting request",
			slog.Int("attempt", attempt),
		)

		if bodyReader != nil {
			_, err := bodyReader.Seek(0, io.SeekStart)
			if err != nil {
				xslog.LogEvent(
					req.Context(),
					rt.Logger,
					slog.LevelError,
					"failed to seek request body",
					slog.Any("error", err),
				)

				return nil, fmt.Errorf("%w", err)
			}

			req.Body = io.NopCloser(bodyReader)
		}

		resp, err := rt.transport.RoundTrip(req)

		// Successful, non-retryable case.
		if err == nil && !rt.Policy.IsRetryable(resp, nil) {
			xslog.LogEvent(
				req.Context(),
				rt.Logger,
				slog.LevelDebug,
				"request succeeded",
				slog.Int("attempt", attempt),
			)

			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		// Retryable case.
		if err != nil || rt.Policy.IsRetryable(resp, err) {
			if resp != nil {
				if drainErr := DrainResponseBody(resp); drainErr != nil {
					xslog.LogEvent(
						req.Context(),
						rt.Logger,
						slog.LevelError,
						"failed to drain response body",
						slog.Any("error", drainErr),
					)

					return nil, fmt.Errorf("%w", drainErr)
				}
			}

			if waitErr := rt.Policy.Wait(req.Context(), attempt); waitErr != nil {
				xslog.LogEvent(
					req.Context(),
					rt.Logger,
					slog.LevelError,
					"failed during retry wait",
					slog.Any("error", waitErr),
				)

				return nil, fmt.Errorf("%w", waitErr)
			}

			req.Body = originalBody

			continue
		}

		// Non-retryable error or response.
		return resp, fmt.Errorf("%w", err)
	}

	return nil, fmt.Errorf("%w", ErrExceededMaxRetries)
}
