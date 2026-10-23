package passwd

import (
	"crypto/rand"
	"fmt"
)

// Encoder defines the interface for password hashing and verification strategies.
type Encoder interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, password string) (bool, error)
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}
