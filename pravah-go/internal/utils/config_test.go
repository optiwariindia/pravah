package utils

import (
	"testing"
	"time"
)

func TestEventRetentionHours(t *testing.T) {
	t.Setenv("EVENT_RETENTION_HOURS", "72")
	if got := EventRetentionHours(); got != 72 {
		t.Fatalf("expected configured retention of 72, got %d", got)
	}
}

func TestEventRetentionHoursDefaultsForInvalidValues(t *testing.T) {
	for _, value := range []string{"", "zero", "0", "-1", "8761"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("EVENT_RETENTION_HOURS", value)
			if got := EventRetentionHours(); got != DefaultRetentionHours {
				t.Fatalf("expected default retention of %d, got %d", DefaultRetentionHours, got)
			}
		})
	}
}

func TestEventExpiration(t *testing.T) {
	createdAt := time.Date(2026, time.November, 21, 0, 0, 0, 0, time.UTC)
	want := createdAt.Add(24 * time.Hour)
	if got := EventExpiration(createdAt, 24); !got.Equal(want) {
		t.Fatalf("expected expiration %s, got %s", want, got)
	}
}
