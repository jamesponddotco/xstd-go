package xrand

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"sync"
)

// Source64 implements the math/rand.Source64 interface as a cryptographically
// secure random number generator using AES-256 in CTR mode.
type Source64 struct {
	// Stream is the underlying AES-256 CTR stream.
	cipher.Stream

	// buf is a reusable buffer for XORKeyStream.
	buf [8]byte

	// mu protects the underlying stream.
	mu sync.Mutex
}

// Compile-time check to ensure Source64 implements the rand.Source64 interface.
var _ mrand.Source64 = (*Source64)(nil)

// NewSource64 returns a new cryptographically secure random number generator
// that satisfies the math/rand.Source64 interface.
func NewSource64() (mrand.Source64, error) {
	var keyiv [32 + 16]byte

	if _, err := rand.Read(keyiv[:]); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	block, err := aes.NewCipher(keyiv[:32])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	stream := cipher.NewCTR(block, keyiv[32:])

	return &Source64{
		Stream: stream,
	}, nil
}

// Uint64 returns a random uint64.
func (s *Source64) Uint64() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.XORKeyStream(s.buf[:], s.buf[:])

	return binary.LittleEndian.Uint64(s.buf[:])
}

// Int63 returns a non-negative random 63-bit integer as an int64.
func (s *Source64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Seed is a no-op that panics.
func (*Source64) Seed(_ int64) {
	panic("xrand: Seed() not supported")
}
