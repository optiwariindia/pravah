# 📄 Product Requirements Document (PRD) - BizBandhan Pravah Website

## Project: BizBandhan Pravah Website (Customer Portal)

### 1. 🎯 Overview
The **BizBandhan Pravah Website** is the primary entry point for the managed BizBandhan Pravah service. It serves a dual purpose: a public-facing marketing and documentation site, and a secure customer portal for managing applications and API keys.

### 2. 👥 Target Audience
*   **Developers:** Looking to integrate real-time notifications via SDK.
*   **Business Owners:** Managing organizational access and monitoring usage.
*   **System Admins:** Managing site policies (Privacy, Terms) via a simple administrative interface.

### 3. 🚀 Goals & Objectives
*   **Dynamic Content Management:** Allow admins to update legal and information policies without code changes.
*   **Comprehensive Onboarding:** Capture detailed organizational data during registration for better customer support.
*   **Self-Service Dashboard:** Enable users to manage multiple event-driven applications independently.
*   **Decoupled Architecture:** Maintain a clean separation between the website and the high-performance SSE engine.

### 4. 🛠️ Features (v1.0.0 - MVP)
| Feature | Description |
| :--- | :--- |
| **Marketing Site** | Public landing page with visualizers and quick-start guides. |
| **Dynamic Policies** | System to manage and serve Privacy Policy, Terms of Service, etc., from the database. |
| **Admin Interface** | Basic CRUD for site policies and monitoring registrations. |
| **Organizational Auth** | Detailed registration (Name, Org, Size, Location, etc.) and secure login. |
| **Application Dashboard** | Manage Publisher/Subscriber applications and rotate API keys. |
| **Usage Insights** | Real-time monitoring of active connections and signal throughput. |
| **Future Scope** | Integrated billing, automated subscription management, and payment processing. |

### 5. 🏗️ High-Level Architecture
*   **Core Stack:** Node.js + Express.js + Twig.
*   **Frameworks:** `express-web-tools` for boilerplate, logging, and modular routing.
*   **Structure:** Modular directory layout (`website/app/modules/`).
*   **Configuration:** Pure environment variables (no `.env` files in production).
*   **Database:** MongoDB (Mongoose) for users, apps, and dynamic policies.

### 6. 📈 Success Metrics
*   **Self-Service Efficiency:** Users should be able to create an account and generate an API key in under 3 minutes.
*   **Content Agility:** Ability for an admin to update a policy and have it reflect live instantly.
*   **Data Integrity:** 100% reliability in API key validation across the portal and SSE engine.

---
**Lead Engineer:** Om Prakash Tiwari
**Date:** May 25, 2026
