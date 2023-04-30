package xurl_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xurl"
)

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "basic normalization",
			input:     "HTTP://www.Example.com:80/Test/",
			expected:  "http://www.example.com/Test",
			expectErr: false,
		},
		{
			name:      "https default port",
			input:     "https://www.example.com:443/test/",
			expected:  "https://www.example.com/test",
			expectErr: false,
		},
		{
			name:      "invalid url",
			input:     "http://192.168.0.%31/",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "sort query parameters",
			input:     "https://example.com/?c=3&b=2&a=1",
			expected:  "https://example.com/?a=1&b=2&c=3",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := xurl.Normalize(tt.input)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Normalize() expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Normalize() unexpected error: %v", err)
				} else if result != tt.expected {
					t.Errorf("Normalize() got = %v, want = %v", result, tt.expected)
				}
			}
		})
	}
}
