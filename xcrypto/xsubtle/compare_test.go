package xsubtle_test

import (
	"os"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xsubtle"
)

func TestConstantTimeStringEqual(t *testing.T) {
	t.Parallel()

	if os.Getenv("JOB_ID") != "" {
		t.Log("skipping expensive test in CI environment")

		t.Skip()
	}

	tests := []struct {
		name        string
		giveStringA string
		giveStringB string
		want        bool
	}{
		{
			name:        "equal strings",
			giveStringA: "secret",
			giveStringB: "secret",
			want:        true,
		},
		{
			name:        "different strings same length",
			giveStringA: "secret",
			giveStringB: "necret",
			want:        false,
		},
		{
			name:        "different lengths",
			giveStringA: "secret",
			giveStringB: "secrets",
			want:        false,
		},
		{
			name:        "very large strings equal",
			giveStringA: strings.Repeat("a", 1<<30),
			giveStringB: strings.Repeat("a", 1<<30),
			want:        true,
		},
		{
			name:        "very large strings different",
			giveStringA: strings.Repeat("a", 1<<30),
			giveStringB: strings.Repeat("b", 1<<30),
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xsubtle.ConstantTimeStringEqual(tt.giveStringA, tt.giveStringB)
			if got != tt.want {
				t.Errorf("ConstantTimeStringEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
