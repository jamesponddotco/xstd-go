package xrand

// Shuffle randomizes the order of elements inside a string slice using the
// Fisher-Yates shuffle algorithm and a ChaCha20-based random number generator.
func Shuffle(str []string) {
	for i := len(str) - 1; i > 0; i-- {
		j := IntChaChaCha(i + 1)

		str[i], str[j] = str[j], str[i]
	}
}
