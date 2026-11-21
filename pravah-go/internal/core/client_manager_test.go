package core

import (
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
