package xsubtle

import (
	"crypto/subtle"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

// ConstantTimeStringEqual performs a constant-time comparison of two strings to
// prevent timing attacks. It returns true if the strings are equal, false
// otherwise. The comparison time depends only on the length of the strings and
// not their contents.
//
// This function is intended for comparing sensitive strings like passwords or
// authentication tokens. For regular string comparison, use the == operator.
func ConstantTimeStringEqual(given, actual string) bool {
	var (
		givenLen    = uint64(len(given))
		actualLen   = uint64(len(actual))
		givenBytes  = xunsafe.StringToBytes(given)
		actualBytes = xunsafe.StringToBytes(actual)
		equal       = ((givenLen ^ actualLen) - 1) >> 63
	)

	if equal == 1 {
		return subtle.ConstantTimeCompare(givenBytes, actualBytes) == 1
	}

	return subtle.ConstantTimeCompare(actualBytes, actualBytes) == 1 && false //nolint:revive // this is intentional
}
