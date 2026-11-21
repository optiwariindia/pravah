package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	tokenIssuer   = "pravah"
	tokenAudience = "pravah-subscribe"
)

type SubscriptionClaims struct {
	AppID     string   `json:"appId"`
	Streams   []string `json:"streams"`
	Issuer    string   `json:"iss"`
	Audience  string   `json:"aud"`
	Subject   string   `json:"sub"`
	ExpiresAt int64    `json:"exp"`
	IssuedAt  int64    `json:"iat"`
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

func CreateSubscriptionToken(appID, subject string, streams []string, lifetime time.Duration) (string, time.Time, error) {
	secret, err := subscriptionSecret()
	if err != nil {
		return "", time.Time{}, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(lifetime)
	header, err := json.Marshal(tokenHeader{Algorithm: "HS256", Type: "JWT"})
	if err != nil {
		return "", time.Time{}, err
	}
	claims, err := json.Marshal(SubscriptionClaims{
		AppID:     appID,
		Streams:   streams,
		Issuer:    tokenIssuer,
		Audience:  tokenAudience,
		Subject:   subject,
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
	})
	if err != nil {
		return "", time.Time{}, err
	}

	unsigned := encodeJWTPart(header) + "." + encodeJWTPart(claims)
	signature := sign(unsigned, secret)
	return unsigned + "." + signature, expiresAt, nil
}

func ValidateSubscriptionToken(raw string) (SubscriptionClaims, error) {
	secret, err := subscriptionSecret()
	if err != nil {
		return SubscriptionClaims{}, err
	}

	parts := strings.Split(raw, ".")
	if len(parts) != 3 {
		return SubscriptionClaims{}, fmt.Errorf("invalid token format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return SubscriptionClaims{}, fmt.Errorf("invalid token header")
	}
	var header tokenHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Algorithm != "HS256" || header.Type != "JWT" {
		return SubscriptionClaims{}, fmt.Errorf("unsupported token header")
	}

	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(sign(unsigned, secret)), []byte(parts[2])) {
		return SubscriptionClaims{}, fmt.Errorf("invalid token signature")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return SubscriptionClaims{}, fmt.Errorf("invalid token claims")
	}
	var claims SubscriptionClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return SubscriptionClaims{}, fmt.Errorf("invalid token claims")
	}

	now := time.Now().Unix()
	if claims.AppID == "" || claims.Subject == "" || len(claims.Streams) == 0 {
		return SubscriptionClaims{}, fmt.Errorf("incomplete token claims")
	}
	if claims.Issuer != tokenIssuer || claims.Audience != tokenAudience {
		return SubscriptionClaims{}, fmt.Errorf("invalid token issuer or audience")
	}
	if claims.IssuedAt > now+30 || claims.ExpiresAt <= now {
		return SubscriptionClaims{}, fmt.Errorf("token is expired or not yet valid")
	}
	return claims, nil
}

func subscriptionSecret() ([]byte, error) {
	secret := os.Getenv("SUBSCRIPTION_JWT_SECRET")
	if len(secret) < 32 {
		return nil, fmt.Errorf("SUBSCRIPTION_JWT_SECRET must contain at least 32 characters")
	}
	return []byte(secret), nil
}

func encodeJWTPart(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}

func sign(value string, secret []byte) string {
	hash := hmac.New(sha256.New, secret)
	_, _ = hash.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}
