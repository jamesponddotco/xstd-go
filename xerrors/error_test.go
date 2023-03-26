package xerrors_test

import (
	"errors"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

func TestError(t *testing.T) {
	t.Parallel()

	const (
		errSimple xerrors.Error = "test error"
	)

	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "simple",
			err:  errSimple,
			want: "test error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.err.Error(); got != tt.want {
				t.Errorf("xerror.Error(): got %v, want %v", got, tt.want)
			}

			if !errors.Is(tt.err, errSimple) {
				t.Errorf("xerror.Error(): got %v, want %v", tt.err, errSimple)
			}
		})
	}
}
