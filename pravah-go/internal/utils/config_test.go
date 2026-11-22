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

func TestOperationalLimitsFromEnvironment(t *testing.T) {
	t.Setenv("SSE_HEARTBEAT_INTERVAL_SECONDS", "45")
	t.Setenv("CLIENT_SIGNAL_BUFFER", "128")
	t.Setenv("MAX_SUBSCRIPTION_STREAMS", "16")
	t.Setenv("MAX_REPLAY_EVENTS", "500")
	t.Setenv("MAX_REQUEST_BODY_BYTES", "8192")
	t.Setenv("DATABASE_OPERATION_TIMEOUT_SECONDS", "5")

	if got := HeartbeatInterval(); got != 45*time.Second {
		t.Fatalf("unexpected heartbeat interval: %s", got)
	}
	if got := ClientSignalBuffer(); got != 128 {
		t.Fatalf("unexpected client buffer: %d", got)
	}
	if got := MaxSubscriptionStreams(); got != 16 {
		t.Fatalf("unexpected stream limit: %d", got)
	}
	if got := MaxReplayEvents(); got != 500 {
		t.Fatalf("unexpected replay limit: %d", got)
	}
	if got := MaxRequestBodyBytes(); got != 8192 {
		t.Fatalf("unexpected request body limit: %d", got)
	}
	if got := DatabaseOperationTimeout(); got != 5*time.Second {
		t.Fatalf("unexpected database operation timeout: %s", got)
	}
}

func TestOperationalLimitsUseDefaultsForInvalidValues(t *testing.T) {
	t.Setenv("SSE_HEARTBEAT_INTERVAL_SECONDS", "1")
	t.Setenv("CLIENT_SIGNAL_BUFFER", "0")
	t.Setenv("MAX_SUBSCRIPTION_STREAMS", "129")
	t.Setenv("MAX_REPLAY_EVENTS", "invalid")
	t.Setenv("MAX_REQUEST_BODY_BYTES", "999999999")
	t.Setenv("DATABASE_OPERATION_TIMEOUT_SECONDS", "0")

	if HeartbeatInterval() != 30*time.Second {
		t.Fatal("unexpected heartbeat fallback")
	}
	if ClientSignalBuffer() != 64 {
		t.Fatal("unexpected buffer fallback")
	}
	if MaxSubscriptionStreams() != 32 {
		t.Fatal("unexpected stream fallback")
	}
	if MaxReplayEvents() != 1000 {
		t.Fatal("unexpected replay fallback")
	}
	if MaxRequestBodyBytes() != 16384 {
		t.Fatal("unexpected body-size fallback")
	}
	if DatabaseOperationTimeout() != 3*time.Second {
		t.Fatal("unexpected database operation timeout fallback")
	}
}

func TestEnvironmentValidationRejectsInvalidOperationalLimit(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("MONGODB", "mongodb://db/pravah")
	t.Setenv("SUBSCRIPTION_JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("ADMIN_API_KEY", "abcdef0123456789abcdef0123456789")
	t.Setenv("CLIENT_SIGNAL_BUFFER", "0")

	if err := CheckEnvironmentVariables(); err == nil {
		t.Fatal("expected invalid CLIENT_SIGNAL_BUFFER to fail startup validation")
	}
}

func TestCookieConfiguration(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "true")
	t.Setenv("COOKIE_SAMESITE", "strict")
	t.Setenv("COOKIE_DOMAIN", "events.example.com")
	t.Setenv("COOKIE_MAX_AGE_SECONDS", "7200")

	if !CookieSecure() {
		t.Fatal("expected secure cookies")
	}
	if CookieSameSite() != "strict" {
		t.Fatalf("unexpected SameSite value: %q", CookieSameSite())
	}
	if CookieDomain() != "events.example.com" {
		t.Fatalf("unexpected cookie domain: %q", CookieDomain())
	}
	if CookieMaxAgeSeconds() != 7200 {
		t.Fatalf("unexpected cookie max age: %d", CookieMaxAgeSeconds())
	}
}

func TestCookieConfigurationDefaults(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "invalid")
	t.Setenv("COOKIE_SAMESITE", "invalid")
	t.Setenv("COOKIE_MAX_AGE_SECONDS", "1")

	if !CookieSecure() {
		t.Fatal("expected secure cookies by default")
	}
	if CookieSameSite() != "lax" {
		t.Fatalf("unexpected SameSite fallback: %q", CookieSameSite())
	}
	if CookieMaxAgeSeconds() != 365*24*60*60 {
		t.Fatalf("unexpected max-age fallback: %d", CookieMaxAgeSeconds())
	}
}
