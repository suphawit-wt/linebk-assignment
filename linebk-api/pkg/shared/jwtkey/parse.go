package jwtkey

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func parseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// Try PKCS#1 format first (BEGIN RSA PRIVATE KEY)
	if priv, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return priv, nil
	}

	// Fallback to PKCS#8 format (BEGIN PRIVATE KEY)
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key (tried PKCS#1 & PKCS#8): %w", err)
	}

	rsaPriv, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key type in PKCS#8 payload is not an RSA private key")
	}

	return rsaPriv, nil
}

func parseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	// Try PKIX format (BEGIN PUBLIC KEY)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaPub, ok := pubInterface.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("key type in PKIX payload is not an RSA public key")
		}
		return rsaPub, nil
	}

	// Fallback to PKCS#1 format (BEGIN RSA PUBLIC KEY)
	rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key (tried PKIX & PKCS#1): %w", err)
	}

	return rsaPub, nil
}
