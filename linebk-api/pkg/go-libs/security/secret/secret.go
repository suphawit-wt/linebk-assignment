package secret

import (
	"crypto/rand"
	"errors"
	"strings"
)

const (
	charsetAlphaNumericLower = "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetAlphaNumericUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetDigits            = "0123456789"

	letterIdxBits = 6                    // 6 bits to represent 64 characters
	letterIdxMask = 1<<letterIdxBits - 1 // All 1's, as many as `letterIdxBits`
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// NewDigits generates a random numeric secret string of the specified length.
func NewDigits(length int) (string, error) {
	return generate(length, charsetDigits)
}

// NewAlphaNumericLower generates a random lowercase alphanumeric secret string of the specified length.
func NewAlphaNumericLower(length int) (string, error) {
	return generate(length, charsetAlphaNumericLower)
}

// NewAlphaNumericUpper generates a random uppercase alphanumeric secret string of the specified length.
func NewAlphaNumericUpper(length int) (string, error) {
	return generate(length, charsetAlphaNumericUpper)
}

func generate(length int, charset string) (string, error) {
	if length < 4 || length > 12 {
		return "", errors.New("secret: length must be between 4 and 12")
	}

	// Create a byte slice to hold enough random noise
	// We need 'length' amount of indices
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	sb := strings.Builder{}
	sb.Grow(length)

	for _, b := range bytes {
		// Use the byte to pick a character from the charset
		// We use modulo here for simplicity with crypto bytes
		idx := int(b) % len(charset)
		sb.WriteByte(charset[idx])
	}

	return sb.String(), nil
}
