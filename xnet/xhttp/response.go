package xhttp

import (
	"fmt"
	"io"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	// ErrCannotDrainResponse is returned when a response body cannot be drained.
	ErrCannotDrainResponse xerrors.Error = "cannot drain response body"

	// ErrCannotCloseResponse is returned when a response body cannot be closed.
	ErrCannotCloseResponse xerrors.Error = "cannot close response body"
)

// DrainResponseBody reads and discards the remaining content of the response
// body until EOF, then closes it. If an error occurs while draining or closing
// the response body, an error is returned.
func DrainResponseBody(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("%w: nil response", ErrCannotDrainResponse)
	}

	_, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotDrainResponse, err)
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCloseResponse, err)
	}

	return nil
}
