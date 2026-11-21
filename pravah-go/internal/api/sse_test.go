package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sampledge-sse/pravah-go/internal/auth"
	"github.com/sampledge-sse/pravah-go/internal/core"
)

func TestWriteSSEEvent(t *testing.T) {
	var output bytes.Buffer
	err := writeSSEEvent(&output, "0123456789abcdef01234567", "bid-updated", map[string]string{"id": "bid-1"})
	if err != nil {
		t.Fatalf("writeSSEEvent returned an error: %v", err)
	}

	want := "id: 0123456789abcdef01234567\n" +
		"event: bid-updated\n" +
		"data: {\"id\":\"bid-1\"}\n\n"
	if output.String() != want {
		t.Fatalf("unexpected SSE output:\nwant: %q\n got: %q", want, output.String())
	}
}

func TestWriteSSEEventWithoutID(t *testing.T) {
	var output bytes.Buffer
	if err := writeSSEEvent(&output, "", "heartbeat", map[string]string{"timestamp": "now"}); err != nil {
		t.Fatalf("writeSSEEvent returned an error: %v", err)
	}
	if strings.HasPrefix(output.String(), "id:") {
		t.Fatalf("heartbeat unexpectedly contained an event ID: %q", output.String())
	}
}

func TestWriteSignalEvent(t *testing.T) {
	var output bytes.Buffer
	signal := core.Signal{
		ID:    "0123456789abcdef01234567",
		Event: "bid-updated",
		Data:  map[string]interface{}{"id": "bid-1"},
	}

	if err := writeSignalEvent(&output, signal); err != nil {
		t.Fatalf("writeSignalEvent returned an error: %v", err)
	}

	expected := "id: 0123456789abcdef01234567\n" +
		"event: signal\n" +
		"data: {\"data\":{\"id\":\"bid-1\"},\"event\":\"bid-updated\"}\n\n"
	if output.String() != expected {
		t.Fatalf("unexpected signal output:\nwant: %q\n got: %q", expected, output.String())
	}
}

func TestIdentifierPattern(t *testing.T) {
	valid := []string{"global", "project:123", "user/456", "bid.updated"}
	for _, value := range valid {
		if !identifierPattern.MatchString(value) {
			t.Errorf("expected %q to be valid", value)
		}
	}

	invalid := []string{"", " project:123", "bad\nvalue", strings.Repeat("a", 129)}
	for _, value := range invalid {
		if identifierPattern.MatchString(value) {
			t.Errorf("expected %q to be invalid", value)
		}
	}
}

func TestAuthorizedStreamsRejectsUnapprovedStream(t *testing.T) {
	claims := auth.SubscriptionClaims{
		AppID:   "app-1",
		Streams: []string{"project:123"},
	}
	ok, streams := authorizedStreams(claims, map[string]struct{}{"project:123": {}, "user:456": {}})
	if ok || streams != nil {
		t.Fatal("expected mixed authorized and unauthorized streams to be rejected")
	}
}

func TestAuthorizedStreamsReturnsApprovedIntersection(t *testing.T) {
	claims := auth.SubscriptionClaims{
		AppID:   "app-1",
		Streams: []string{"project:123", "user:456"},
	}
	ok, streams := authorizedStreams(claims, map[string]struct{}{"project:123": {}})
	if !ok || len(streams) != 1 {
		t.Fatalf("expected approved stream, got ok=%v streams=%v", ok, streams)
	}
}
