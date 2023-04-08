package xrand

import (
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xstrings"
)

const (
	_indexBits = 6
	_indexMask = 1<<_indexBits - 1
	_indexMax  = 63 / _indexBits
)

// String returns a random string of the given length from the given characters.
//
// If no characters are given, [xstrings.AllCharacters] is used. If the length is less than 1, an empty string is returned.
//
// Based on [code written by @icza on StackOverflow].
//
// [xstrings.AllCharacters]: https://godocs.io/git.sr.ht/~jamesponddotco/xstd-go/xstrings#pkg-constants
// [code written by @icza on StackOverflow]: https://stackoverflow.com/a/31832326
func String(characters string, length int) string {
	if length < 1 {
		return ""
	}

	if strings.TrimSpace(characters) == "" {
		characters = xstrings.AllCharacters
	}

	var builder strings.Builder

	builder.Grow(length)

	for iter, cache, remain := length-1, Int63(), _indexMax; iter >= 0; {
		if remain == 0 {
			cache, remain = Int63(), _indexMax
		}

		if index := int(cache & _indexMask); index < len(characters) {
			builder.WriteByte(characters[index])
			iter--
		}

		cache >>= _indexBits
		remain--
	}

	return builder.String()
}
