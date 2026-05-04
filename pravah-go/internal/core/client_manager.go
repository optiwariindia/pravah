package core

import (
	"sync"
	"time"
)

type Client struct {
	ID          string
	PravahID    string
	Streams     map[string]struct{}
	ConnectedAt time.Time
	Signal      chan string
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

func (cm *ClientManager) AddClient(connectionID, pravahID string, streams map[string]struct{}) *Client {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client := &Client{
		ID:          connectionID,
		PravahID:    pravahID,
		Streams:     streams,
		ConnectedAt: time.Now(),
		Signal:      make(chan string, 10), // Buffered channel
	}
	cm.clients[connectionID] = client
	return client
}

func (cm *ClientManager) RemoveClient(connectionID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.clients, connectionID)
}

func (cm *ClientManager) GetClientsByStream(streamID string) []*Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var matchingClients []*Client
	for _, client := range cm.clients {
		if _, ok := client.Streams[streamID]; ok {
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
			"pravahId":    c.PravahID,
			"streams":     streams,
			"connectedAt": c.ConnectedAt,
		})
	}
	return allClients
}
