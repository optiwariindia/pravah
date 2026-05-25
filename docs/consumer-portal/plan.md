# Implementation Plan: BizBandhan Pravah Consumer Portal

## 1. Background & Motivation
BizBandhan Pravah is transitioning from a self-hosted open-source engine to a Managed SaaS offering. The current static HTML landing page is insufficient. We need a robust, dynamic web application to handle public marketing, legal documentation, user registration, subscription management (billing), and a customer dashboard for managing API keys. 

## 2. Scope & Impact
- **Location:** The existing `landing/` directory will be transformed into a full Express.js application.
- **Impact:** This introduces a new microservice to the stack (the "Portal"), isolating customer management from the high-performance core SSE engine (`pravah-go`).
- **Tech Stack:** Node.js, Express.js, Twig (Templating Engine), MongoDB (Mongoose).

## 3. Proposed Architecture

### 3.1 Directory Structure (`landing/`)
```text
website/
├── Dockerfile
└── app/               
    ├── package.json
    ├── server.js 
    ├── modules/  
    |   ├── visitor/   # For public pages, marketing content & legal documentation
    |   |   ├── models
    |   |   ├── controllers
    |   |   └── routes
    |   ├── auth/
    |   |   ├── models
    |   |   ├── controllers
    |   |   └── routes
    |   ├── dashboard/
    |   |   ├── models
    |   |   ├── controllers
    |   |   └── routes
    |   └── subscription/
    |       ├── models
    |       ├── controllers
    |       └── routes
    ├── views/
    |   ├── layouts/
    |   └── pages/
    ├── public/
    │   ├── css/
    │   ├── js/
    │   └── images/
    └── app/               
```

### 3.2 Integration Points
- **Database:** The portal will connect to the same MongoDB instance used by the core engine. It will create and manage records in the `apps` collection (generating `appId` and `apiKey`) on behalf of users.
- **Docker:** A new service (`portal`) will be added to `docker-compose.yml`.

## 4. Documentation & Site Map Requirements

The project documentation is split into several focused files:
- [Product Requirements Document (PRD)](./prd.md)
- [Functional Requirements Document (FRD)](./frd.md)
- [Non-Functional Requirements (NFR)](./nfr.md)
- [User Stories](./user_stories.md)

Before coding, we must define the user journey and necessary pages:

### 4.1 Public Pages (Marketing & Legal)
- **`/` (Home):** The main landing page with visualizer and features.
- **`/pricing`:** Subscription tiers and plan comparisons.
- **`/docs`:** API documentation, SDK integration guides, and proxy configuration.
- **`/privacy-policy`:** Crucial for explaining the "Signal-Only" data isolation policy.
- **`/terms-of-service`:** Standard legal terms and usage limits.
- **`/contact`:** Support and sales inquiries.

### 4.2 Authentication Flow
- **`/login`:** User login.
- **`/register`:** New account creation.
- **`/forgot-password`:** Password recovery.

### 4.3 Customer Dashboard (Protected routes `/dashboard/*`)
- **`/dashboard`:** Overview of account usage, active connections, and recent signals.
- **`/dashboard/apps`:** Create and manage applications, view and regenerate API keys.
- **`/dashboard/billing`:** Manage subscriptions, view invoices, update payment methods.
- **`/dashboard/settings`:** User profile and security settings.

## 5. Phased Implementation Plan

### Phase 1: Foundation & Public Pages
1. Initialize `package.json` with `express`, `twig`, `mongoose`. Do not use dotenv for configuration, instead, rely on environment variables passed via Docker / Docker Compose.
2. Setup `server.js` with static file serving and Twig engine. Use express-web-tools for predefined utilities logging, server boilerplate, models, controllers, and routes. for details about the library, refer to https://npmjs.com/package/express-web-tools and https://github.com/optiwariindia/express-web-tools.
3. define common css files in `public/css/` and common js files in `public/js/` for reuse across all pages.
4. Create boilerplate Twig templates for legal requirements: `policy.twig`, and create model named policy with fields name, content, createdAt, updatedAt. Create a controller and route to serve these policies at `/privacy-policy` and `/terms-of-service`.
5. Create an admin interface to manage these policies (CRUD operations) at `/admin/policies` (this will be protected by a simple admin authentication, which can be hardcoded for now). view registrations and statistics of the portal can also be added to this admin interface in future iterations.

### Phase 2: Authentication System

- We need to manage user accounts to provide personalized dashboards and secure access to API key management and billing information.
- We consider that client who wants to use our services will let us know about a contact person and basic details about their organization. This contact person will be the one who will register on our portal and manage the account. Hence, we can have a simple registration form that captures the name, email, contact number, designation, password, and organization name, nature of business, company size, and location (country, state, city). This will help us in better understanding our customers and providing them with personalized support and services.
- We can allow multiple users from the same organization but the primary account will have a right to manage all other accounts and view billing information. This will help in better account management and security.

### Phase 3: Customer Dashboard
- The dashboard will allow users to manage their applications (both, publisher and subscriber) and API keys. Users can create multiple applications, each with its own API key. This allows for better organization and security, as users can regenerate keys for specific applications without affecting others.
- The dashboard will also provide insights into API usage, active connections, and recent signals. This will help users understand their usage patterns and optimize their applications accordingly.
- The dashboard will also have a billing section where users can view their subscription details, invoices, and manage their payment methods. This will provide a seamless experience for users to manage their accounts and subscriptions.
- The dashboard will be designed with a user-friendly interface, making it easy for users to navigate and manage their accounts effectively.

### Phase 4: Billing Scaffolding
- A subscription management system will handle billing. A feature based subscription model or volume based subscription model can be implemented based on the pravah usage patterns and customer feedback. In MVP this part shall be suppressed. (this part shall be in future scope of the product.)

## 5. Verification
- Verify the landing page loads correctly with Twig.
- Verify user registration inserts correctly into the database.
- Verify that API keys generated via the new portal dashboard are recognized and valid when publishing to the `pravah-go` engine.
