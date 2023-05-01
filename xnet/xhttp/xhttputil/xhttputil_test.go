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
		key      string
		expected string
	}{
		{
			name:     "No secret in dump",
			dump:     []byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"),
			key:      "SecretKey",
			expected: "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n",
		},
		{
			name:     "Redact secret in headers",
			dump:     []byte("GET / HTTP/1.1\r\nHost: example.com\r\nAuthorization: SecretKey\r\n\r\n"),
			key:      "SecretKey",
			expected: "GET / HTTP/1.1\r\nHost: example.com\r\nAuthorization: REDACTED\r\n\r\n",
		},
		{
			name:     "Redact secret in body",
			dump:     []byte("POST /login HTTP/1.1\r\nHost: example.com\r\n\r\n{\"username\":\"user\",\"password\":\"SecretKey\"}"),
			key:      "SecretKey",
			expected: "POST /login HTTP/1.1\r\nHost: example.com\r\n\r\n{\"username\":\"user\",\"password\":\"REDACTED\"}",
		},
		{
			name:     "Redact multiple secrets",
			dump:     []byte("GET /SecretKey HTTP/1.1\r\nHost: SecretKey.example.com\r\n\r\n"),
			key:      "SecretKey",
			expected: "GET /REDACTED HTTP/1.1\r\nHost: REDACTED.example.com\r\n\r\n",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			redacted := xhttputil.RedactSecret(tt.dump, tt.key)
			if redacted != tt.expected {
				t.Errorf("Expected redacted dump to be '%s', got '%s'", tt.expected, redacted)
			}
		})
	}
}
