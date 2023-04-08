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
