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
