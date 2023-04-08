package xrand_test

import (
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xrand"
)

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		characters  string
		length      int
		expectedLen int
	}{
		{
			name:        "Length 0",
			characters:  "",
			length:      0,
			expectedLen: 0,
		},
		{
			name:        "Length 1",
			characters:  "",
			length:      1,
			expectedLen: 1,
		},
		{
			name:        "Length 10",
			characters:  "",
			length:      10,
			expectedLen: 10,
		},
		{
			name:        "Custom Characters, Length 0",
			characters:  "abc",
			length:      0,
			expectedLen: 0,
		},
		{
			name:        "Custom Characters, Length 5",
			characters:  "abc",
			length:      5,
			expectedLen: 5,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xrand.String(tt.characters, tt.length)

			if len(got) != tt.expectedLen {
				t.Errorf("String() got length = %d, want length = %d", len(got), tt.expectedLen)
			}

			if tt.length > 0 {
				for _, char := range got {
					if tt.characters != "" && !strings.ContainsRune(tt.characters, char) {
						t.Errorf("String() generated character not in the given set: %c", char)
					}
				}
			}
		})
	}
}

func FuzzString(f *testing.F) {
	f.Fuzz(func(t *testing.T, chars string, length int) {
		s := xrand.String(chars, length)

		if length > 1 && len(s) != length {
			t.Errorf("invalid string length: %d", len(s))
		}
	})
}
