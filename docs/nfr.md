# ⚡ Non-Functional Requirements (NFR)

## 1. 🏎️ Performance
*   **Latency:** The system must deliver events from the `/pravah/publish` endpoint to connected SSE clients in less than **1 second** (excluding network latency).
*   **Throughput:** Must handle at least **50 event publications per second** on a standard single-instance setup.
*   **Concurrency:** Must support at least **100 concurrent SSE connections** per instance without performance degradation.
*   **Resource Usage:** CPU usage should stay below 60% and memory usage below 512MB under normal load (100 connections).

## 2. 🛡️ Security
*   **Encryption:** All communication (Publish API and SSE Stream) must be over **HTTPS/TLS**.
*   **Authentication:**
    *   Publish API: Requires a secure API key or internal service token.
    *   SSE Stream: Must use session cookies or secure tokens for authentication.
*   **Authorization:** The system must verify that a user has permission to access a specific stream (e.g., `project:id`) before allowing a subscription.
*   **Input Validation:** All payloads received at the `/pravah/publish` endpoint must be validated against a predefined schema.

## 3. 🔄 Reliability & Availability
*   **Uptime:** The service should aim for **99.9% availability**.
*   **Graceful Degradation:** If the database (MongoDB) is unavailable, the service should still allow live streaming of new events, even if replay is temporarily disabled.
*   **Connection Management:** The system must gracefully handle abrupt client disconnections and clean up memory resources.
*   **Data Retention:** Events must be stored in MongoDB with a **TTL (e.g., 24 hours)** to prevent database bloat.

## 4. 🛠️ Maintainability
*   **Modularity:** The code should follow the modular structure defined in `file-structure.md` (controllers, services, routes).
*   **Logging:** Detailed logging of event flow, connection errors, and authentication failures using a standard logger (e.g., Winston/Pino).
*   **Code Quality:** Adherence to ESLint rules and a minimum of 80% test coverage for core logic (ClientManager, EventBus).

## 5. 📈 Scalability
*   **Horizontal Scaling:** While the MVP is single-instance, the architecture must be designed to support **Redis Pub/Sub** for horizontal scaling in Phase 2.
*   **Stateless Publish:** The `/pravah/publish` API must remain stateless.
*   **Stateful SSE:** Connection state is managed in-memory but should be transferable to a distributed store (like Redis) in future versions.

## 6. 🧪 Usability
*   **Developer Experience:** Clear error messages for the `/pravah/publish` API.
*   **Standard Protocol:** Use standard SSE format (`id`, `event`, `data`) to ensure compatibility with all modern browser EventSource APIs.
