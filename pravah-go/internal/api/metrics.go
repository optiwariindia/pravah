package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/metrics"
)

func MetricsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"metrics": metrics.GetSnapshot()})
}
