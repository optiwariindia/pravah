# 📄 Product Requirements Document (PRD)

## Project: Pravah (Real-time Event Notification Service)

### 1. 🎯 Overview
**Pravah** is a standalone microservice designed to provide real-time event delivery to connected clients using Server-Sent Events (SSE). It bridges the gap between the main application backend and real-time dashboards/browsers, ensuring instant notification and historical event replay upon reconnection.

### 2. 👥 Target Audience
*   **Internal Developers:** Building dashboards and real-time UI components.
*   **System Integrators:** Connecting third-party services that require real-time updates.
*   **End Users:** Consumers of dashboards who need instant updates (e.g., "respondent added", "quota reached").

### 3. 🚀 Goals & Objectives
*   **Instant Notification:** Deliver domain events to clients within 1 second.
*   **Reliable Reconnect:** Allow clients to recover events missed during short network interruptions using `Last-Event-ID`.
*   **Scalable Architecture:** Provide a foundation that can eventually handle thousands of connections via Redis pub/sub.
*   **Security:** Ensure that only authorized users can subscribe to specific event streams.

### 4. 🛠️ Features (MVP)
| Feature | Description |
| :--- | :--- |
| **Event Publishing API** | A REST endpoint (`POST /pravah/publish`) for the main app to send events. |
| **SSE Subscription** | A persistent connection endpoint (`GET /pravah/subscribe`) for clients. |
| **Stream-based Routing** | Logical grouping of events (e.g., `project:123`, `user:456`). |
| **Event Persistence** | Storing events in MongoDB with TTL (Time-To-Live) for replay. |
| **Missed Event Replay** | Automatic delivery of events missed during disconnection. |
| **Heartbeat Mechanism** | Keeping connections alive and detecting dead clients. |
| **Client Registry** | Real-time tracking of active connections in-memory. |

### 5. 🏗️ High-Level Architecture
*   **Backend:** Node.js + Express.
*   **Database:** MongoDB (for event storage).
*   **Protocol:** Server-Sent Events (SSE).
*   **Communication:** Internal HTTP for publishing, SSE for streaming to clients.

### 6. 📈 Success Metrics
*   **Latency:** >95% of events delivered in < 1 second.
*   **Availability:** 99.9% uptime for the event stream.
*   **Accuracy:** Zero dropped events for connected clients (within TTL window).
*   **Scalability:** Support 100+ concurrent connections on a single instance.

### 7. 🗓️ Roadmap (Summary)
*   **Phase 1 (MVP):** Core SSE, Event Persistence, Replay, Auth.
*   **Phase 2 (Scaling):** Redis Pub/Sub integration for horizontal scaling.
*   **Phase 3 (Enterprise):** Advanced rate limiting, Admin UI, Batching.

---
**Lead Engineer:** Om Prakash Tiwari
**Date:** April 30, 2026
