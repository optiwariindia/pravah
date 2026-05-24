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
  logLevel: 'info' // Optional: 'none' (default), 'error', 'warn', 'info'
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
    streams: ['user:123', 'global']
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

## React Context Provider

You can wrap your application with `PravahProvider` to avoid passing the `pravah` instance to every hook.

```javascript
import { PravahProvider, usePravahListener } from 'pravah-sdk';

const options = { baseUrl: 'https://your-pravah-service.com' };

const Root = () => (
  <PravahProvider options={options}>
    <MyComponent />
  </PravahProvider>
);

const MyComponent = () => {
  // Automatically uses the Pravah instance from context
  usePravahListener(null, 'event-name', (data) => {
    console.log(data);
  });
  
  return <div>My Component</div>;
};
```

## API Reference

### Core Class: `Pravah`
- `constructor(options)`: `{ baseUrl: string, streams: string[], logLevel?: 'none' | 'error' | 'warn' | 'info' }`
- `on(event, callback)`: Subscribe to an event.
- `off(event, callback)`: Unsubscribe from an event.

### React Hooks
- `PravahProvider`: Context provider. Props: `{ options: PravahOptions, children: ReactNode }`.
- `usePravah(options)`: Initializes and manages a Pravah instance.
- `usePravahContext()`: Access the Pravah instance from the provider.
- `usePravahListener(pravah, event, callback)`: Subscribes to an event within a component. `pravah` is optional if inside `PravahProvider`.
- `usePravahState(pravah, event)`: Returns the latest data for an event. `pravah` is optional if inside `PravahProvider`.

## License
ISC
