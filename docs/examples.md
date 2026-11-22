# Integration Examples

## Publish from curl

```bash
curl -X POST "$PRAVAH_URL/pravah/publish" \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: $PRAVAH_API_KEY" \
  -d '{"streamId":"project:123","event":"resource-updated","data":{"id":"resource-456"}}'
```

## Publish from Go

```go
body := strings.NewReader(`{"streamId":"project:123","event":"resource-updated","data":{"id":"resource-456"}}`)
request, _ := http.NewRequest(http.MethodPost, pravahURL+"/pravah/publish", body)
request.Header.Set("Content-Type", "application/json")
request.Header.Set("X-API-KEY", apiKey)
response, err := http.DefaultClient.Do(request)
```

## Publish from Node.js

```javascript
const response = await fetch(`${pravahUrl}/pravah/publish`, {
  method: 'POST',
  headers: {
    'content-type': 'application/json',
    'x-api-key': apiKey
  },
  body: JSON.stringify({
    streamId: 'project:123',
    event: 'resource-updated',
    data: { id: 'resource-456' }
  })
});
```

## Subscribe in vanilla JavaScript

```javascript
import Pravah from 'pravah-sdk';

const pravah = new Pravah({
  baseUrl: 'https://pravah.example.com',
  streams: ['project:123'],
  getSubscriptionToken: async () => {
    const response = await fetch('/api/pravah-token', { credentials: 'include' });
    return (await response.json()).token;
  }
});

pravah.on('resource-updated', ({ id }) => refreshResource(id));
```

## Subscribe in React

```jsx
const pravah = usePravah({
  baseUrl: 'https://pravah.example.com',
  streams: ['project:123'],
  getSubscriptionToken
});

usePravahListener(pravah, 'resource-updated', ({ id }) => {
  queryClient.invalidateQueries(['resource', id]);
});
```

Only issue subscription tokens from an authenticated application backend. Do
not expose `X-API-KEY`, `X-ADMIN-KEY`, or signing secrets to the browser.
