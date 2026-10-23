package jwtkey

import (
	"context"
	"fmt"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/vault"
)

func loadKeyPair(
	ctx context.Context,
	vaultClient *vault.Client,
	path string,
) (KeyPair, error) {
	privStr, err := vaultClient.Get(ctx, "jwt", path, "private")
	if err != nil {
		return KeyPair{}, fmt.Errorf("vault: %s:private fetch error: %w", path, err)
	}
	privKey, err := parseRSAPrivateKey(privStr)
	if err != nil {
		return KeyPair{}, fmt.Errorf("vault: %s:private parse error: %w", path, err)
	}

	pubStr, err := vaultClient.Get(ctx, "jwt", path, "public")
	if err != nil {
		return KeyPair{}, fmt.Errorf("vault: %s:public fetch error: %w", path, err)
	}
	pubKey, err := parseRSAPublicKey(pubStr)
	if err != nil {
		return KeyPair{}, fmt.Errorf("vault: %s:public parse error: %w", path, err)
	}

	return KeyPair{
		Private: privKey,
		Public:  pubKey,
	}, nil
}
