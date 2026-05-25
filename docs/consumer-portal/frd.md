# ⚙️ Functional Requirements Document (FRD) - BizBandhan Pravah Website

## 1. Visitor Module
### 1.1 Landing Page
- **F-1.1.1:** Render the main page using Twig layouts.
- **F-1.1.2:** Display features, benefits, and "Signal Flow" visualizer.
- **F-1.1.3:** Provide a `/pricing` page (informational/static for MVP).

### 1.2 Policy System
- **F-1.2.1:** Serve Privacy Policy and Terms of Service dynamically from the `Policy` model.
- **F-1.2.2:** Simple Admin interface for CRUD operations on policies at `/admin/policies`.

## 2. Auth Module
### 2.1 Registration
- **F-2.1.1:** Capture Name, Email, Contact Number, and Designation.
- **F-2.1.2:** Capture Organization Name, Nature of Business, Company Size, and Location (Country, State, City).
- **F-2.1.3:** Enforce password hashing and basic validation.

### 2.2 Organization Management
- **F-2.2.1:** Allow primary account to manage other sub-accounts within the same organization.

## 3. Dashboard Module
### 3.1 Application Management
- **F-3.1.1:** Create and list Applications (Publisher/Subscriber roles).
- **F-3.1.2:** View and rotate API Keys for specific applications.

### 3.2 Monitoring
- **F-3.2.1:** Display active SSE connection counts.
- **F-3.2.2:** Display signal volume/throughput metrics.

## 4. Subscription Module (Future Scope)
- **F-4.1.1:** System for selecting and paying for subscription plans (Suppressed in MVP).

## 5. Technical Requirements
- **F-5.1.1:** Utilize `express-web-tools` for server, models, and controller patterns.
- **F-5.1.2:** Use Twig `extends` and `blocks` for all views.
- **F-5.1.3:** Configuration via environment variables only (Docker-compatible).

---
**Lead Engineer:** Om Prakash Tiwari
**Date:** May 25, 2026
