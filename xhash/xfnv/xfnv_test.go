package xfnv_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xhash/xfnv"
)

// TestString tests the String function of the xfnv package.
func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want string
	}{
		{
			name: "Empty string",
			give: "",
			want: "cbf29ce484222325",
		},
		{
			name: "Single character",
			give: "a",
			want: "af63dc4c8601ec8c",
		},
		{
			name: "Hello, World!",
			give: "Hello, World!",
			want: "6ef05bd7cc857c54",
		},
		{
			name: "Long string",
			give: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			want: "7c44b7f23520aa23",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xfnv.String(tt.give)
			if got != tt.want {
				t.Errorf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
