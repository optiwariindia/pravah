package metrics

import "testing"

func TestPublicationSnapshot(t *testing.T) {
	before := GetSnapshot()

	Publication(25_000_000, 1, 3, 2, true)

	after := GetSnapshot()
	if after.SignalsPublished != before.SignalsPublished+1 {
		t.Fatal("published metric was not updated")
	}
	if after.SignalsDelivered != before.SignalsDelivered+3 {
		t.Fatal("delivery metric was not updated")
	}
	if after.SignalsDropped != before.SignalsDropped+2 {
		t.Fatal("drop metric was not updated")
	}
	if after.PersistenceFailures != before.PersistenceFailures+1 {
		t.Fatal("persistence metric was not updated")
	}
	if after.AveragePublicationMs <= 0 {
		t.Fatal("average publication latency was not updated")
	}
}

func TestConnectionSnapshot(t *testing.T) {
	before := GetSnapshot()
	ConnectionOpened()
	ConnectionClosed()
	after := GetSnapshot()

	if after.ActiveConnections != before.ActiveConnections {
		t.Fatal("active connections did not return to baseline")
	}
	if after.ConnectionsOpened != before.ConnectionsOpened+1 {
		t.Fatal("opened connection metric was not updated")
	}
	if after.ConnectionsClosed != before.ConnectionsClosed+1 {
		t.Fatal("closed connection metric was not updated")
	}
}
