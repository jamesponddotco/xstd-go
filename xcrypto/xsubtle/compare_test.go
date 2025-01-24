package xsubtle_test

import (
	"os"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xsubtle"
)

func TestConstantTimeStringEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveStringA string
		giveStringB string
		skipCI      bool
		want        bool
	}{
		{
			name:        "equal strings",
			giveStringA: "secret",
			giveStringB: "secret",
			skipCI:      false,
			want:        true,
		},
		{
			name:        "different strings same length",
			giveStringA: "secret",
			giveStringB: "necret",
			skipCI:      false,
			want:        false,
		},
		{
			name:        "different lengths",
			giveStringA: "secret",
			giveStringB: "secrets",
			skipCI:      false,
			want:        false,
		},
		{
			name:        "very large strings equal",
			giveStringA: strings.Repeat("a", 1<<30),
			giveStringB: strings.Repeat("a", 1<<30),
			skipCI:      true,
			want:        true,
		},
		{
			name:        "very large strings different",
			giveStringA: strings.Repeat("a", 1<<30),
			giveStringB: strings.Repeat("b", 1<<30),
			skipCI:      true,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.skipCI && os.Getenv("JOB_ID") != "" {
				t.Skip()
			}

			got := xsubtle.ConstantTimeStringEqual(tt.giveStringA, tt.giveStringB)
			if got != tt.want {
				t.Errorf("ConstantTimeStringEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
