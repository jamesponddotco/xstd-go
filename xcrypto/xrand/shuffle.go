package xrand

import (
	"crypto/rand"
	"io"
	"math/big"
)

// Shuffle randomizes the order of elements inside a string slice using the
// Fisher-Yates shuffle algorithm and a ChaCha20-based random number generator.
func Shuffle(str []string, reader io.Reader) {
	if reader == nil {
		reader = rand.Reader
	}

	for i := len(str) - 1; i > 0; i-- {
		j, err := rand.Int(reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}

		str[i], str[j.Int64()] = str[j.Int64()], str[i]
	}
}
