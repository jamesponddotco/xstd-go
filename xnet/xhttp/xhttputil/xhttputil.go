// Package xhttputil provides utility functions for HTTP requests and responses.
package xhttputil

import (
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

// RedactSecret replaces sensitive information in request and response dumps
// with the redacted version.
func RedactSecret(dump []byte, key string) string {
	var (
		str      = xunsafe.BytesToString(dump)
		redacted = strings.ReplaceAll(str, key, "REDACTED")
	)

	return redacted
}
