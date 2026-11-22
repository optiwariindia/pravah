package api

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/models"
	"github.com/sampledge-sse/pravah-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type RegisterRequest struct {
	Name           string `json:"name"`
	RetentionHours int    `json:"retentionHours,omitempty"`
}

func RegisterAppHandler(c *gin.Context) {
	var req RegisterRequest
	if err := decodeStrictJSON(c, &req); err != nil || strings.TrimSpace(req.Name) == "" {
		badRequest(c, "A valid JSON body with a name is required.")
		return
	}

	retentionHours := req.RetentionHours
	if retentionHours == 0 {
		retentionHours = utils.EventRetentionHours()
	}
	if !validRetentionHours(retentionHours) {
		badRequest(c, "retentionHours must be between 1 and 8760.")
		return
	}

	app := models.App{
		Name:      strings.TrimSpace(req.Name),
		AppID:     uuid.New().String(),
		APIKey:    strings.ReplaceAll(uuid.New().String(), "-", ""),
		Status:    "active",
		CreatedAt: time.Now(),
	}

	_, err := db.Database.Collection("apps").InsertOne(c.Request.Context(), app)
	if err != nil {
		slog.Error("failed to register application", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Unable to register the application.",
		})
		return
	}

	config := models.AppConfig{
		AppID:          app.AppID,
		RetentionHours: retentionHours,
		UpdatedAt:      time.Now(),
	}
	_, err = db.Database.Collection("app_configs").InsertOne(c.Request.Context(), config)
	if err != nil {
		_, _ = db.Database.Collection("apps").DeleteOne(c.Request.Context(), bson.M{"appId": app.AppID})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Unable to create tenant configuration."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"app":    app,
			"config": config,
		},
	})
}

func validRetentionHours(hours int) bool {
	return hours >= 1 && hours <= 8760
}
