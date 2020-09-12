package utils

import "unicode"

// does't worth making this code more generic
func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func IsLetters(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}