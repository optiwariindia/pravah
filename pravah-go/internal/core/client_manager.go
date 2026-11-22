package core

import (
	"sync"
	"time"

	"github.com/sampledge-sse/pravah-go/internal/utils"
)

type Client struct {
	ID          string
	AppID       string
	PravahID    string
	Streams     map[string]struct{}
	ConnectedAt time.Time
	Signal      chan Signal
	Done        chan struct{}
	closeOnce   sync.Once
	closeReason string
}

type Signal struct {
	ID    string
	Event string
	Data  map[string]interface{}
}

func (c *Client) Enqueue(signal Signal) bool {
	select {
	case <-c.Done:
		return false
	default:
	}

	select {
	case c.Signal <- signal:
		return true
	default:
		c.Close("buffer_full")
		return false
	}
}

func (c *Client) Close(reason string) {
	c.closeOnce.Do(func() {
		c.closeReason = reason
		close(c.Done)
	})
}

func (c *Client) CloseReason() string {
	select {
	case <-c.Done:
		return c.closeReason
	default:
		return ""
	}
}

type ClientManager struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var Instance *ClientManager
var once sync.Once
var StartTime time.Time

func init() {
	StartTime = time.Now()
}

func GetClientManager() *ClientManager {
	once.Do(func() {
		Instance = &ClientManager{
			clients: make(map[string]*Client),
		}
	})
	return Instance
}

func (cm *ClientManager) AddClient(connectionID, appID, pravahID string, streams map[string]struct{}) *Client {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client := &Client{
		ID:          connectionID,
		AppID:       appID,
		PravahID:    pravahID,
		Streams:     streams,
		ConnectedAt: time.Now(),
		Signal:      make(chan Signal, utils.ClientSignalBuffer()),
		Done:        make(chan struct{}),
	}
	cm.clients[connectionID] = client
	return client
}

func (cm *ClientManager) RemoveClient(connectionID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.clients, connectionID)
}

func (cm *ClientManager) CloseAll(reason string) {
	cm.mu.RLock()
	clients := make([]*Client, 0, len(cm.clients))
	for _, client := range cm.clients {
		clients = append(clients, client)
	}
	cm.mu.RUnlock()

	for _, client := range clients {
		client.Close(reason)
	}
}

func (cm *ClientManager) GetClientsByStream(appID, streamID string) []*Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var matchingClients []*Client
	for _, client := range cm.clients {
		_, subscribed := client.Streams[streamID]
		if client.AppID == appID && subscribed {
			matchingClients = append(matchingClients, client)
		}
	}
	return matchingClients
}

func (cm *ClientManager) GetConnectionCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.clients)
}

func (cm *ClientManager) GetAllClients() []map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	allClients := make([]map[string]interface{}, 0, len(cm.clients))
	for _, c := range cm.clients {
		streams := make([]string, 0, len(c.Streams))
		for s := range c.Streams {
			streams = append(streams, s)
		}
		allClients = append(allClients, map[string]interface{}{
			"id":          c.ID,
			"appId":       c.AppID,
			"pravahId":    c.PravahID,
			"streams":     streams,
			"connectedAt": c.ConnectedAt,
		})
	}
	return allClients
}
