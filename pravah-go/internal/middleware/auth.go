package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/optiwariindia/pravah/pravah-go/internal/db"
	"github.com/optiwariindia/pravah/pravah-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidatePublishAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mode 1: Cookie-based
		pravahID, err := c.Cookie("pravah_id")
		if err == nil && pravahID != "" {
			c.Next()
			return
		}

		// Mode 2: API Key-based
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != "" {
			var app models.App
			err := db.Database.Collection("apps").FindOne(context.TODO(), bson.M{"apiKey": apiKey, "status": "active"}).Decode(&app)
			if err == nil {
				c.Set("appContext", app)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "Unauthorized source. Missing pravah_id cookie or valid X-API-KEY.",
		})
		c.Abort()
	}
}
