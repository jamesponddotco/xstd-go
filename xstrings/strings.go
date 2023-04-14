package xstrings

import "strings"

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

	builder.Grow(len(str) - len(set))

	for _, r := range str {
		if !strings.ContainsRune(set, r) {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
