# 🔄 Use Cases

## UC.1: Subscribe to Event Stream
**Actor:** Frontend Client
**Pre-condition:** User is authenticated.

1.  Client sends a `GET` request to `/pravah/subscribe?streams=project:123,user:456`.
2.  Server validates the user's session/token.
3.  Server checks if the user has permission to access the requested streams.
4.  Server checks for the `Last-Event-ID` header.
5.  If `Last-Event-ID` exists, Server fetches missed events from MongoDB and sends them.
6.  Server establishes a persistent SSE connection.
7.  Server adds the client to the `ClientManager` registry.
8.  Server begins sending periodic heartbeat messages.

## UC.2: Publish an Event
**Actor:** Main Application Backend
**Pre-condition:** Service API Key is valid.

1.  Main App sends a `POST` request to `/pravah/publish` with `streamId`, `event`, and `data`.
2.  Server validates the API key.
3.  Server validates the event payload structure.
4.  Server saves the event to MongoDB (with TTL).
5.  Server passes the event to `EventBus`.
6.  `EventBus` identifies all active clients in `ClientManager` subscribed to the `streamId`.
7.  Server sends the event data to each matching client over their SSE connection.

## UC.3: Client Reconnection (Replay)
**Actor:** Frontend Client
**Pre-condition:** Client was previously connected but lost connection.

1.  Client detects disconnection.
2.  Client re-establishes connection via `GET /pravah/subscribe`.
3.  Client includes the ID of the last successfully received event in the `Last-Event-ID` header.
4.  Server identifies the client and validates access.
5.  Server queries MongoDB for all events in the subscribed streams with an ID greater than `Last-Event-ID`.
6.  Server sends these historical events to the client in order.
7.  Server resumes live event streaming.

## UC.4: Graceful Disconnection
**Actor:** Frontend Client / Server

1.  Client closes the browser tab or network drops.
2.  Server detects the closed socket.
3.  Server removes the client from the `ClientManager` registry.
4.  Server stops heartbeat and releases resources (in-memory objects).

---
**Lead Engineer:** Om Prakash Tiwari
