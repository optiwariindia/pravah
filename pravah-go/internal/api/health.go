package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/core"
	"github.com/sampledge-sse/pravah-go/internal/db"
)

func HealthHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mongoStatus := "connected"
	if !db.Ping(ctx) {
		mongoStatus = "disconnected"
	} else {
		db.EnsureIndexes()
	}

	cm := core.GetClientManager()

	status := "ok"
	httpStatus := http.StatusOK
	if mongoStatus == "disconnected" {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status":       status,
		"liveDelivery": true,
		"replay":       mongoStatus == "connected",
		"database":     mongoStatus,
		"connections":  cm.GetConnectionCount(),
		"uptime":       time.Since(core.StartTime).String(),
	})
}
