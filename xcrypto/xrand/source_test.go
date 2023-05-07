package xrand_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xrand"
)

func TestNewSource64(t *testing.T) {
	t.Parallel()

	src, err := xrand.NewSource64()
	if err != nil {
		t.Fatalf("NewSource64() failed: %v", err)
	}

	if src == nil {
		t.Fatal("NewSource64() returned nil")
	}
}

func TestSource64_Int63(t *testing.T) {
	t.Parallel()

	src, err := xrand.NewSource64()
	if err != nil {
		t.Fatalf("NewCryptoSource() failed: %v", err)
	}

	prevInt63 := src.Int63()

	for i := 0; i < 100; i++ {
		currentInt63 := src.Int63()

		if currentInt63 == prevInt63 {
			t.Errorf("Int63() repeated value: %d", currentInt63)
		}

		prevInt63 = currentInt63
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	src, err := xrand.NewSource64()
	if err != nil {
		t.Fatalf("NewCryptoSource() failed: %v", err)
	}

	prevUint64 := src.Uint64()

	for i := 0; i < 100; i++ {
		currentUint64 := src.Uint64()

		if currentUint64 == prevUint64 {
			t.Errorf("Uint64() repeated value: %d", currentUint64)
		}

		prevUint64 = currentUint64
	}
}

func TestSeed(t *testing.T) {
	t.Parallel()

	src, err := xrand.NewSource64()
	if err != nil {
		t.Fatalf("NewCryptoSource() failed: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Seed() should have panicked, but didn't")
		}
	}()

	src.Seed(0)
}
