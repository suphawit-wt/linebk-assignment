package passwd

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptParams holds configuration options for Bcrypt password hashing.
type BcryptParams struct {
	Cost int // bcrypt algorithm's complexity
}

// DefaultBcryptParams returns recommended default parameters for Bcrypt.
func DefaultBcryptParams() BcryptParams {
	return BcryptParams{
		Cost: bcrypt.DefaultCost,
	}
}

// NewBcryptEncoder constructs a new Bcrypt password encoder.
func NewBcryptEncoder(params BcryptParams) Encoder {
	return &params
}

// Hash creates a Bcrypt hash of the given password string.
func (e *BcryptParams) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), e.Cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt: failed to generate hash: %w", err)
	}

	return string(hashedPassword), nil
}

// Verify compares a raw password against an existing Bcrypt hash.
func (e *BcryptParams) Verify(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}

		return false, fmt.Errorf("bcrypt: error verifying hash: %w", err)
	}

	return true, nil
}
