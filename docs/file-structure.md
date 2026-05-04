# File Structure

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