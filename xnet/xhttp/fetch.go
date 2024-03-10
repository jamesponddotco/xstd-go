package xhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"git.sr.ht/~jamesponddotco/xstd-go/xlog/xslog"
)

const (
	// ErrNotHTTPS is returned when trying to fetch resources from a non-secure
	// URI.
	ErrNotHTTPS xerrors.Error = "uri is not secure; must use HTTPS"

	// ErrDomainNotTrusted is returned when trying to fetch resources from a
	// domain outside the list of allowed domains.
	ErrDomainNotTrusted xerrors.Error = "domain is not trusted"

	// ErrNilHTTPClient is returned when trying to fetch resources without an
	// HTTP client.
	ErrNilHTTPClient xerrors.Error = "HTTP client cannot be nil"

	// ErrRequestFailed is returned when a request fails for unknown reasons.
	ErrRequestFailed xerrors.Error = "request failed with status code"

	// ErrExceededMaxContentLength is returned when trying to fetch resources
	// bigger than the allowed MaxContentLength.
	ErrExceededMaxContentLength xerrors.Error = "response body exceeds maximum content length"
)

// Fetcher is an interface implemented by HTTP clients to fetch resources from a
// given URL.
type Fetcher interface {
	// Fetch fetches the resource at the given URL and returns the response as
	// a byte array.
	Fetch(ctx context.Context, uri string) ([]byte, error)
}

// SafeFetcher is an implementation of Fetcher that tries to provide
// secure-by-default fetching of remote resources.
type SafeFetcher struct {
	// Client is the underlying HTTP client. If nil, ErrNilHTTPClient is
	// returned.
	Client *http.Client

	// Logger is the logger to use when logging events. If nil, no events are
	// logged.
	Logger *slog.Logger

	// Header specifies additional headers to be included in each request.
	Header http.Header

	// TrustedDomain specifies domains from which fetching is permitted.
	TrustedDomain []string

	// MaxContentLength limits the size of the response body. Fetch operations
	// exceeding this limit are aborted to prevent potential resource exhaustion
	// attacks.
	MaxContentLength int64
}

// Compile-time check to ensure SafeFetcher implements the Fetcher interface.
var _ Fetcher = (*SafeFetcher)(nil)

// Fetch fetches the resource at the given URL and returns the response as a
// byte array.
func (sf *SafeFetcher) Fetch(ctx context.Context, uri string) ([]byte, error) {
	if sf.Client == nil {
		return nil, ErrNilHTTPClient
	}

	if !strings.HasPrefix(uri, "https://") {
		return nil, ErrNotHTTPS
	}

	if !sf.isDomainTrusted(ctx, uri) {
		return nil, ErrDomainNotTrusted
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	for key, values := range sf.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	xslog.LogEvent(
		ctx,
		sf.Logger,
		slog.LevelInfo,
		"fetching resource",
		slog.String("url", uri),
	)

	resp, err := sf.Client.Do(req)
	if err != nil {
		select {
		case <-req.Context().Done():
			return nil, fmt.Errorf("%w", req.Context().Err())
		default:
			return nil, fmt.Errorf("%w", err)
		}
	}

	defer func() {
		if err = DrainResponseBody(resp); err != nil {
			xslog.LogEvent(
				ctx,
				sf.Logger,
				slog.LevelError,
				"failed to drain and close response body",
				slog.Any("error", err),
			)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		xslog.LogEvent(
			ctx,
			sf.Logger,
			slog.LevelError,
			"failed to fetch resource",
			slog.String("url", uri),
			slog.Int("status_code", resp.StatusCode),
		)

		return nil, fmt.Errorf("%w: %d", ErrRequestFailed, resp.StatusCode)
	}

	if resp.ContentLength > 0 && resp.ContentLength > sf.MaxContentLength {
		xslog.LogEvent(
			ctx,
			sf.Logger,
			slog.LevelError,
			"response body exceeds maximum content length",
			slog.String("url", uri),
			slog.Int64("content_length", resp.ContentLength),
			slog.Int64("max_content_length", sf.MaxContentLength),
		)

		return nil, ErrExceededMaxContentLength
	}

	var buffer *bytes.Buffer

	if resp.ContentLength < 0 {
		buffer = bytes.NewBuffer(make([]byte, 0, sf.MaxContentLength))
	} else {
		buffer = bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	}

	_, err = io.Copy(buffer, io.LimitReader(resp.Body, sf.MaxContentLength))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	xslog.LogEvent(
		ctx,
		sf.Logger,
		slog.LevelInfo,
		"fetched resource successfully",
		slog.String("url", uri),
	)

	return buffer.Bytes(), nil
}

// isDomainTrusted checks if the given URL's domain is listed in the trusted
// domains.
func (sf *SafeFetcher) isDomainTrusted(ctx context.Context, uri string) bool {
	if len(sf.TrustedDomain) == 0 || uri == "" {
		return false
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		xslog.LogEvent(
			ctx,
			sf.Logger,
			slog.LevelError,
			"failed to parse URL",
			slog.String("url", uri),
			slog.Any("error", err),
		)

		return false
	}

	for _, domain := range sf.TrustedDomain {
		if parsed.Hostname() == domain {
			return true
		}
	}

	return false
}
