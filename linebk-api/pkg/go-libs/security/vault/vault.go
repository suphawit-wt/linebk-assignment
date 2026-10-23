package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
)

type Client struct {
	Client *api.Client
}

func NewClient(appMode string, cfg *config.VaultConfig) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if cfg.Address == "" {
		return nil, fmt.Errorf("vault setup error: address cannot be empty")
	}

	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = cfg.Address

	vaultConfig.Timeout = 5 * time.Second

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("vault client creation error: %w", err)
	}

	if cfg.Namespace != "" {
		client.SetNamespace(cfg.Namespace)
	}

	switch strings.ToLower(appMode) {
	case "prod", "production", "staging":
		role := cfg.Role

		k8sAuth, err := auth.NewKubernetesAuth(
			role,
			auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize k8s auth provider: %w", err)
		}

		secret, err := client.Auth().Login(ctx, k8sAuth)
		if err != nil {
			return nil, fmt.Errorf("failed to login to vault via k8s auth: %w", err)
		}

		if secret == nil || secret.Auth == nil {
			return nil, fmt.Errorf("vault login response contained no auth information")
		}

		client.SetToken(secret.Auth.ClientToken)

	case "dev", "development", "local":
		if cfg.Token != "" {
			client.SetToken(cfg.Token)
		}

	default:
		return nil, fmt.Errorf(
			"vault setup error: invalid env %q (must be dev, staging, or prod)",
			appMode,
		)
	}

	return &Client{Client: client}, nil
}

func (v *Client) Get(ctx context.Context, mountPath, secretPath, secretKey string) (string, error) {
	if mountPath == "" || secretPath == "" || secretKey == "" {
		return "", errors.New(
			"vault get error: mountPath, secretPath, and secretKey must not be empty",
		)
	}

	secret, err := v.Client.KVv2(mountPath).Get(ctx, secretPath)
	if err != nil {
		return "", fmt.Errorf("vault kv2 get error: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("secret not found at path %q in mount %q", secretPath, mountPath)
	}

	rawVal, ok := secret.Data[secretKey]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret at path %q", secretKey, secretPath)
	}

	valStr, ok := rawVal.(string)
	if !ok {
		return "", fmt.Errorf(
			"key %q exists at path %q but value is not a string",
			secretKey,
			secretPath,
		)
	}

	return valStr, nil
}
