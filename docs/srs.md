# 📘 Software Requirements Specification (SRS)

## Pravah: Real-time Event Notification Service (SSE-based)

---

# 1. 📌 Introduction

## 1.1 Purpose

This system, named **Pravah**, provides **real-time event delivery** to connected clients (dashboards, browsers) using Server-Sent Events (SSE).

It enables the main application to:

* publish domain events (e.g., respondent added, quota reached)
* notify relevant users instantly
* allow clients to recover missed events after reconnect

---

## 1.2 Scope

The service acts as a **standalone microservice** responsible for:

* managing SSE connections
* handling client subscriptions to event streams
* delivering real-time events
* storing events temporarily for replay
* enforcing access control

---

## 1.3 Definitions

| Term   | Meaning                                         |
| ------ | ----------------------------------------------- |
| Event  | A structured message describing a system action |
| Stream | Logical channel grouping related events         |
| Client | Connected consumer (browser/dashboard)          |
| Replay | Sending missed events after reconnect           |
| SSE    | Server-Sent Events protocol                     |

---

# 2. 🏗️ System Overview

## 2.1 Architecture

```
[Main Application]
        |
        |  HTTP (publish events)
        ↓
[Pravah Service]
        |
        |  SSE
        ↓
[Clients]
```

---

## 2.2 Design Principles

* simplicity over premature scaling
* stream-based routing
* stateless publish API
* stateful connection handling
* eventual scalability support

---

# 3. ⚙️ Functional Requirements

---

## 3.1 Event Publishing

### Description

The service shall accept events from the main application.

### Endpoint

```
POST /pravah/publish
```

### Requirements

* must accept structured event payload
* must validate authentication
* must persist event in database
* must dispatch event to relevant clients

---

## 3.2 Event Subscription (SSE)

### Description

Clients connect via SSE and subscribe to one or more streams.

### Endpoint

```
GET /pravah/subscribe?streams=<comma-separated>
```

### Requirements

* must establish SSE connection
* must support multiple stream subscriptions
* must validate user access per stream
* must maintain persistent connection
* must send heartbeat messages

---

## 3.3 Stream Management

### Description

Streams represent logical event channels.

### Requirements

* system shall support dynamic stream identifiers
* streams shall follow naming convention:

  * `project:<id>`
  * `user:<id>`
* clients shall subscribe to one or more streams
* events shall be routed based on streamId

---

## 3.4 Event Delivery

### Requirements

* events must be delivered in real-time
* only relevant clients shall receive events
* events shall follow SSE format:

  * id
  * event
  * data

---

## 3.5 Event Replay

### Description

Clients shall receive missed events after reconnect.

### Requirements

* system shall support `Last-Event-ID` header
* system shall fetch events newer than last ID
* system shall replay events before live streaming

---

## 3.6 Event Persistence

### Requirements

* events must be stored in database
* events shall include expiration time
* expired events shall be automatically deleted

---

## 3.7 Client Management

### Requirements

* system shall maintain active client registry (in-memory)
* each client shall include:

  * unique id
  * userId
  * subscribed streams
  * connection reference
* system shall remove clients on disconnect

---

## 3.8 Authentication & Authorization

### Requirements

* all endpoints must require authentication
* SSE authentication via cookies
* system shall validate:

  * user identity
  * stream access permissions

---

## 3.9 Health Monitoring

### Endpoint

```
GET /pravah/health
```

### Requirements

* must return service status
* must indicate active connection count

---

# 4. 📦 Data Requirements

---

## 4.1 Event Schema

```json
{
  "_id": "ObjectId",
  "streamId": "project:123",
  "event": "respondent-added",
  "data": {},
  "createdAt": "Date",
  "expiresAt": "Date"
}
```

---

## 4.2 Client Model (in-memory)

```js
{
  id: uuid,
  userId: string,
  streams: Set<string>,
  connectedAt: Date,
  res: object
}
```

---

# 5. 🔄 System Workflows

---

## 5.1 Subscription Flow

1. Client sends request with streams
2. Server validates authentication
3. Server validates stream access
4. SSE connection established
5. Replay missed events (if any)
6. Start heartbeat
7. Add client to registry

---

## 5.2 Event Flow

1. Main app sends event
2. Event stored in database
3. Event dispatched to matching clients
4. Clients receive event in real-time

---

## 5.3 Reconnection Flow

1. Client reconnects
2. Sends `Last-Event-ID`
3. Server fetches missed events
4. Server replays events
5. Resume live stream

---

# 6. ⚡ Non-Functional Requirements

---

## 6.1 Performance

* support at least 100 concurrent connections
* event delivery latency < 1 second

---

## 6.2 Reliability

* must handle client disconnect gracefully
* must avoid memory leaks
* must ensure no duplicate connections

---

## 6.3 Security

* must use HTTPS
* must avoid token exposure in URLs
* must validate all inputs

---

## 6.4 Maintainability

* modular architecture
* clear separation of concerns
* configurable event TTL

---

## 6.5 Scalability (Future)

* design shall allow integration with pub/sub system
* system shall support horizontal scaling in future

---

# 7. 🚫 Out of Scope (for current version)

* multi-instance synchronization
* Redis pub/sub
* WebSocket support
* advanced rate limiting
* event ordering guarantees across instances

---

# 8. 🧱 Technology Stack

* Node.js (runtime)
* Express (framework)
* MongoDB (event storage with TTL)

---

# 9. 🎯 Assumptions

* system will run as a single instance
* concurrent users < 100
* events are lightweight JSON payloads

---

# 10. 📌 Future Enhancements

* Redis-based scaling
* event batching
* admin dashboard
* metrics & monitoring
* stream management APIs

---

# ✅ Final Note

This SRS defines a system that is:

* **simple enough to build quickly**
* **robust enough for production (your scale)**
* **extensible for future scaling**

---
