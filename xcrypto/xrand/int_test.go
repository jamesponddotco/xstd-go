package xrand_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xrand"
)

func TestInt63(t *testing.T) {
	t.Parallel()

	for i := 0; i < 32; i++ {
		x := xrand.Int63()

		if x < 0 || x > 9223372036854775807 {
			t.Fatalf("random number is not within the range of an int64: %d", x)
		}
	}
}

func TestIntChaChaCha(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		n           int
		shouldPanic bool
	}{
		{
			name:        "n=0",
			n:           0,
			shouldPanic: true,
		},
		{
			name: "n=1",
			n:    1,
		},
		{
			name: "n=2",
			n:    2,
		},
		{
			name: "n=100",
			n:    100,
		},
		{
			name: "n=1000",
			n:    1000,
		},
		{
			name: "n=10000",
			n:    10000,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("IntChaChaCha(%d) did not panic", tt.n)
					}
				}()
			}

			got := xrand.IntChaChaCha(tt.n)

			if got < 0 || got >= tt.n {
				t.Errorf("IntChaChaCha(%d) = %d; want in range [0, %d)", tt.n, got, tt.n)
			}
		})
	}
}
