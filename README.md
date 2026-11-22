# 🌊 Pravah

**Pravah** (meaning *Flow* or *Stream*) is a high-performance, secure, and project-agnostic real-time event notification microservice based on Server-Sent Events (SSE).

Now powered by a **Go engine**, Pravah provides a minimal, self-hosted SSE runtime for reliable signal delivery.

## 🚀 The Problem
Modern web applications face major challenges when implementing real-time features:
1.  **Server Overload:** Each open tab establishes its own persistent connection.
2.  **Data Security:** Broadcasting sensitive data over event streams increases the risk of leaks.
3.  **Complexity:** Managing identity and multi-tab synchronization is repetitive and error-prone.

## 💡 Our Solution
Pravah solves these issues through a specialized architecture:
-   **Multi-Tab Synchronization:** A smart SDK that uses Leader Election to maintain only **ONE** active SSE connection per browser.
-   **Signal-Only Policy:** Pravah broadcasts minimalist "Signals" (Event Name + Resource ID). Sensitive data is fetched via your existing secure REST APIs.
-   **Go Engine (v1.0.0):** A compiled, scratch-based backend validated with 100 concurrent clients at 50 publications per second.
-   **Reliable Replay:** Tenant-configurable retention replays missed SSE events through `Last-Event-ID`.

## ✨ Key Features
-   ✅ **High-Performance Go Backend:** Built with Go 1.25 for speed and safety.
-   ✅ **Scratch-Based Hardening:** The production image contains *zero* OS overhead, making it incredibly secure.
-   ✅ **Tab-Level Optimization:** Drastically reduces server load via `BroadcastChannel` synchronization.
-   ✅ **Security by Design:** Minimalist payloads ensure data privacy.
-   ✅ **Multi-Tenancy:** Manage multiple apps with unique, database-backed API keys.
-   ✅ **React Ready:** Native `usePravah()` hook for modern frontend integration.
-   ✅ **Project Agnostic:** Generic event handling for any domain.

## 🛠️ Tech Stack
-   **Backend:** Go 1.25 (Gin Framework)
-   **Database:** MongoDB
-   **Frontend:** JavaScript (Native SSE) + React Hook
-   **Deployment:** Docker (Scratch-based minimal images)

## 📦 React Integration
Pravah provides a dedicated hook for React projects:

```javascript
import { usePravah, usePravahListener } from 'pravah-sdk';

const MyComponent = () => {
    const pravah = usePravah({
        streams: ['test-stream'],
        getSubscriptionToken: fetchSubscriptionToken
    });

    usePravahListener(pravah, 'bid-updated', (data) => {
        console.log('Bid updated:', data.id);
    });

    return (
        <div>
            Status: {pravah?.connectionState || 'Initializing'}
        </div>
    );
};
```

## 🧪 Run locally

The Go service is the supported backend:

```bash
export SUBSCRIPTION_JWT_SECRET="$(openssl rand -hex 32)"
export ADMIN_API_KEY="$(openssl rand -hex 32)"
docker compose up --build pravah-go db
```

MongoDB is pinned to `9.0`. Override `MONGODB_DATA_PATH` to use a different
persistent data directory:

```bash
MONGODB_DATA_PATH=/srv/pravah/mongodb docker compose up -d db pravah-go
```

Run the complete first-event quickstart:

```bash
bash scripts/quickstart.sh
```

The OpenAPI contract is served at `/openapi.yaml` and stored at
[pravah-go/openapi.yaml](pravah-go/openapi.yaml).

Production operations are covered by [Operations](docs/operations.md),
[Troubleshooting](docs/troubleshooting.md), and
[MVP1 Known Limitations](docs/known-limitations.md).

The legacy Node implementation is disabled by default. It can be started
explicitly with `docker compose --profile legacy up pravah-legacy`.

## 🤝 Community & Contribution
Pravah is currently in its **Refinement Phase**. We are actively looking for developers to use, break, and help us improve the system.

- **Found a bug?** [Open an issue](.github/ISSUE_TEMPLATE/bug_report.yml)
- **Want a feature?** [Suggest it here](.github/ISSUE_TEMPLATE/feature_request.yml)
- **Want to contribute?** Check out [CONTRIBUTING.md](CONTRIBUTING.md)

---
**Lead Engineer:** Om Prakash Tiwari
**Vision:** Building a world where real-time flow is efficient, secure, and effortless.
