# 👥 User Stories

## 👩‍💻 Developer (Internal Application)

*   **US.1 - Publish Event:** As a developer, I want to send an event via a REST API so that I can notify users of important system actions (e.g., respondent added).
*   **US.2 - Secure Publishing:** As a developer, I want the publish endpoint to be authenticated so that unauthorized services cannot inject fake events into my streams.
*   **US.3 - Stream Targeting:** As a developer, I want to specify a `streamId` when publishing an event so that only users interested in that specific context (e.g., a specific project) receive the notification.

## 👨‍💻 Frontend Developer (Client-side)

*   **US.4 - Real-time Subscription:** As a frontend developer, I want to subscribe to one or more streams using SSE so that my UI updates instantly without polling.
*   **US.5 - Automatic Replay:** As a frontend developer, I want the system to automatically send me events I missed during a brief disconnect (using `Last-Event-ID`) so that my dashboard stays consistent.
*   **US.6 - Connection Heartbeat:** As a frontend developer, I want to receive heartbeats from the server so that I can detect a stalled connection and initiate a reconnect.

## 👤 End User (Dashboard User)

*   **US.7 - Instant Updates:** As a dashboard user, I want to see new data appear on my screen immediately after it happens so that I can make timely decisions.
*   **US.8 - Seamless Experience:** As a dashboard user, I want my data to be up-to-date even after my laptop wakes up from sleep so that I don't miss any critical events.

---

## 🛠️ Infrastructure / DevOps

*   **US.9 - Health Check:** As a DevOps engineer, I want a health endpoint so that I can monitor the service status and active connection count in production.
*   **US.10 - Event Expiration:** As a system administrator, I want old events to be automatically deleted after a certain period (TTL) so that the database doesn't grow indefinitely.
