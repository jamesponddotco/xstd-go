// Package xio provides functions and utilities that extends the standard
// library with additional functionality.
package xio

import (
	"bytes"
	"fmt"
	"io"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	// ErrNilReader is returned when a nil reader is passed to
	// ReaderToReadSeeker.
	ErrNilReader xerrors.Error = "provided reader is nil"

	// ErrCopy is returned when an error occurs during the copying of data from
	// the reader.
	ErrCopy xerrors.Error = "failed to buffer reader's content"
)

// ReaderToReadSeeker checks if the provided io.Reader also implements
// io.ReadSeeker. If it does not, it reads the reader's content into memory and
// returns a bytes.Reader that implements io.ReadSeeker.
func ReaderToReadSeeker(r io.Reader) (io.ReadSeeker, error) {
	if r == nil {
		return nil, ErrNilReader
	}

	if rs, ok := r.(io.ReadSeeker); ok {
		return rs, nil
	}

	var buf bytes.Buffer

	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCopy, err)
	}

	return bytes.NewReader(buf.Bytes()), nil
}
