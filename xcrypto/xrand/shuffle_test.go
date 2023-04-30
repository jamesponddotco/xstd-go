package xrand_test

import (
	"reflect"
	"sort"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xrand"
)

func TestShuffle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  []string
	}{
		{
			name: "Empty slice",
			str:  []string{},
		},
		{
			name: "Single element",
			str:  []string{"a"},
		},
		{
			name: "Two elements",
			str:  []string{"a", "b"},
		},
		{
			name: "Three elements",
			str:  []string{"a", "b", "c"},
		},
		{
			name: "Ten elements",
			str:  []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			original := make([]string, len(tt.str)) //nolint:makezero // @TODO: make zeroed
			copy(original, tt.str)

			xrand.Shuffle(tt.str, nil)

			// Test that the length of the original and shuffled slices are equal
			if len(original) != len(tt.str) {
				t.Errorf("Shuffle() modified the length of the slice: original length = %d, shuffled length = %d", len(original), len(tt.str))
			}

			// Test that the shuffled slice still contains the same elements as the original
			sort.Strings(original)
			sortedShuffled := make([]string, len(tt.str)) //nolint:makezero // @TODO: make zeroed
			copy(sortedShuffled, tt.str)
			sort.Strings(sortedShuffled)

			if !reflect.DeepEqual(original, sortedShuffled) {
				t.Errorf("Shuffle() changed the elements in the slice: original = %v, shuffled = %v", original, tt.str)
			}
		})
	}
}
