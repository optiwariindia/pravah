package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/models"
)

type RegisterRequest struct {
	Name string `json:"name" binding:"required"`
}

func RegisterAppHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Name is required."})
		return
	}

	app := models.App{
		Name:      req.Name,
		AppID:     uuid.New().String(),
		APIKey:    strings.ReplaceAll(uuid.New().String(), "-", ""),
		Status:    "active",
		CreatedAt: time.Now(),
	}

	_, err := db.Database.Collection("apps").InsertOne(context.TODO(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    app,
	})
}
