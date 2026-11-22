# Troubleshooting

## `401 Unauthorized` from subscription

- Call `/pravah/init` first and include the resulting cookie.
- Issue a fresh signed token from the tenant backend.
- Check token audience, expiry, signing secret, and system clock.

## `403 Forbidden` from subscription

The token does not authorize one or more requested streams. Request a token
containing the exact streams sent to `/pravah/subscribe`.

## Cookie is not stored or sent

- Use HTTPS and `COOKIE_SECURE=true` in production.
- Set `COOKIE_DOMAIN` only when a shared parent domain is required.
- Use credentials on the client request.
- Check proxy handling of `Set-Cookie`.
- For cross-site requests, explicitly configure `SameSite=None`, secure cookies,
  and the exact allowed CORS origin.

## CORS preflight fails

Set `CORS_ALLOWED_ORIGINS` to the exact comma-separated browser origins. Wildcard
origins are not enabled with credentials.

## Events arrive repeatedly or stop after reconnect

- Preserve the SSE event ID reported as `lastEventId`.
- Do not replace event IDs with resource IDs.
- Reconnect without `Last-Event-ID` and refresh full application state after a
  `410 Replay Limit Exceeded` response.

## Health returns `503` with `degraded`

Live delivery remains available for clients already connected and previously
validated publishers. Replay and persistence are unavailable. Restore MongoDB
and wait until `/pravah/health` returns `200`.

## Client is disconnected unexpectedly

A client buffer closes when it fills. Reduce publication bursts, increase
`CLIENT_SIGNAL_BUFFER`, or improve downstream processing. Check logs for a
`buffer_full` close reason.

## Reverse proxy closes the SSE stream

Disable response buffering for `/pravah/subscribe`, increase proxy read
timeouts beyond the heartbeat interval, and preserve the `Connection` stream.
For Nginx, use `proxy_buffering off` and `proxy_read_timeout` appropriate to the
deployment.

## Cannot access port 3000 from the host

Ports are intentionally not published. Use the container IP, another container
on `pravah_default`, or place a reverse proxy in front of the service.

## Subscription tokens appear in access logs

Pravah's request logger records only URL paths, not query strings. Also disable
query-string logging in proxies and replace tokens promptly if an external log
captured one.
