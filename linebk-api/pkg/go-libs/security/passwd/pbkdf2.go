package passwd

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// PBKDF2Params holds settings for PBKDF2 key derivation.
type PBKDF2Params struct {
	Iterations int // number of iterations
	KeyLength  int
	SaltLength int
}

// DefaultPBKDF2Params returns recommended standard parameters for PBKDF2.
func DefaultPBKDF2Params() PBKDF2Params {
	return PBKDF2Params{
		Iterations: 600_000,
		KeyLength:  32,
		SaltLength: 16,
	}
}

// NewPBKDF2Encoder creates a new PBKDF2 password encoder with specified parameters.
func NewPBKDF2Encoder(params PBKDF2Params) Encoder {
	return &params
}

// Hash generates a PBKDF2 hash of the raw password string.
func (e *PBKDF2Params) Hash(password string) (string, error) {
	salt, err := generateSalt(e.SaltLength)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(password), salt, e.Iterations, e.KeyLength, sha256.New)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf("$pbkdf2$iter=%d$%s$%s", e.Iterations, b64Salt, b64Key)

	return hash, nil
}

// Verify checks if the given password matches an encoded PBKDF2 hash.
func (e *PBKDF2Params) Verify(hashedPassword, password string) (bool, error) {
	params, salt, storedKey, err := decodePBKDF2Hash(hashedPassword)
	if err != nil {
		return false, err
	}

	derivedKey := pbkdf2.Key(
		[]byte(password),
		salt,
		params.Iterations,
		params.KeyLength,
		sha256.New,
	)

	isMatch := subtle.ConstantTimeCompare(storedKey, derivedKey) == 1

	return isMatch, nil
}

func decodePBKDF2Hash(hash string) (params *PBKDF2Params, salt, key []byte, err error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 5 {
		return nil, nil, nil, errors.New("pbkdf2: hash is not in the correct format")
	}

	if parts[1] != "pbkdf2" {
		return nil, nil, nil, errors.New("pbkdf2: incompatible variant of pbkdf2")
	}

	params = &PBKDF2Params{}
	_, err = fmt.Sscanf(parts[2], "iter=%d", &params.Iterations)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(parts[3])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = len(salt)

	key, err = base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = len(key)

	return params, salt, key, nil
}
