package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidatePublishAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != "" {
			var app models.App
			err := db.Database.Collection("apps").FindOne(c.Request.Context(), bson.M{"apiKey": apiKey, "status": "active"}).Decode(&app)
			if err == nil {
				c.Set("appContext", app)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "A valid X-API-KEY is required.",
		})
		c.Abort()
	}
}
