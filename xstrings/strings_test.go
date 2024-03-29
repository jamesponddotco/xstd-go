package xstrings_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xstrings"
)

func TestRemove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		str    string
		set    string
		result string
	}{
		{
			name:   "EmptyString",
			str:    "",
			set:    xstrings.Letters,
			result: "",
		},
		{
			name:   "RemoveAllLetters",
			str:    "Hello123!@#",
			set:    xstrings.Letters,
			result: "123!@#",
		},
		{
			name:   "RemoveAllNumbers",
			str:    "H3ll0W0rld",
			set:    xstrings.Numbers,
			result: "HllWrld",
		},
		{
			name:   "RemoveSymbols",
			str:    "Hello@World!",
			set:    xstrings.Symbols,
			result: "HelloWorld",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := xstrings.Remove(tt.str, tt.set)
			if result != tt.result {
				t.Errorf("Expected %q, got %q", tt.result, result)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		str    []string
		result string
	}{
		{
			name:   "EmptySlice",
			str:    []string{},
			result: "",
		},
		{
			name:   "SingleString",
			str:    []string{"Hello"},
			result: "Hello",
		},
		{
			name:   "MultipleStrings",
			str:    []string{"Hello", "World", "!"},
			result: "HelloWorld!",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := xstrings.Join(tt.str...)
			if result != tt.result {
				t.Errorf("Expected %q, got %q", tt.result, result)
			}
		})
	}
}

func TestJoinWithSeparator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		sep    string
		str    []string
		result string
	}{
		{
			name:   "EmptySlice",
			sep:    "",
			str:    []string{},
			result: "",
		},
		{
			name:   "SingleString",
			sep:    "",
			str:    []string{"Hello"},
			result: "Hello",
		},
		{
			name:   "Multiple strings",
			sep:    " ",
			str:    []string{"Hello", "World", "!"},
			result: "Hello World !",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := xstrings.JoinWithSeparator(tt.sep, tt.str...)
			if result != tt.result {
				t.Errorf("Expected %q, got %q", tt.result, result)
			}
		})
	}
}
