package endpoints

import (
	"fmt"
	"strings"
)

/*
 * Parameter / Request Validation
 */

// Check if the value is not longer than a given length
func ValidateLength(value string, maxLength int) error {
	if len(value) > maxLength {
		return fmt.Errorf("Provided param value is too long.")
	}
	return nil
}

func ValidateCharset(value string, alphabet string) error {
	for _, check := range value {
		ok := false
		for _, char := range alphabet {
			if char == check {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("Invalid character in param value")
		}
	}
	return nil
}

func ValidateLengthAndCharset(value string, maxLength int, alphabet string) (string, error) {
	// Check length
	if err := ValidateLength(value, maxLength); err != nil {
		return "", err
	}

	// Check input
	if err := ValidateCharset(value, alphabet); err != nil {
		return "", err
	}

	return value, nil
}

func ValidateProtocolParam(value string) (string, error) {
	return ValidateLengthAndCharset(value, 80, "ABCDEFGHIJKLMNOPQRSTUVWXYZ_:.abcdefghijklmnopqrstuvwxyz1234567890")
}

func ValidatePrefixParam(value string) (string, error) {
	value = strings.Replace(value, "m", "/", 1)
	return ValidateLengthAndCharset(value, 80, "1234567890abcdef.:/")
}

func ValidateNetMaskParam(value string) (string, error) {
	return ValidateLengthAndCharset(value, 3, "1234567890")
}
