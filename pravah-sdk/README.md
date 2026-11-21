# Pravah SDK 🌊

The high-performance, secure, real-time event notification SDK for **Pravah**. This package includes the core JavaScript SDK and React hooks for seamless integration.

## Installation

```bash
npm install pravah-sdk
```

## Vanilla JavaScript Usage

```javascript
import Pravah from 'pravah-sdk';

const pravah = new Pravah({
  baseUrl: 'https://your-pravah-instance.com',
  streams: ['global'],
  // Fetch a fresh token after reconnects or token expiry.
  getSubscriptionToken: async () => {
    const response = await fetch('/api/pravah/subscription-token', { credentials: 'include' });
    return (await response.json()).token;
  }
});

pravah.on('bid-updated', (data) => {
  console.log('New bid:', data.price);
});

pravah.on('statechange', (state) => {
  console.log('Connection state:', state);
});

// Call during application teardown.
pravah.close();
```

## React Usage

```javascript
import { usePravah, usePravahListener, usePravahState } from 'pravah-sdk';

const MyComponent = () => {
  const pravah = usePravah({
    baseUrl: 'https://your-pravah-instance.com',
    streams: ['user:123', 'global'],
    getSubscriptionToken: fetchSubscriptionToken
  });

  // 1. Listen for events using a callback
  usePravahListener(pravah, 'message-received', (data) => {
    console.log('New message signal:', data.id);
  });

  // 2. Or track event data in state automatically
  const latestAlert = usePravahState(pravah, 'system-alert');

  return (
    <div>
      <p>Latest alert: {latestAlert?.message || 'None'}</p>
    </div>
  );
};
```

## API Reference

### Core Class: `Pravah`
- `constructor(options)`: `{ baseUrl: string, streams: string[], subscriptionToken?: string, getSubscriptionToken?: () => Promise<string> }`
- `on(event, callback)`: Subscribe to an event.
- `off(event, callback)`: Unsubscribe from an event.
- `close()`: Stop reconnects and release EventSource, timers, listeners, and BroadcastChannel resources.
- `updateSubscriptionToken(token)`: Replace a static token and reconnect if needed.

The SDK emits `statechange` with `initializing`, `connecting`, `connected`,
`reconnecting`, `disconnected`, or `closed`.

### React Hooks
- `usePravah(options)`: Initializes and manages a Pravah instance.
- `usePravahListener(pravah, event, callback)`: Subscribes to an event within a component.
- `usePravahState(pravah, event)`: Returns the latest data for an event.

## License
ISC
