package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/optiwariindia/pravah/pravah-go/internal/core"
	"github.com/optiwariindia/pravah/pravah-go/internal/db"
	"github.com/optiwariindia/pravah/pravah-go/internal/models"
)

type PublishRequest struct {
	StreamID string           `json:"streamId" binding:"required"`
	Event    string           `json:"event" binding:"required"`
	Data     models.EventData `json:"data" binding:"required"`
}

func PublishHandler(c *gin.Context) {
	var req PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "streamId, event, and data.id are required."})
		return
	}

	if req.Data.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "data.id is required."})
		return
	}

	// 2. Persistence
	event := models.Event{
		StreamID:  req.StreamID,
		Event:     req.Event,
		Data:      req.Data,
		CreatedAt: time.Now(),
	}

	_, err := db.Database.Collection("events").InsertOne(context.TODO(), event)
	if err != nil {
		log.Printf("[Pravah] DB Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	// 3. Broadcast
	clients := core.GetClientManager().GetClientsByStream(req.StreamID)

	signal := map[string]interface{}{
		"event": req.Event,
		"data": map[string]interface{}{
			"id":        req.Data.ID,
			"timestamp": event.CreatedAt.Format(time.RFC3339),
		},
	}
	signalBytes, _ := json.Marshal(signal)
	signalString := string(signalBytes)

	for _, client := range clients {
		select {
		case client.Signal <- signalString:
		default:
			log.Printf("[Pravah] Client %s buffer full, dropping message", client.ID)
		}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Event published",
		"count":   len(clients),
	})
}
