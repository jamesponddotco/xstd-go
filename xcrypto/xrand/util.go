package xrand

import (
	"encoding/binary"
)

// Int63 returns a non-negative cryptographically-secure random 63-bit integer
// as an int64. It panics if crypto/rand fails to generate random bytes.
func Int63() int64 {
	return int64(binary.BigEndian.Uint64(Bytes(8)) & (1<<63 - 1))
}
