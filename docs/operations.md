# Pravah Operations

## Deployment

Pravah is deployed with Docker Compose. Host ports are intentionally not
published; connect through the Compose network, container IP, or a reverse
proxy.

Required secrets:

- `SUBSCRIPTION_JWT_SECRET`: at least 32 characters
- `ADMIN_API_KEY`: at least 32 characters

Production examples:

```bash
export SUBSCRIPTION_JWT_SECRET="$(openssl rand -hex 32)"
export ADMIN_API_KEY="$(openssl rand -hex 32)"
export COOKIE_SECURE=true
export COOKIE_SAMESITE=lax
export CORS_ALLOWED_ORIGINS="https://app.example.com"
export MONGODB_DATA_PATH=/srv/pravah/mongodb
docker compose up -d --build db pravah-go
```

Use `GET /pravah/health` for readiness and `GET /pravah/metrics` for in-memory
operational counters. A `503` health response means live delivery is still
available but replay/persistence is degraded.

## HTTPS and Reverse Proxy

Terminate TLS at a reverse proxy and connect it to `pravah-go:3000` on the
Compose network. Disable proxy buffering for `/pravah/subscribe`, preserve
`Set-Cookie`, and allow long-lived responses.

Set these for production:

```bash
COOKIE_SECURE=true
COOKIE_SAMESITE=lax
```

Use `COOKIE_SAMESITE=none` only with `COOKIE_SECURE=true` and an explicit
cross-site requirement.

## Backup

Back up the configured `MONGODB_DATA_PATH` while MongoDB is stopped:

```bash
docker compose stop pravah-go db
tar -C /srv/pravah -czf "pravah-mongodb-$(date +%Y%m%d%H%M%S).tar.gz" mongodb
docker compose up -d db pravah-go
```

Test restoration in a separate directory and environment before relying on a
backup. Application API keys and tenant retention configuration live in
MongoDB.

## Restore

1. Stop Pravah and MongoDB.
2. Move the current database directory aside.
3. Extract the tested backup into `MONGODB_DATA_PATH`.
4. Start MongoDB with the same major release.
5. Verify `/pravah/health` and a subscription before restarting normal traffic.

## Secret Rotation

### Subscription signing secret

Changing `SUBSCRIPTION_JWT_SECRET` immediately invalidates all existing
subscription tokens. Coordinate the rollout: deploy with the new secret and
have tenant backends issue new tokens, then reconnect browsers.

### Admin key

Rotate `ADMIN_API_KEY` by updating the secret provider and restarting Pravah.
Only application registration uses this key.

### Tenant API key

There is no API-key revocation endpoint in MVP1. Revoke by changing the app
document's `status` from `active` and updating retained deployment
configuration. Cached authorization remains valid for at most five minutes
during MongoDB outages.

## Upgrade and Rollback

1. Back up MongoDB.
2. Build and test the candidate image.
3. Stop traffic and preserve the database volume.
4. Start the candidate and verify health, registration, publish, subscribe, and
   replay.
5. Roll back the Pravah image with the same database version if needed.

Do not downgrade MongoDB across major feature compatibility versions. Pin and
upgrade MongoDB deliberately.
