package xrand

import (
	"hash/maphash"
	"math/rand"
	"time"
)

// NewRand returns a seeded rand.Rand instance. The seed is based on the current
// time and a hash computed from runtime.fastrand64()'s byte sequence.
//
// The returned rand.Rand instance provides pseudo-random numbers and is not
// suitable for security-sensitive work. See [crypto/rand] for a
// cryptographically secure random number generator.
//
// [crypto/rand]: https://golang.org/pkg/crypto/rand/
func NewRand() *rand.Rand {
	var (
		hash   maphash.Hash
		seed   = time.Now().UnixNano() + int64(hash.Sum64())
		source = rand.NewSource(seed)
	)

	return rand.New(source) //nolint:gosec // we don't need crypto/rand here
}
