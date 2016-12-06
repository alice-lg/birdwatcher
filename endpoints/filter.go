package endpoints

import (
	"fmt"
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

func ValidateProtocolParam(value string) (string, error) {

	// Check length
	if err := ValidateLength(value, 80); err != nil {
		return "", err
	}

	// Check input
	allowed := "ID_AS:.abcdef1234567890"
	if err := ValidateCharset(value, allowed); err != nil {
		return "", err
	}

	return value, nil
}
