# Changelog

Pravah follows semantic versioning. Backend/container and SDK releases are
versioned independently; this changelog groups related release work.

## Unreleased

### Added

- Tenant-scoped signed subscription tokens.
- Per-tenant replay retention and TTL expiration.
- Event replay with `Last-Event-ID`.
- Strict signal-only publishing.
- Deterministic multi-tab SDK leader election.
- Token-refresh reconnect and explicit SDK lifecycle cleanup.
- Degraded live delivery during MongoDB outages.
- Operational metrics, health checks, and graceful shutdown.
- MongoDB integration, SDK lifecycle, security, race, and load validation.
- OpenAPI, quickstart, operations, troubleshooting, and benchmark documentation.

### Security

- Protected application registration.
- Configurable secure cookies and CORS allowlist.
- Sanitized internal API errors and token-redacted access logs.
- Bounded request, replay, stream, and client-buffer sizes.

## 1.0.0 - Initial

- Initial Go and legacy Node services.
- MongoDB event storage.
- SSE publishing and subscription endpoints.
- JavaScript and React integration experiments.

## SDK 1.1.0

- Core JavaScript SDK and React hooks.
- Multi-tab leader election using `BroadcastChannel`.
