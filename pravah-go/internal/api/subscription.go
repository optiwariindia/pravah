package api

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/auth"
	"github.com/sampledge-sse/pravah-go/internal/utils"
)

const (
	defaultTokenLifetime = 15 * time.Minute
	maxTokenLifetime     = time.Hour
)

type SubscriptionTokenRequest struct {
	Streams   []string `json:"streams"`
	ExpiresIn int      `json:"expiresIn,omitempty"`
	Subject   string   `json:"subject,omitempty"`
}

func CreateSubscriptionTokenHandler(c *gin.Context) {
	app, ok := appFromContext(c)
	if !ok {
		return
	}

	var req SubscriptionTokenRequest
	if err := decodeStrictJSON(c, &req); err != nil || len(req.Streams) == 0 || len(req.Streams) > utils.MaxSubscriptionStreams() {
		badRequest(c, "A valid JSON body with 1 to the configured maximum number of streams is required.")
		return
	}

	streams := uniqueValidStreams(req.Streams)
	if len(streams) != len(req.Streams) {
		badRequest(c, "streams must be unique, valid identifiers of at most 128 characters.")
		return
	}

	lifetime := time.Duration(req.ExpiresIn) * time.Second
	if req.ExpiresIn == 0 {
		lifetime = defaultTokenLifetime
	} else if lifetime < time.Second || lifetime > maxTokenLifetime {
		badRequest(c, "expiresIn must be between 1 and 3600 seconds.")
		return
	}

	subject := strings.TrimSpace(req.Subject)
	if len(subject) > 128 {
		badRequest(c, "subject must be at most 128 characters.")
		return
	}
	if subject == "" {
		subject = app.AppID
	}
	token, expiresAt, err := auth.CreateSubscriptionToken(app.AppID, subject, streams, lifetime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Subscription tokens are not configured."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"token":     token,
			"expiresAt": expiresAt,
			"streams":   streams,
		},
	})
}

func uniqueValidStreams(streams []string) []string {
	seen := make(map[string]struct{}, len(streams))
	result := make([]string, 0, len(streams))
	for _, stream := range streams {
		if !identifierPattern.MatchString(stream) {
			return nil
		}
		if _, exists := seen[stream]; exists {
			return nil
		}
		seen[stream] = struct{}{}
		result = append(result, stream)
	}
	sort.Strings(result)
	return result
}

func authorizedStreams(claims auth.SubscriptionClaims, requested map[string]struct{}) (bool, map[string]struct{}) {
	allowed := make(map[string]struct{}, len(claims.Streams))
	for _, stream := range claims.Streams {
		allowed[stream] = struct{}{}
	}
	intersection := make(map[string]struct{}, len(requested))
	for stream := range requested {
		if _, ok := allowed[stream]; !ok {
			return false, nil
		}
		intersection[stream] = struct{}{}
	}
	return true, intersection
}
