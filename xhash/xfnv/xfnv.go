// Package xfnv implements the 64-bit FNV-1a non-cryptographic hash algorithm.
//
// The implementation isn't compatible with [Go's hash.Hash].
//
// [Go's hash.Hash]: https://godocs.io/hash#Hash
package xfnv

import "strconv"

const (
	// _FNVaOffset64 is the 64-bit FNV-1a offset basis.
	_FNVaOffset64 uint64 = 14695981039346656037

	// _FNVaPrime64 is the 64-bit FNV-1a prime.
	_FNVaPrime64 uint64 = 1099511628211

	// _EmptyString is the 64-bit FNV-1a hash of an empty string.
	_EmptyString = "cbf29ce484222325"
)

// String returns a new 64-bit FNV-1a hash as a string.
func String(str string) string {
	if str == "" {
		return _EmptyString
	}

	var (
		hash   = _FNVaOffset64
		keyLen = len(str)
	)

	for i := 0; i < keyLen; i++ {
		hash ^= uint64(str[i])
		hash *= _FNVaPrime64
	}

	return strconv.FormatUint(hash, 16)
}
