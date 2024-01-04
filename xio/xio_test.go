package xio_test

import (
	"errors"
	"io"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xio"
)

type readSeeker struct {
	io.Reader
}

func (*readSeeker) Seek(_ int64, _ int) (int64, error) {
	return 0, nil
}

type reader struct {
	content string
}

func (r *reader) Read(b []byte) (int, error) {
	n := copy(b, r.content)

	r.content = r.content[n:]

	if r.content == "" {
		return n, io.EOF
	}

	return n, nil
}

type errorReader struct{}

func (*errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestReaderToReadSeeker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    io.Reader
		wantErr bool
		wantNil bool
	}{
		{
			name:    "Nil Reader",
			give:    nil,
			wantErr: true,
			wantNil: true,
		},
		{
			name: "Reader Implements io.ReadSeeker",
			give: &readSeeker{
				Reader: &reader{content: "test content"},
			},
			wantNil: false,
		},
		{
			name: "Reader Does Not Implement io.ReadSeeker",
			give: &reader{
				content: "test content",
			},
			wantNil: false,
		},
		{
			name:    "Cannot Read From Reader",
			give:    &errorReader{},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := xio.ReaderToReadSeeker(tt.give)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReaderToReadSeeker() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (got == nil) != tt.wantNil {
				t.Fatalf("ReaderToReadSeeker() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}
