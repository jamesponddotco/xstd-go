package xstrings

import (
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xunsafe"
)

const (
	LowercaseLetters string = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letters          string = LowercaseLetters + UppercaseLetters
	Numbers          string = "0123456789"
	Symbols          string = `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`
	Space            string = " "
	AllCharacters    string = Letters + Numbers + Symbols + Space
)

// ContainsOnly returns true if the string contains only the characters in the
// given set.
func ContainsOnly(s, set string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if !strings.ContainsRune(set, r) {
			return false
		}
	}

	return true
}

// Remove returns a copy of the string with all occurrences of the given characters removed.
func Remove(str, set string) string {
	if str == "" {
		return ""
	}

	var builder strings.Builder

	builder.Grow(max(len(str)-len(set), 0))

	for _, r := range str {
		if !strings.ContainsRune(set, r) {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// Join concatenates strings to create a single string faster than
// strings.Join().
func Join(str ...string) string {
	length := len(str)

	if length == 0 {
		return ""
	}

	n := 0
	for i := 0; i < length; i++ {
		n += len(str[i])
	}

	buff := make([]byte, 0, n)
	for i := 0; i < length; i++ {
		buff = append(buff, str[i]...)
	}

	return xunsafe.BytesToString(buff)
}

// JoinWithSeparator concatenates strings to create a single string with a separator between each string.
func JoinWithSeparator(separator string, str ...string) string {
	length := len(str)

	if length == 0 {
		return ""
	}

	n := 0
	for i := 0; i < length; i++ {
		n += len(str[i])
	}

	n += len(separator) * (length - 1)

	buff := make([]byte, 0, n)
	for i := 0; i < length; i++ {
		buff = append(buff, str[i]...)

		if i < length-1 {
			buff = append(buff, separator...)
		}
	}

	return xunsafe.BytesToString(buff)
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
