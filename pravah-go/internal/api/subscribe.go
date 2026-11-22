package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sampledge-sse/pravah-go/internal/auth"
	"github.com/sampledge-sse/pravah-go/internal/core"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/metrics"
	"github.com/sampledge-sse/pravah-go/internal/models"
	"github.com/sampledge-sse/pravah-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	rawToken := c.Query("token")
	claims, err := auth.ValidateSubscriptionToken(rawToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "A valid subscription token is required.",
		})
		return
	}

	connectionID := uuid.New().String()

	streamsParam := c.Query("streams")
	streams := make(map[string]struct{})
	if streamsParam != "" {
		for _, s := range strings.Split(streamsParam, ",") {
			stream := strings.TrimSpace(s)
			if stream != "" {
				streams[stream] = struct{}{}
			}
		}
	}
	if len(streams) == 0 {
		streams["global"] = struct{}{}
	}
	if len(streams) > utils.MaxSubscriptionStreams() {
		badRequest(c, fmt.Sprintf("A subscription may contain at most %d streams.", utils.MaxSubscriptionStreams()))
		return
	}
	authorized, streams := authorizedStreams(claims, streams)
	if !authorized {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "The subscription token does not authorize one or more requested streams.",
		})
		return
	}

	lastEventID := c.GetHeader("Last-Event-ID")
	if lastEventID == "" {
		lastEventID = c.Query("lastEventId")
	}
	if lastEventID != "" {
		if _, err := primitive.ObjectIDFromHex(lastEventID); err != nil {
			badRequest(c, "Last-Event-ID must be a valid event identifier.")
			return
		}
	}

	client := core.GetClientManager().AddClient(connectionID, claims.AppID, pravahID, streams)
	defer core.GetClientManager().RemoveClient(connectionID)

	replay, err := loadReplayEvents(c, claims.AppID, streams, lastEventID)
	if err != nil {
		metrics.ReplayError()
		if !db.IsHealthy() && lastEventID != "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Replay Unavailable",
				"message": "MongoDB is unavailable. Reconnect without Last-Event-ID to receive live events.",
			})
			return
		}
		if errors.Is(err, errReplayLimitExceeded) {
			c.JSON(http.StatusGone, gin.H{
				"error":   "Replay Limit Exceeded",
				"message": "Reconnect without Last-Event-ID and refresh the application state.",
			})
			return
		}
		slog.Error("replay query failed", "connectionId", connectionID, "appId", claims.AppID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Unable to replay missed events."})
		return
	}

	// Set SSE Headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Send Initial Connection Data
	initData := map[string]interface{}{
		"connectionId": connectionID,
		"pravahId":     pravahID,
		"streams":      getKeys(streams),
		"timestamp":    time.Now().Format(time.RFC3339),
		"degraded":     !db.IsHealthy(),
	}
	if err := writeSSEEvent(c.Writer, "", "connected", initData); err != nil {
		return
	}
	c.Writer.Flush()

	metrics.Replay(uint64(len(replay)))
	slog.Info("stream opened", "connectionId", connectionID, "userId", pravahID, "appId", claims.AppID, "streams", getKeys(streams))

	replayedThrough := ""
	for _, signal := range replay {
		if err := writeSignalEvent(c.Writer, signal); err != nil {
			return
		}
		replayedThrough = signal.ID
		c.Writer.Flush()
	}

	ticker := time.NewTicker(utils.HeartbeatInterval())
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-client.Signal:
			if replayedThrough != "" && msg.ID <= replayedThrough {
				return true
			}
			if err := writeSignalEvent(w, msg); err != nil {
				return false
			}
			return true
		case <-client.Done:
			slog.Info("stream closed", "connectionId", connectionID, "reason", client.CloseReason())
			return false
		case <-ticker.C:
			heartbeat := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
			}
			if err := writeSSEEvent(w, "", "heartbeat", heartbeat); err != nil {
				return false
			}
			return true
		case <-c.Request.Context().Done():
			slog.Info("stream closed", "connectionId", connectionID, "reason", "request_done")
			return false
		}
	})
}

func getKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

var errReplayLimitExceeded = errors.New("replay event limit exceeded")

func loadReplayEvents(c *gin.Context, appID string, streams map[string]struct{}, lastEventID string) ([]core.Signal, error) {
	if lastEventID == "" {
		return nil, nil
	}

	lastID, err := primitive.ObjectIDFromHex(lastEventID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":       bson.M{"$gt": lastID},
		"appId":     appID,
		"streamId":  bson.M{"$in": getKeys(streams)},
		"expiresAt": bson.M{"$gt": time.Now()},
	}
	cursor, err := db.Database.Collection("events").Find(
		c.Request.Context(),
		filter,
		options.Find().
			SetSort(bson.D{{Key: "_id", Value: 1}}).
			SetLimit(int64(utils.MaxReplayEvents()+1)),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c.Request.Context())

	var events []models.Event
	if err := cursor.All(c.Request.Context(), &events); err != nil {
		return nil, err
	}
	if len(events) > utils.MaxReplayEvents() {
		return nil, errReplayLimitExceeded
	}

	signals := make([]core.Signal, 0, len(events))
	for _, event := range events {
		signals = append(signals, core.Signal{
			ID:    event.ID.Hex(),
			Event: event.Event,
			Data: map[string]interface{}{
				"id":        event.Data.ID,
				"timestamp": event.CreatedAt.Format(time.RFC3339),
			},
		})
	}

	return signals, nil
}

func writeSSEEvent(w io.Writer, id, event string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if id != "" {
		if _, err := fmt.Fprintf(w, "id: %s\n", id); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, payload)
	return err
}

func writeSignalEvent(w io.Writer, signal core.Signal) error {
	return writeSSEEvent(w, signal.ID, "signal", map[string]interface{}{
		"event": signal.Event,
		"data":  signal.Data,
	})
}
