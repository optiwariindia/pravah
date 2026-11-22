# File Structure

> The supported implementation is `pravah-go/`. The diagram below describes
> the original Node prototype and remains as historical context.

```
event-service/
│
├── src/
│   ├── server.js
│   │
│   ├── modules/
│   │
│   │   ├── subscribe/
│   │   │   ├── route.js
│   │   │   ├── controller.js
│   │   │   ├── service.js
│   │   │   └── validator.js   // optional
│   │   │
│   │   ├── publish/
│   │   │   ├── route.js
│   │   │   ├── controller.js
│   │   │   └── service.js
│   │   │
│   │   └── health/
│   │       └── route.js
│   │
│   ├── core/
│   │   ├── clientManager.js
│   │   ├── eventBus.js
│   │   ├── replayService.js
│   │   └── sse.js
│   │
│   ├── models/
│   │   └── event.js
│   │
│   ├── middleware/
│   │   └── auth.js
│   │
│   ├── config/
│   │   └── index.js
│   │
│   └── utils/
│       ├── logger.js
│       └── stream.js
│
├── package.json
└── README.md
```