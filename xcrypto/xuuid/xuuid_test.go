package xuuid_test

import (
	"crypto/rand"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xuuid"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

func TestUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		genError error
		wantErr  error
	}{
		{
			name:     "Successful UUID generation",
			genError: nil,
			wantErr:  nil,
		},
		{
			name:     "Failed UUID generation",
			genError: xerrors.Error("rand.Reader error"),
			wantErr:  xuuid.ErrGenerateUUID,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Mock the rand.Reader if needed.
			if tt.genError != nil {
				randReaderOrig := rand.Reader
				defer func() { rand.Reader = randReaderOrig }()
				rand.Reader = errorReader{err: tt.genError}
			}

			uuid, err := xuuid.New()

			if (tt.wantErr != nil) != (err != nil) {
				t.Errorf("New() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if err == nil {
				uuidStr := uuid.String()
				if len(uuidStr) != 36 {
					t.Errorf("UUID string length incorrect: got %d, want 36", len(uuidStr))
				}
			}
		})
	}
}

type errorReader struct {
	err error
}

func (r errorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}
