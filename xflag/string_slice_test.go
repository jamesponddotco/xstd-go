package xflag_test

import (
	"flag"
	"testing"

	"git.sr.ht/~jamesponddotco/xstd-go/xflag"
)

func TestStringSlice_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give xflag.StringSlice
		want string
	}{
		{
			name: "Empty slice",
			give: xflag.StringSlice{},
			want: "",
		},
		{
			name: "Single item",
			give: xflag.StringSlice{
				"hello",
			},
			want: "hello",
		},
		{
			name: "Multiple items",
			give: xflag.StringSlice{
				"hello",
				"world",
				"test",
			},
			want: "hello, world, test",
		},
		{
			name: "Items with spaces",
			give: xflag.StringSlice{
				"hello world",
				"foo bar",
			},
			want: "hello world, foo bar",
		},
		{
			name: "Items with special characters",
			give: xflag.StringSlice{
				"hello!",
				"world?",
				"#test",
			},
			want: "hello!, world?, #test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.give.String(); got != tt.want {
				t.Errorf("StringSlice.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStringSlice_Set(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveInitial xflag.StringSlice
		giveValues  []string
		wantFull    string
	}{
		{
			name:        "Add to empty slice",
			giveInitial: xflag.StringSlice{},
			giveValues: []string{
				"test",
			},
			wantFull: "test",
		},
		{
			name:        "Add multiple giveValues",
			giveInitial: xflag.StringSlice{},
			giveValues: []string{
				"value1",
				"value2",
				"value3",
			},
			wantFull: "value1, value2, value3",
		},
		{
			name:        "Add empty string",
			giveInitial: xflag.StringSlice{},
			giveValues:  []string{""},
			wantFull:    "",
		},
		{
			name: "Add to existing slice",
			giveInitial: xflag.StringSlice{
				"existing",
			},
			giveValues: []string{
				"new",
			},
			wantFull: "existing, new",
		},
		{
			name:        "Add strings with special characters",
			giveInitial: xflag.StringSlice{},
			giveValues: []string{
				"hello!",
				"world?",
				"#test",
			},
			wantFull: "hello!, world?, #test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := tt.giveInitial
			for _, v := range tt.giveValues {
				err := s.Set(v)
				if err != nil {
					t.Errorf("StringSlice.Set(%q) unexpected error: %v", v, err)
				}
			}

			if got := s.String(); got != tt.wantFull {
				t.Errorf("After Set() operations, StringSlice.String() = %q, want %q", got, tt.wantFull)
			}
		})
	}
}

func TestStringSlice_FlagIntegration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveArgs []string
		want     string
		wantLen  int
	}{
		{
			name:     "Single value",
			giveArgs: []string{"prog", "-tags", "test"},
			want:     "test",
			wantLen:  1,
		},
		{
			name:     "Multiple values",
			giveArgs: []string{"prog", "-tags", "test1", "-tags", "test2", "-tags", "test3"},
			want:     "test1, test2, test3",
			wantLen:  3,
		},
		{
			name:     "Empty value",
			giveArgs: []string{"prog", "-tags", ""},
			want:     "",
			wantLen:  1,
		},
		{
			name:     "Values with spaces",
			giveArgs: []string{"prog", "-tags", "hello world", "-tags", "foo bar"},
			want:     "hello world, foo bar",
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				flags = flag.NewFlagSet("test", flag.ContinueOnError)
				s     xflag.StringSlice
			)

			flags.Var(&s, "tags", "a slice of strings")

			err := flags.Parse(tt.giveArgs[1:])
			if err != nil {
				t.Fatalf("failed to parse flags: %v", err)
			}

			if got := len(s); got != tt.wantLen {
				t.Errorf("len(StringSlice) = %d, want %d", got, tt.wantLen)
			}

			if got := s.String(); got != tt.want {
				t.Errorf("StringSlice.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
