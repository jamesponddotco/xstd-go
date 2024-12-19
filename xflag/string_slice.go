package xflag

import "git.sr.ht/~jamesponddotco/xstd-go/xstrings"

// StringSlice is a custom flag type that allows for repeated string flags.
type StringSlice []string

// String satisfies the flag.Value interface.
func (s *StringSlice) String() string {
	return xstrings.JoinWithSeparator(", ", *s...)
}

// Set satisfies the flag.Value interface.
func (s *StringSlice) Set(value string) error { //nolint:unparam // value is unused
	*s = append(*s, value)

	return nil
}
