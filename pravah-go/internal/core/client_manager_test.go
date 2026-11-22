package core

import (
	"fmt"
	"sync"
	"testing"
)

func TestClientEnqueueClosesFullBuffer(t *testing.T) {
	t.Setenv("CLIENT_SIGNAL_BUFFER", "1")
	manager := &ClientManager{clients: make(map[string]*Client)}
	client := manager.AddClient("connection-1", "app-1", "browser-1", map[string]struct{}{"global": {}})

	if !client.Enqueue(Signal{ID: "event-1"}) {
		t.Fatal("expected first signal to be queued")
	}
	if client.Enqueue(Signal{ID: "event-2"}) {
		t.Fatal("expected a full client buffer to reject the second signal")
	}
	if client.CloseReason() != "buffer_full" {
		t.Fatalf("expected buffer_full close reason, got %q", client.CloseReason())
	}

	select {
	case <-client.Done:
	default:
		t.Fatal("expected client to be marked closed")
	}
}

func TestClientEnqueueRejectsClosedClient(t *testing.T) {
	client := &Client{
		Signal: make(chan Signal, 1),
		Done:   make(chan struct{}),
	}
	client.Close("shutdown")

	if client.Enqueue(Signal{ID: "event-1"}) {
		t.Fatal("expected closed client to reject signals")
	}
	if len(client.Signal) != 0 {
		t.Fatal("closed client unexpectedly queued a signal")
	}
}

func TestCloseAllClosesEveryClient(t *testing.T) {
	manager := &ClientManager{clients: make(map[string]*Client)}
	first := manager.AddClient("connection-1", "app-1", "browser-1", nil)
	second := manager.AddClient("connection-2", "app-1", "browser-2", nil)

	manager.CloseAll("server_shutdown")

	for _, client := range []*Client{first, second} {
		select {
		case <-client.Done:
			if client.CloseReason() != "server_shutdown" {
				t.Fatalf("unexpected close reason: %q", client.CloseReason())
			}
		default:
			t.Fatalf("client %s was not closed", client.ID)
		}
	}
}

func TestConcurrentPublishLookupAndRemoval(t *testing.T) {
	manager := &ClientManager{clients: make(map[string]*Client)}
	const iterations = 500
	for index := 0; index < iterations; index++ {
		id := fmt.Sprintf("connection-%d", index)
		manager.AddClient(id, "app-1", "browser", map[string]struct{}{"stream": {}})
	}
	var workers sync.WaitGroup

	workers.Add(2)
	go func() {
		defer workers.Done()
		for index := 0; index < iterations; index++ {
			for _, client := range manager.GetClientsByStream("app-1", "stream") {
				client.Enqueue(Signal{ID: "event"})
			}
		}
	}()
	go func() {
		defer workers.Done()
		for index := 0; index < iterations; index++ {
			manager.RemoveClient(fmt.Sprintf("connection-%d", index))
		}
	}()

	workers.Wait()
	if count := manager.GetConnectionCount(); count != 0 {
		t.Fatalf("expected all clients removed, got %d", count)
	}
}
