# 📈 Marketing Requirements Document (MRD) - Pravah

This document outlines the strategic requirements to transform **Pravah** from an internal microservice into a marketable product (SaaS or Premium Open Source).

---

## 1. 🎯 Market Positioning
**Value Proposition:** "Pravah is a lightweight, self-hosted SSE service that shares one browser connection across tabs and delivers secure resource-update signals."

**Target Audience:**
*   SaaS Platforms with complex dashboards.
*   Fintech & Real-time bidding applications.
*   Collaborative tools (Project management, CRMs).
*   DevOps teams looking for a secure, internal signaling microservice.

---

## 2. 📦 Product Packaging & Distribution
To reach developers at scale, the product must be easy to consume:
*   **Official NPM Package:** Publish `pravah-sdk` to the NPM registry for standard integration (`npm install pravah-sdk`).
*   **One-Click Cloud Deployment:** Provide "Deploy to AWS/DigitalOcean" templates and official Docker images on Docker Hub.
*   **CDN Support:** Host the SDK on a global CDN for quick `<script>` tag inclusion.

---

## 3. 🛠️ Developer Experience (DX)
The success of a developer-tool depends on the first 5 minutes of usage:
*   **Pravah Playground:** An online interactive sandbox where users can trigger signals via a web-UI and see them arrive in real-time.
*   **Multi-Framework Support:** Official wrappers for:
    *   **Frontend:** React Hooks (`usePravah`), Vue Composables, Angular Services.
    *   **Backend:** Client libraries for Python, PHP, Go, and Java.
*   **Automatic API Documentation:** Swagger/OpenAPI integration for all ingestion and management endpoints.

---

## 4. 🚀 Scalability & Enterprise Readiness
To compete with industry leaders like Pusher or Ably, Pravah must support global scale:
*   **Redis Pub/Sub Integration:** Enable horizontal scaling so events are synchronized across a cluster of Pravah instances.
*   **Advanced Rate Limiting:** Prevent system abuse by implementing per-app or per-IP throttles.
*   **Event Ordering & Reliability:** Implement delivery guarantees and robust error-handling for high-throughput environments.

---

## 5. 🖥️ Management & Observability
A professional product requires a "Command Center":
*   **Admin Dashboard:** A GUI for managing applications, revoking API keys, and monitoring live connection maps.
*   **Real-time Analytics:** Visual charts showing event throughput, latency, and active connection growth.
*   **Usage Billing (SaaS):** Logic to track message counts and peak connections for tiered pricing.

---

## 6. 🛡️ Security & Compliance
*   **End-to-End Encryption:** Options for encrypting signal payloads if metadata is sensitive.
*   **Compliance Ready:** Documentation for GDPR/SOC2 compliance regarding connection logging and data persistence.
*   **Signed Subscriptions:** Enhancing the handshake to support JWT-based stream authorization.

---

## 7. 💰 Monetization Strategy (Go-To-Market)
*   **Open-Core Model:** The SSE engine remains open-source. The Admin Dashboard, Redis Scaling, and Premium Support are paid add-ons.
*   **Managed Pravah (SaaS):** A hosted version where users pay a monthly fee based on the number of concurrent connections (e.g., Free for 100, Pro for 5000+).

---
**Lead Engineer:** Om Prakash Tiwari
**Date:** April 30, 2026
