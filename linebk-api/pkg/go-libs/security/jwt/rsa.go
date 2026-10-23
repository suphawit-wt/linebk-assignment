package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SignRS256 generates a signed JWT token using RSA with SHA-256.
func (s *jwtService) SignRS256(
	userID string,
	sessionID string,
	email string,
	roles []string,
	clientType string,
	ttl time.Duration,
	audience []string,
	privateKey *rsa.PrivateKey,
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

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return jwtToken.SignedString(privateKey)
}

// VerifyRS256 parses and validates an RS256 signed JWT token against an RSA public key.
func (s *jwtService) VerifyRS256(tokenStr string, publicKey *rsa.PublicKey) (*AuthClaims, error) {
	if publicKey == nil {
		return nil, errors.New("public key cannot be nil")
	}

	claims := &AuthClaims{}

	jwtToken, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
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
