package passwd

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2IDEncoder holds parameters for generating Argon2id password hashes.
type Argon2IDEncoder struct {
	Time       uint32 // number of iterations
	Memory     uint32 // amount of memory (in KiB)
	Threads    uint8  // number of parallel threads
	KeyLength  uint32
	SaltLength int
}

// DefaultArgon2IDParams returns recommended standard parameters for Argon2id hashing.
func DefaultArgon2IDParams() Argon2IDEncoder {
	return Argon2IDEncoder{
		Time:       1,
		Memory:     65_536,
		Threads:    4,
		KeyLength:  32,
		SaltLength: 16,
	}
}

// NewArgon2IDEncoder creates a new Argon2id password encoder using the provided parameters.
func NewArgon2IDEncoder(params Argon2IDEncoder) Encoder {
	return &params
}

// Hash generates an Argon2id hash for the given password.
func (e *Argon2IDEncoder) Hash(password string) (string, error) {
	salt, err := generateSalt(e.SaltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, e.Time, e.Memory, e.Threads, e.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		e.Memory,
		e.Time,
		e.Threads,
		b64Salt,
		b64Key,
	)

	return hash, nil
}

// Verify checks if the raw password matches the encoded Argon2id hash.
func (e *Argon2IDEncoder) Verify(hashedPassword, password string) (bool, error) {
	params, salt, storedKey, err := decodeArgon2idHash(hashedPassword)
	if err != nil {
		return false, err
	}

	derivedKey := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLength,
	)

	isMatch := subtle.ConstantTimeCompare(storedKey, derivedKey) == 1

	return isMatch, nil
}

func decodeArgon2idHash(hash string) (params *Argon2IDEncoder, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("argon2id: hash is not in the correct format")
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, errors.New("argon2id: incompatible variant of argon2")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, errors.New("argon2id: incompatible version of argon2")
	}

	params = &Argon2IDEncoder{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = len(salt)

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	keyLength := len(key)
	if keyLength > math.MaxUint32 {
		return nil, nil, nil, errors.New("argon2id: key length exceeds uint32 max limit")
	}
	params.KeyLength = uint32(keyLength)

	return params, salt, key, nil
}
