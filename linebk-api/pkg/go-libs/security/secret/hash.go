package secret

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken returns the SHA-256 hex string representation of the given token.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
