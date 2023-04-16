package xstrings

import "testing"

func TestMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		giveA int
		giveB int
		want  int
	}{
		{
			name:  "Zero",
			giveA: 0,
			giveB: 0,
			want:  0,
		},
		{
			name:  "Positive",
			giveA: 1,
			giveB: 2,
			want:  2,
		},
		{
			name:  "Negative",
			giveA: -1,
			giveB: -2,
			want:  -1,
		},
		{
			name:  "Mixed",
			giveA: 1,
			giveB: -2,
			want:  1,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := max(tt.giveA, tt.giveB)
			if got != tt.want {
				t.Errorf("max(%d, %d) = %d, want %d", tt.giveA, tt.giveB, got, tt.want)
			}
		})
	}
}
