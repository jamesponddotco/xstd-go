package xhttputil_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp/xhttputil"
)

func TestRedactSecret(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		dump     []byte
		header   string
		expected string
	}{
		{
			name:     "Simple redaction",
			dump:     []byte("Authorization: Bearer mysecrettoken\nContent-Type: application/json\n"),
			header:   "Authorization",
			expected: "Authorization: REDACTED\nContent-Type: application/json\n",
		},
		{
			name:     "Case insensitive header",
			dump:     []byte("authorization: Bearer mysecrettoken\nContent-Type: application/json\n"),
			header:   "Authorization",
			expected: "authorization: REDACTED\nContent-Type: application/json\n",
		},
		{
			name:     "No matching header",
			dump:     []byte("Content-Type: application/json\n"),
			header:   "Authorization",
			expected: "Content-Type: application/json\n",
		},
		{
			name:     "Empty input",
			dump:     []byte(""),
			header:   "Authorization",
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := xhttputil.RedactSecret(tt.dump, tt.header)
			if result != tt.expected {
				t.Errorf("RedactSecret(%q, %q) = %q; want %q", tt.dump, tt.header, result, tt.expected)
			}
		})
	}
}
