// Package xhttputil provides utility functions for HTTP requests and responses.
package xhttputil

import (
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xstrings"
	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

// RedactSecret replaces the value of the specific header in the HTTP request or
// response dumps with a redacted version. This is particularly useful for
// ensuring sensitive information like access keys do not get logged or
// displayed unintentionally.
//
// The dump parameter expects the output of httputil.DumpRequest or DumpResponse
// functions.
func RedactSecret(dump []byte, header string) string {
	var (
		str         = xunsafe.BytesToString(dump)
		lines       = strings.Split(str, "\n")
		lowerHeader = strings.ToLower(header)
	)

	for i, line := range lines {
		lowerLine := strings.ToLower(line)

		if strings.HasPrefix(lowerLine, lowerHeader+":") {
			index := strings.Index(line, ":")

			if index != -1 {
				lines[i] = line[:index+1] + " REDACTED"
			}
		}
	}

	return xstrings.JoinWithSeparator("\n", lines...)
}
