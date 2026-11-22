package metrics

import "sync/atomic"

type counters struct {
	activeConnections      atomic.Int64
	connectionsOpened      atomic.Uint64
	connectionsClosed      atomic.Uint64
	signalsPublished       atomic.Uint64
	signalsDelivered       atomic.Uint64
	signalsDropped         atomic.Uint64
	replayEvents           atomic.Uint64
	replayErrors           atomic.Uint64
	authenticationFailures atomic.Uint64
	persistenceFailures    atomic.Uint64
	publicationCount       atomic.Uint64
	publicationNanos       atomic.Uint64
}

var current counters

func ConnectionOpened() {
	current.activeConnections.Add(1)
	current.connectionsOpened.Add(1)
}

func ConnectionClosed() {
	current.activeConnections.Add(-1)
	current.connectionsClosed.Add(1)
}

func Publication(durationNanos int64, published, delivered, dropped uint64, persistenceFailed bool) {
	current.publicationCount.Add(1)
	if durationNanos > 0 {
		current.publicationNanos.Add(uint64(durationNanos))
	}
	current.signalsPublished.Add(published)
	current.signalsDelivered.Add(delivered)
	current.signalsDropped.Add(dropped)
	if persistenceFailed {
		current.persistenceFailures.Add(1)
	}
}

func Replay(events uint64) {
	current.replayEvents.Add(events)
}

func ReplayError() {
	current.replayErrors.Add(1)
}

func AuthenticationFailure() {
	current.authenticationFailures.Add(1)
}

type Snapshot struct {
	ActiveConnections      int64   `json:"activeConnections"`
	ConnectionsOpened      uint64  `json:"connectionsOpened"`
	ConnectionsClosed      uint64  `json:"connectionsClosed"`
	SignalsPublished       uint64  `json:"signalsPublished"`
	SignalsDelivered       uint64  `json:"signalsDelivered"`
	SignalsDropped         uint64  `json:"signalsDropped"`
	ReplayEvents           uint64  `json:"replayEvents"`
	ReplayErrors           uint64  `json:"replayErrors"`
	AuthenticationFailures uint64  `json:"authenticationFailures"`
	PersistenceFailures    uint64  `json:"persistenceFailures"`
	AveragePublicationMs   float64 `json:"averagePublicationMs"`
}

func GetSnapshot() Snapshot {
	count := current.publicationCount.Load()
	average := float64(0)
	if count > 0 {
		average = float64(current.publicationNanos.Load()) / float64(count) / float64(1e6)
	}
	return Snapshot{
		ActiveConnections:      current.activeConnections.Load(),
		ConnectionsOpened:      current.connectionsOpened.Load(),
		ConnectionsClosed:      current.connectionsClosed.Load(),
		SignalsPublished:       current.signalsPublished.Load(),
		SignalsDelivered:       current.signalsDelivered.Load(),
		SignalsDropped:         current.signalsDropped.Load(),
		ReplayEvents:           current.replayEvents.Load(),
		ReplayErrors:           current.replayErrors.Load(),
		AuthenticationFailures: current.authenticationFailures.Load(),
		PersistenceFailures:    current.persistenceFailures.Load(),
		AveragePublicationMs:   average,
	}
}
