package xrand

import (
	"crypto/rand"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/chacha20"
)

// Int63 returns a non-negative cryptographically-secure random 63-bit integer
// as an int64. It panics if crypto/rand fails to generate random bytes.
func Int63() int64 {
	return int64(binary.BigEndian.Uint64(Bytes(8)) & (1<<63 - 1))
}

// IntChaChaCha returns a cryptographically secure random integer in the range
// [0, n) using the ChaCha20 stream cipher. If the provided io.Reader is nil,
// it defaults to using crypto/rand.Reader.
func IntChaChaCha(n int, reader io.Reader) int {
	if reader == nil {
		reader = rand.Reader
	}

	var (
		key   = make([]byte, 32)
		nonce = make([]byte, 12)
	)

	_, err := io.ReadFull(reader, key)
	if err != nil {
		panic(err)
	}

	_, err = io.ReadFull(reader, nonce)
	if err != nil {
		panic(err)
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		panic(err)
	}

	randomBytes := make([]byte, 4) //nolint:makezero // By design.

	cipher.XORKeyStream(randomBytes, randomBytes)

	randomInt := int(binary.BigEndian.Uint32(randomBytes))

	return randomInt % n
}
