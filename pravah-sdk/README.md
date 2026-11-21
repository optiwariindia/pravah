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
  subscriptionToken: '<short-lived-token>'
});

pravah.on('bid-updated', (data) => {
  console.log('New bid:', data.price);
});
```

## React Usage

```javascript
import { usePravah, usePravahListener, usePravahState } from 'pravah-sdk';

const MyComponent = () => {
  const pravah = usePravah({
    baseUrl: 'https://your-pravah-instance.com',
    streams: ['user:123', 'global'],
    subscriptionToken: '<short-lived-token>'
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
- `constructor(options)`: `{ baseUrl: string, streams: string[], subscriptionToken: string }`
- `on(event, callback)`: Subscribe to an event.
- `off(event, callback)`: Unsubscribe from an event.

### React Hooks
- `usePravah(options)`: Initializes and manages a Pravah instance.
- `usePravahListener(pravah, event, callback)`: Subscribes to an event within a component.
- `usePravahState(pravah, event)`: Returns the latest data for an event.

## License
ISC
