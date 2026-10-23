package jwt

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SignES256 generates a signed JWT token using ECDSA with P-256 and SHA-256.
func (s *jwtService) SignES256(
	userID string,
	sessionID string,
	email string,
	roles []string,
	clientType string,
	ttl time.Duration,
	audience []string,
	privateKey *ecdsa.PrivateKey,
) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key cannot be nil")
	}

	now := time.Now().UTC()

	claims := AuthClaims{
		UserID:     userID,
		Email:      email,
		Roles:      roles,
		ClientType: clientType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        sessionID,
			Subject:   userID,
			Issuer:    s.cfg.Issuer,
			Audience:  audience,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return jwtToken.SignedString(privateKey)
}

// VerifyES256 parses and validates an ES256 signed JWT token against a public key.
func (s *jwtService) VerifyES256(tokenStr string, publicKey *ecdsa.PublicKey) (*AuthClaims, error) {
	if publicKey == nil {
		return nil, errors.New("public key cannot be nil")
	}

	claims := &AuthClaims{}

	jwtToken, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		},
		jwt.WithIssuer(s.cfg.Issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
