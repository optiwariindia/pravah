package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type cachedApp struct {
	app       models.App
	expiresAt time.Time
}

var publishAuthCache sync.Map

func ValidatePublishAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != "" {
			cacheKey := hashAPIKey(apiKey)
			var app models.App
			err := db.Database.Collection("apps").FindOne(c.Request.Context(), bson.M{"apiKey": apiKey, "status": "active"}).Decode(&app)
			if err == nil {
				publishAuthCache.Store(cacheKey, cachedApp{
					app:       app,
					expiresAt: time.Now().Add(5 * time.Minute),
				})
				c.Set("appContext", app)
				c.Next()
				return
			}
			if err != mongo.ErrNoDocuments && db.Database != nil {
				db.MarkUnhealthy()
				if cached, ok := publishAuthCache.Load(cacheKey); ok {
					entry := cached.(cachedApp)
					if time.Now().Before(entry.expiresAt) {
						c.Header("X-Pravah-Mode", "degraded")
						c.Set("appContext", entry.app)
						c.Next()
						return
					}
					publishAuthCache.Delete(cacheKey)
				}
			}
			if err == mongo.ErrNoDocuments {
				publishAuthCache.Delete(cacheKey)
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "A valid X-API-KEY is required.",
		})
		c.Abort()
	}
}

func hashAPIKey(apiKey string) string {
	sum := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(sum[:])
}

func ValidateAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := os.Getenv("ADMIN_API_KEY")
		if len(expected) >= 32 && hmac.Equal([]byte(c.GetHeader("X-ADMIN-KEY")), []byte(expected)) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "A valid X-ADMIN-KEY is required.",
		})
		c.Abort()
	}
}
