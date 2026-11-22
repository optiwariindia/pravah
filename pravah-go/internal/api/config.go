package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/models"
	"github.com/sampledge-sse/pravah-go/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppConfigRequest struct {
	RetentionHours int `json:"retentionHours"`
}

func GetAppConfigHandler(c *gin.Context) {
	app, ok := appFromContext(c)
	if !ok {
		return
	}

	config, err := getAppConfig(c.Request.Context(), app.AppID)
	if err != nil {
		log.Printf("[Pravah] Failed to load application config %s: %v", app.AppID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Unable to load tenant configuration."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": config})
}

func UpdateAppConfigHandler(c *gin.Context) {
	app, ok := appFromContext(c)
	if !ok {
		return
	}

	var req AppConfigRequest
	if err := decodeStrictJSON(c, &req); err != nil || !validRetentionHours(req.RetentionHours) {
		badRequest(c, "retentionHours must be between 1 and 8760.")
		return
	}

	config := models.AppConfig{
		AppID:          app.AppID,
		RetentionHours: req.RetentionHours,
		UpdatedAt:      time.Now(),
	}
	_, err := db.Database.Collection("app_configs").ReplaceOne(
		c.Request.Context(),
		bson.M{"appId": app.AppID},
		config,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		log.Printf("[Pravah] Failed to update application config %s: %v", app.AppID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Unable to update tenant configuration."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": config})
}

func getAppConfig(ctx context.Context, appID string) (models.AppConfig, error) {
	var config models.AppConfig
	err := db.Database.Collection("app_configs").FindOne(ctx, bson.M{"appId": appID}).Decode(&config)
	if err == nil {
		return config, nil
	}
	if err == mongo.ErrNoDocuments {
		return models.AppConfig{
			AppID:          appID,
			RetentionHours: utils.EventRetentionHours(),
		}, nil
	}
	return models.AppConfig{}, err
}

func appFromContext(c *gin.Context) (models.App, bool) {
	value, exists := c.Get("appContext")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden", "message": "A valid X-API-KEY is required."})
		return models.App{}, false
	}
	app, ok := value.(models.App)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Invalid application context."})
		return models.App{}, false
	}
	return app, true
}
