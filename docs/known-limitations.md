# MVP1 Known Limitations

- Pravah runs as a single stateful instance. In-memory client connections are not
  synchronized across multiple service replicas.
- MongoDB stores replay events, but a multi-instance Redis/NATS transport is not
  implemented.
- Delivery guarantees cover connected clients and retained replay only. There is
  no end-to-end consumer acknowledgement protocol.
- Events accepted during a database outage are delivered live only and cannot be
  replayed.
- Subscription tokens are JWTs with short expiration but no revocation list.
- Publisher authentication uses an application API key; there is no UI for key
  rotation or application administration.
- Stream authorization is enforced by token claims. Pravah does not independently
  understand application-specific roles or object ownership.
- The SDK relies on modern browser `EventSource`, `BroadcastChannel`, and `fetch`
  support. It does not provide a polyfill.
- The in-memory metrics endpoint resets on restart and is not a long-term
  monitoring backend.
- Events are signal-only. Applications remain responsible for securely fetching
  resource data after receiving an identifier.
- Hosted SaaS, admin UI, billing, horizontal scaling, delivery batching, and
  formal event-order guarantees across instances are Phase 2 work.
