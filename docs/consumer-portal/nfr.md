# 🛡️ Non-Functional Requirements (NFR) - BizBandhan Pravah Website

## 1. Performance
- **NFR-1.1:** Average page load time for public marketing pages should be < 1.0 second.
- **NFR-1.2:** Database queries for API Keys and Policy retrieval should respond in < 50ms.

## 2. Security
- **NFR-2.1:** Password storage using `bcrypt` (10+ rounds).
- **NFR-2.2:** Environment variable based configuration (no local `.env` storage).
- **NFR-2.3:** Protection against common web vulnerabilities (XSS, CSRF) handled via `express-web-tools` middleware.
- **NFR-2.4:** Secure session management for the customer dashboard.

## 3. Reliability & Availability
- **NFR-3.1:** 99.5% uptime for the website and dashboard.
- **NFR-3.2:** Stateless application logic to allow horizontal scaling via Docker clusters.

## 4. Portability & Deployment
- **NFR-4.1:** Full Dockerization of the `website/` directory.
- **NFR-4.2:** Compatibility with MongoDB v7.0+.

## 5. Maintainability
- **NFR-5.1:** Modular codebase using the specified module pattern (`app/modules`).
- **NFR-5.2:** Use of `express-web-tools` to ensure consistent patterns across controllers and models.

## 6. Compliance
- **NFR-6.1:** Dynamic policy delivery to ensure legal agility.
- **NFR-6.2:** Clear data isolation policy ("Signal-Only") highlighted on the Privacy page.

---
**Lead Engineer:** Om Prakash Tiwari
**Date:** May 25, 2026
