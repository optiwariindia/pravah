package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/optiwariindia/pravah/pravah-go/internal/core"
	"github.com/optiwariindia/pravah/pravah-go/internal/db"
)

func HealthHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mongoStatus := "connected"
	err := db.Client.Ping(ctx, nil)
	if err != nil {
		mongoStatus = "disconnected"
	}

	cm := core.GetClientManager()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"database": mongoStatus,
		"connections": cm.GetConnectionCount(),
		"uptime": time.Since(core.StartTime).String(),
	})
}
