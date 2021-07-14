package regexutils

import (
	"regexp"
)

// IsIntegerNumber - Verify if the string is a integer number
func IsIntegerNumber(word string) (bool, error) {
	return regexp.Match(`^\d*$`, []byte(word))
}

// IsFloatNumber - Verify if the string is a float number
func IsFloatNumber(word string) (bool, error) {
	return regexp.Match(`(^\d*)(,|\.)(\d*$)`, []byte(word))
}
