package jwtkey

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/vault"
)

// KeyPair represents a pair of RSA keys (private and public).
type KeyPair struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

// UserTokens holds the JWT keys for user access and refresh tokens.
type UserTokens struct {
	Access  KeyPair
	Refresh KeyPair
}

// Keys holds all JWT keys for both users and partners.
type Keys struct {
	User UserTokens
}

// LoadJwtKeysFromVault retrieves JWT keys from Vault and returns a Keys struct.
func LoadJwtKeysFromVault(vaultClient *vault.Client) (*Keys, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userAccess, err := loadKeyPair(ctx, vaultClient, "user-access")
	if err != nil {
		return nil, err
	}

	userRefresh, err := loadKeyPair(ctx, vaultClient, "user-refresh")
	if err != nil {
		return nil, err
	}

	return &Keys{
		User: UserTokens{
			Access:  userAccess,
			Refresh: userRefresh,
		},
	}, nil
}
