package api

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/core"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/models"
	"github.com/sampledge-sse/pravah-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublishRequest struct {
	StreamID string           `json:"streamId"`
	Event    string           `json:"event"`
	Data     models.EventData `json:"data"`
}

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9:._/-]{0,127}$`)

func PublishHandler(c *gin.Context) {
	var req PublishRequest
	if err := decodeStrictJSON(c, &req); err != nil {
		badRequest(c, "A valid JSON signal with streamId, event, and data.id is required; metadata is not accepted.")
		return
	}

	if !identifierPattern.MatchString(req.StreamID) ||
		!identifierPattern.MatchString(req.Event) ||
		!identifierPattern.MatchString(req.Data.ID) {
		badRequest(c, "streamId, event, and data.id must be valid identifiers of at most 128 characters.")
		return
	}

	app, ok := appFromContext(c)
	if !ok {
		return
	}
	config, err := getAppConfig(c.Request.Context(), app.AppID)
	if err != nil {
		db.MarkUnhealthy()
		log.Printf("[Pravah] Tenant config unavailable; using default retention for live delivery: %v", err)
		config = models.AppConfig{
			AppID:          app.AppID,
			RetentionHours: utils.EventRetentionHours(),
		}
	}

	createdAt := time.Now()
	event := models.Event{
		ID:        primitive.NewObjectID(),
		AppID:     app.AppID,
		StreamID:  req.StreamID,
		Event:     req.Event,
		Data:      req.Data,
		CreatedAt: createdAt,
		ExpiresAt: utils.EventExpiration(createdAt, config.RetentionHours),
	}

	persistContext, cancel := context.WithTimeout(c.Request.Context(), utils.DatabaseOperationTimeout())
	_, err = db.Database.Collection("events").InsertOne(persistContext, event)
	cancel()
	persisted := err == nil
	if err != nil {
		db.MarkUnhealthy()
		log.Printf("[Pravah] Persistence unavailable; publishing live-only signal: %v", err)
	}

	// 3. Broadcast
	clients := core.GetClientManager().GetClientsByStream(app.AppID, req.StreamID)

	signal := core.Signal{
		ID:    event.ID.Hex(),
		Event: req.Event,
		Data: map[string]interface{}{
			"id":        req.Data.ID,
			"timestamp": event.CreatedAt.Format(time.RFC3339),
		},
	}

	delivered := 0
	for _, client := range clients {
		if client.Enqueue(signal) {
			delivered++
		} else {
			log.Printf("[Pravah] Client %s closed: %s", client.ID, client.CloseReason())
		}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success":    true,
		"message":    "Event published",
		"count":      delivered,
		"eventId":    event.ID.Hex(),
		"persisted":  persisted,
		"replayable": persisted,
	})
}
