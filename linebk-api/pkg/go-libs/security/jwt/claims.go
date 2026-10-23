package jwt

import "github.com/golang-jwt/jwt/v5"

// AuthClaims represents custom and standard JWT claims for authentication.
//
//nolint:tagliatelle // AuthClaims fields use snake_case by design
type AuthClaims struct {
	UserID     string   `json:"user_id"`
	Email      string   `json:"email"`
	Roles      []string `json:"roles"`
	ClientType string   `json:"client_type"`
	jwt.RegisteredClaims
}
