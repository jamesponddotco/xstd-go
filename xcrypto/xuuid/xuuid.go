package xuuid

import (
	"crypto/rand"
	"fmt"
	"io"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	ErrGenerateUUID xerrors.Error = "could not generate UUID"
	ErrNotSupported xerrors.Error = "operating system not supported"
)

// _hexDigits is a lookup table for converting a byte into a hex string.
const _hexDigits = "0123456789abcdef"

// UUID represents a 128-bit Universal Unique Identifier (UUID) as defined in
// [RFC 4122]. Specifically, it represents the random variant, UUIDv4.
//
// [RFC 4122]: https://tools.ietf.org/html/rfc4122
type UUID [16]byte

// New returns a new UUID.
func New() (UUID, error) {
	var uuid UUID

	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return UUID{}, fmt.Errorf("%w: %w", ErrGenerateUUID, err)
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

// Strings returns the string representation of the UUID.
func (u UUID) String() string {
	buf := make([]byte, 0, 36)

	for i := 0; i < 16; i++ {
		if i == 4 || i == 6 || i == 8 || i == 10 {
			buf = append(buf, '-')
		}

		buf = append(buf, _hexDigits[u[i]>>4], _hexDigits[u[i]&0x0f])
	}

	return string(buf)
}
