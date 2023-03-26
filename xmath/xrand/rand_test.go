package xrand_test

import (
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/xstd-go/xmath/xrand"
)

// TestNewRand tests the NewRand function from the xrand package.
func TestNewRand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
	}{
		{
			name:     "test without interval",
			interval: 0,
		},
		{
			name:     "test with interval",
			interval: 10 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.interval > 0 {
				time.Sleep(tt.interval)
			}

			r1 := xrand.NewRand()
			r2 := xrand.NewRand()

			// Check if the generated numbers are different for two instances of rand.Rand
			n1, n2 := r1.Intn(1000), r2.Intn(1000)
			if n1 == n2 {
				t.Errorf("expected different random numbers, got the same: %d", n1)
			}

			// Check if the generated float64 numbers are different for two instances of rand.Rand
			f1, f2 := r1.Float64(), r2.Float64()
			if f1 == f2 {
				t.Errorf("expected different random float64 numbers, got the same: %f", f1)
			}
		})
	}
}
