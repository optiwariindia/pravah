package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/optiwariindia/pravah/pravah-go/internal/core"
)

func SubscribeHandler(c *gin.Context) {
	pravahID, err := c.Cookie("pravah_id")
	if err != nil || pravahID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "No Pravah ID found. Please init first.",
		})
		return
	}

	connectionID := uuid.New().String()

	streamsParam := c.Query("streams")
	streams := make(map[string]struct{})
	if streamsParam != "" {
		for _, s := range strings.Split(streamsParam, ",") {
			streams[strings.TrimSpace(s)] = struct{}{}
		}
	} else {
		streams["global"] = struct{}{}
	}

	// Set SSE Headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Send Initial Connection Data
	initData := map[string]interface{}{
		"event": "connected",
		"data": map[string]interface{}{
			"connectionId": connectionID,
			"pravahId":     pravahID,
			"streams":      getKeys(streams),
			"timestamp":    time.Now().Format(time.RFC3339),
		},
	}
	data, _ := json.Marshal(initData)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	c.Writer.Flush()

	// Register in ClientManager
	client := core.GetClientManager().AddClient(connectionID, pravahID, streams)

	log.Printf("[Pravah] Stream Opened: %s (User: %s) Streams: %s", connectionID, pravahID, strings.Join(getKeys(streams), ", "))

	// Create a channel for heartbeats
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-client.Signal:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			return true
		case <-ticker.C:
			heartbeat := map[string]interface{}{
				"event": "heartbeat",
				"data":  map[string]interface{}{"timestamp": time.Now().Format(time.RFC3339)},
			}
			data, _ := json.Marshal(heartbeat)
			fmt.Fprintf(w, "data: %s\n\n", data)
			return true
		case <-c.Request.Context().Done():
			core.GetClientManager().RemoveClient(connectionID)
			log.Printf("[Pravah] Stream Closed: %s (User: %s)", connectionID, pravahID)
			return false
		}
	})
	
	_ = client // client is stored in manager, we might use it for direct writes if we didn't use c.Stream
}

func getKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
