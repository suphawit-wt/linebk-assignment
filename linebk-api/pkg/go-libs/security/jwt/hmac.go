package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SignHS256 generates a signed JWT token using HMAC SHA-256 and a secret key.
func (s *jwtService) SignHS256(
	userID string,
	sessionID string,
	email string,
	roles []string,
	clientType string,
	ttl time.Duration,
	audience []string,
	secretKey string,
) (string, error) {
	if secretKey == "" {
		return "", errors.New("secret key cannot be empty")
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

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secretKey))
}

// VerifyHS256 parses and validates an HS256 signed JWT token using a secret key.
func (s *jwtService) VerifyHS256(tokenStr, secretKey string) (*AuthClaims, error) {
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}

	claims := &AuthClaims{}

	jwtToken, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
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
