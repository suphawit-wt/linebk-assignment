package passwd

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// ScryptParams holds configuration settings for Scrypt password hashing.
type ScryptParams struct {
	N          int // CPU/memory cost parameter (must be a power of 2)
	r          int // Block size parameter
	p          int // Parallelization parameter
	KeyLength  int
	SaltLength int
}

// DefaultScryptParams returns recommended default parameters for Scrypt.
func DefaultScryptParams() ScryptParams {
	return ScryptParams{
		N:          16_384,
		r:          8,
		p:          1,
		KeyLength:  32,
		SaltLength: 16,
	}
}

// NewScryptEncoder creates a new Scrypt password encoder.
func NewScryptEncoder(params ScryptParams) Encoder {
	return &params
}

// Hash generates a Scrypt key derivation hash for a password.
func (e *ScryptParams) Hash(password string) (string, error) {
	salt, err := generateSalt(e.SaltLength)
	if err != nil {
		return "", err
	}

	key, err := scrypt.Key([]byte(password), salt, e.N, e.r, e.p, e.KeyLength)
	if err != nil {
		return "", err
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf("$scrypt$N=%d,r=%d,p=%d$%s$%s", e.N, e.r, e.p, b64Salt, b64Key)

	return hash, nil
}

// Verify compares a raw password against an encoded Scrypt hash.
func (e *ScryptParams) Verify(hashedPassword, password string) (bool, error) {
	params, salt, storedKey, err := decodeScryptHash(hashedPassword)
	if err != nil {
		return false, err
	}

	derivedKey, err := scrypt.Key(
		[]byte(password),
		salt,
		params.N,
		params.r,
		params.p,
		params.KeyLength,
	)
	if err != nil {
		return false, fmt.Errorf("error computing scrypt hash: %w", err)
	}

	isMatch := subtle.ConstantTimeCompare(storedKey, derivedKey) == 1

	return isMatch, nil
}

func decodeScryptHash(hash string) (params *ScryptParams, salt, key []byte, err error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 5 {
		return nil, nil, nil, errors.New("scrypt: hash is not in the correct format")
	}

	if parts[1] != "scrypt" {
		return nil, nil, nil, errors.New("scrypt: incompatible variant of scrypt")
	}

	params = &ScryptParams{}
	_, err = fmt.Sscanf(parts[2], "N=%d,r=%d,p=%d", &params.N, &params.r, &params.p)
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
