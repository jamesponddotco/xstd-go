package xhttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	// ErrMissingRequest indicates that a nil http.Request was provided to the
	// ClientIP function.
	ErrMissingRequest xerrors.Error = "missing http.Request"

	// ErrInvalidIP indicates that no valid IP address could be extracted from
	// the given http.Request.
	ErrInvalidIP xerrors.Error = "invalid IP address"
)

// ClientIP implements a [best-effort algorithm] for determining the real IP
// address of the client.
//
// Under the hood, it prioritizes the rightmost IP from the last instance of the
// X-Forwarded-For header and falls back to the RemoteAddr field of http.Request
// if the header is missing.
//
// [best-effort algorithm]: https://adam-p.ca/blog/2022/03/x-forwarded-for/
func ClientIP(req *http.Request) (string, error) {
	if req == nil {
		return "", ErrMissingRequest
	}

	xff := req.Header[XForwardedFor]

	if len(xff) == 0 {
		// Fallback to http.Request.RemoteAddr if X-Forwarded-For header is
		// missing.
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrInvalidIP, err)
		}

		return ip, nil
	}

	var (
		lastInstance = xff[len(xff)-1]
		ips          = strings.Split(lastInstance, ",")
		ip           = strings.TrimSpace(ips[len(ips)-1])
	)

	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidIP, ip)
	}

	return ip, nil
}
