package auth

import (
	"testing"
	"time"
)

func TestSubscriptionTokenRoundTrip(t *testing.T) {
	t.Setenv("SUBSCRIPTION_JWT_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef")

	token, expiresAt, err := CreateSubscriptionToken("app-1", "user-1", []string{"project:123"}, time.Minute)
	if err != nil {
		t.Fatalf("CreateSubscriptionToken returned an error: %v", err)
	}

	claims, err := ValidateSubscriptionToken(token)
	if err != nil {
		t.Fatalf("ValidateSubscriptionToken returned an error: %v", err)
	}
	if claims.AppID != "app-1" || claims.Subject != "user-1" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if len(claims.Streams) != 1 || claims.Streams[0] != "project:123" {
		t.Fatalf("unexpected streams: %v", claims.Streams)
	}
	if claims.ExpiresAt != expiresAt.Unix() {
		t.Fatalf("expected expiration %d, got %d", expiresAt.Unix(), claims.ExpiresAt)
	}
}

func TestSubscriptionTokenRejectsTampering(t *testing.T) {
	t.Setenv("SUBSCRIPTION_JWT_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef")

	token, _, err := CreateSubscriptionToken("app-1", "user-1", []string{"project:123"}, time.Minute)
	if err != nil {
		t.Fatalf("CreateSubscriptionToken returned an error: %v", err)
	}
	tampered := token[:len(token)-1] + "x"
	if tampered == token {
		tampered = token[:len(token)-1] + "y"
	}
	if _, err := ValidateSubscriptionToken(tampered); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}

func TestSubscriptionTokenRejectsMissingSecret(t *testing.T) {
	t.Setenv("SUBSCRIPTION_JWT_SECRET", "short")
	_, _, err := CreateSubscriptionToken("app-1", "user-1", []string{"project:123"}, time.Minute)
	if err == nil {
		t.Fatal("expected short secret to be rejected")
	}
}

func TestSubscriptionTokenRejectsExpiredToken(t *testing.T) {
	t.Setenv("SUBSCRIPTION_JWT_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef")
	token, _, err := CreateSubscriptionToken("app-1", "user-1", []string{"project:123"}, -time.Second)
	if err != nil {
		t.Fatalf("CreateSubscriptionToken returned an error: %v", err)
	}
	if _, err := ValidateSubscriptionToken(token); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}
