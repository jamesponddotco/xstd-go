// Package xargon2 provides functions and utilities to extend
// [Go's argon2 module].
//
// [Go's argon2 module]: https://pkg.go.dev/golang.org/x/crypto/argon2
package xargon2

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xrand"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
	"golang.org/x/crypto/argon2"
)

const (
	// ErrInvalidPassword is returned when the provided password is invalid.
	ErrInvalidPassword xerrors.Error = "invalid password"

	// ErrInvalidHash is returned when the provided hash is invalid.
	ErrInvalidHash xerrors.Error = "invalid hash format"

	// ErrParseHashParameters is returned when CompareHashAndPassword fails to
	// parse the hash parameters.
	ErrParseHashParameters xerrors.Error = "failed to parse hash parameters"

	// ErrDecode is returned when CompareHashAndPassword fails to decode the
	// hash or salt.
	ErrDecode xerrors.Error = "failed to decode"
)

// DefaultSaltAndPepperLength is the default length of the salt and
// pepper used by the HashPassword method in bytes.
const DefaultSaltAndPepperLength int = 16

// Parameters holds the configuration parameters for the Argon2 hashing algorithm.
type Parameters struct {
	// Threads is the number of threads used by the Argon2 hashing algorithm for
	// computing the hash. It should be less than or equal to runtime.NumCPU().
	Threads uint8

	// Time is the computational cost factor used by the Argon2 hashing algorithm.
	Time uint32

	// Memory is the amount of memory used by the Argon2 hashing algorithm in KiB.
	Memory uint32

	// KeyLen is the length of the generated key in bytes.
	KeyLen uint32
}

// NewParameters returns a new Parameters instance with default values.
func NewParameters() Parameters {
	return Parameters{
		Threads: 4,
		Time:    3,
		Memory:  64 * 1024,
		KeyLen:  32,
	}
}

// CompareHashAndPassword compares an argon2id hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func CompareHashAndPassword(hash, password string, pepper []byte) error {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return ErrInvalidHash
	}

	var memory, time, threads uint32

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrParseHashParameters, err)
	}

	var (
		saltBase64 = parts[4]
		hashBase64 = parts[5]
	)

	salt, err := base64.RawStdEncoding.DecodeString(saltBase64)
	if err != nil {
		return fmt.Errorf("%w: salt: %w", ErrDecode, err)
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(hashBase64)
	if err != nil {
		return fmt.Errorf("%w: hash: %w", ErrDecode, err)
	}

	var (
		passwordBytes    = xunsafe.StringToBytes(password)
		combinedPassword = passwordBytes
	)

	if pepper != nil {
		combinedPassword = append(combinedPassword, pepper...)
	}

	newHash := argon2.IDKey(combinedPassword, salt, time, memory, uint8(threads), uint32(len(hashBytes)))

	if subtle.ConstantTimeCompare(hashBytes, newHash) != 1 {
		return ErrInvalidPassword
	}

	return nil
}

// GenerateFromPassword returns the argon2id hash of the given password using
// the Argon2id hashing algorithm with the given configuration, and the optional
// salt and pepper.
//
// If salt is nil, a random salt is generated using crypto/rand.
func GenerateFromPassword(password string, parameters Parameters, salt, pepper []byte) string {
	if salt == nil {
		salt = xrand.Bytes(DefaultSaltAndPepperLength)
	}

	var (
		passwordBytes    = xunsafe.StringToBytes(password)
		combinedPassword = passwordBytes
	)

	if pepper != nil {
		combinedPassword = append(combinedPassword, pepper...)
	}

	hash := argon2.IDKey(
		combinedPassword,
		salt,
		parameters.Time,
		parameters.Memory,
		parameters.Threads,
		parameters.KeyLen,
	)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		parameters.Memory, parameters.Time, parameters.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}
