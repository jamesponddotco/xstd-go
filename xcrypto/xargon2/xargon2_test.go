package xargon2_test

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xargon2"
)

func TestCompareHashAndPassword(t *testing.T) {
	t.Parallel()

	var (
		testParams   = xargon2.NewParameters()
		testPepper   = []byte("somePepper")
		testPassword = "password123"
		testHash     = xargon2.GenerateFromPassword(testPassword, testParams, nil, testPepper)
	)

	tests := []struct {
		name     string
		hash     string
		password string
		pepper   []byte
		wantErr  error
	}{
		{
			name:     "Valid",
			hash:     testHash,
			password: testPassword,
			pepper:   testPepper,
			wantErr:  nil,
		},
		{
			name:     "InvalidPassword",
			hash:     testHash,
			password: "wrongPassword",
			pepper:   testPepper,
			wantErr:  xargon2.ErrInvalidPassword,
		},
		{
			name:     "InvalidHash",
			hash:     "invalid$hash",
			password: testPassword,
			pepper:   testPepper,
			wantErr:  xargon2.ErrInvalidHash,
		},
		{
			name:     "MalformedParameters",
			hash:     "$argon2id$v=19$m=bad,t=param,p=string$someSalt$someHash",
			password: testPassword,
			pepper:   testPepper,
			wantErr:  xargon2.ErrParseHashParameters,
		},
		{
			name:     "DecodeErrorSalt",
			hash:     "$argon2id$v=19$m=65536,t=3,p=4$!nvalidSalt$someHash",
			password: testPassword,
			pepper:   testPepper,
			wantErr:  xargon2.ErrDecode,
		},
		{
			name:     "DecodeErrorHash",
			hash:     "$argon2id$v=19$m=65536,t=3,p=4$someSalt$!nvalidHash",
			password: testPassword,
			pepper:   testPepper,
			wantErr:  xargon2.ErrDecode,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := xargon2.CompareHashAndPassword(tt.hash, tt.password, tt.pepper)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateFromPassword(t *testing.T) {
	t.Parallel()

	var (
		testParams = xargon2.NewParameters()
		testPepper = []byte("somePepper")
	)

	tests := []struct {
		name       string
		password   string
		parameters xargon2.Parameters
		pepper     []byte
	}{
		{
			name:       "DefaultParameters",
			password:   "password123",
			parameters: testParams,
			pepper:     testPepper,
		},
		{
			name:     "EmptyPassword",
			password: "",
			parameters: xargon2.Parameters{
				Threads: 4,
				Time:    3,
				Memory:  64 * 1024,
				KeyLen:  32,
			},
			pepper: testPepper,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash := xargon2.GenerateFromPassword(tt.password, tt.parameters, nil, tt.pepper)

			parts := strings.Split(hash, "$")
			if len(parts) != 6 {
				t.Errorf("invalid hash format: %s", hash)

				return
			}

			saltBase64, hashBase64 := parts[4], parts[5]

			// Decode salt and hash to ensure they are properly base64 encoded.
			_, err := base64.RawStdEncoding.DecodeString(saltBase64)
			if err != nil {
				t.Errorf("failed to decode salt: %v", err)
			}

			_, err = base64.RawStdEncoding.DecodeString(hashBase64)
			if err != nil {
				t.Errorf("failed to decode hash: %v", err)
			}
		})
	}
}
