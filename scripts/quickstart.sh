#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

export SUBSCRIPTION_JWT_SECRET="${SUBSCRIPTION_JWT_SECRET:-$(openssl rand -hex 32)}"
export ADMIN_API_KEY="${ADMIN_API_KEY:-$(openssl rand -hex 32)}"

docker compose up -d --build db pravah-go

for _ in $(seq 1 60); do
    STATUS="$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' pravah-pravah-go-1 2>/dev/null || true)"
    [ "$STATUS" = "healthy" ] && break
    sleep 1
done
if [ "${STATUS:-}" != "healthy" ]; then
    echo "Pravah did not become healthy" >&2
    docker compose logs pravah-go >&2
    exit 1
fi

IP="$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' pravah-pravah-go-1)"
PRAVAH_URL="http://${IP}:3000"
COOKIE_FILE="$(mktemp)"
STREAM_FILE="$(mktemp)"
trap 'rm -f "$COOKIE_FILE" "$STREAM_FILE"' EXIT

REGISTER_RESPONSE="$(curl -fsS -X PUT "$PRAVAH_URL/api/v1/register" \
    -H 'Content-Type: application/json' \
    -H "X-ADMIN-KEY: $ADMIN_API_KEY" \
    -d '{"name":"quickstart-app","retentionHours":24}')"
API_KEY="$(printf '%s' "$REGISTER_RESPONSE" | node -e 'let s="";process.stdin.on("data",d=>s+=d).on("end",()=>process.stdout.write(JSON.parse(s).data.app.apiKey))')"

TOKEN="$(curl -fsS -X POST "$PRAVAH_URL/api/v1/subscriptions/token" \
    -H 'Content-Type: application/json' \
    -H "X-API-KEY: $API_KEY" \
    -d '{"streams":["quickstart:demo"],"subject":"quickstart-user","expiresIn":300}' |
    node -e 'let s="";process.stdin.on("data",d=>s+=d).on("end",()=>process.stdout.write(JSON.parse(s).data.token))')"

curl -fsS -c "$COOKIE_FILE" "$PRAVAH_URL/pravah/init" >/dev/null
curl -sS -N --max-time 8 -b "$COOKIE_FILE" \
    "$PRAVAH_URL/pravah/subscribe?streams=quickstart:demo&token=$TOKEN" >"$STREAM_FILE" &
STREAM_PID=$!
sleep 1

curl -fsS -X POST "$PRAVAH_URL/pravah/publish" \
    -H 'Content-Type: application/json' \
    -H "X-API-KEY: $API_KEY" \
    -d '{"streamId":"quickstart:demo","event":"resource-updated","data":{"id":"demo-123"}}'
echo

sleep 1
grep -F '"event":"resource-updated"' "$STREAM_FILE"
kill "$STREAM_PID" 2>/dev/null || true

echo "Quickstart passed. Pravah remains available at $PRAVAH_URL on the Docker network."
