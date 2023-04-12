package xunsafe_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

func TestBytesToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		bytes []byte
		want  string
	}{
		{
			name:  "empty byte slice",
			bytes: []byte{},
			want:  "",
		},
		{
			name:  "non-empty byte slice",
			bytes: []byte{72, 101, 108, 108, 111},
			want:  "Hello",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xunsafe.BytesToString(tt.bytes)

			if got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
		want []byte
	}{
		{
			name: "empty string",
			str:  "",
			want: []byte{},
		},
		{
			name: "non-empty string",
			str:  "Hello",
			want: []byte{72, 101, 108, 108, 111},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xunsafe.StringToBytes(tt.str)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
