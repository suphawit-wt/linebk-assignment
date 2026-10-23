#! /bin/sh

set -e

export VAULT_ADDR=http://127.0.0.1:8200

# give some time for Vault to start and be ready
sleep 3

# login with root token at $VAULT_ADDR
vault login root

# enable vault jwt kv engine
vault secrets enable -path=jwt kv-v2

# create jwt user-access key
vault kv put jwt/user-access \
  private="$(cat /vault/secrets/jwt/user/access_private.pem)" \
  public="$(cat /vault/secrets/jwt/user/access_public.pem)"

# create jwt user-refresh key
vault kv put jwt/user-refresh \
  private="$(cat /vault/secrets/jwt/user/refresh_private.pem)" \
  public="$(cat /vault/secrets/jwt/user/refresh_public.pem)"
