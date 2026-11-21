package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/api"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/middleware"
	"github.com/sampledge-sse/pravah-go/internal/utils"
)

func main() {
	// Validate environment variables
	if err := utils.CheckEnvironmentVariables(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	mongoURI := os.Getenv("MONGODB")

	// Connect to Database
	db.Connect(mongoURI)

	// Initialize Gin
	r := gin.New()
	r.Use(
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s [GIN] %3d | %13v | %s | %-7s %s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Request.URL.Path,
			)
		}),
		gin.Recovery(),
	)

	// Middleware
	corsOrigins := utils.CORSAllowedOrigins()
	if len(corsOrigins) > 0 {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     corsOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-KEY", "X-ADMIN-KEY"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	// Static files
	r.Static("/public", "./public")
	r.StaticFile("/test.html", "./public/test.html")
	r.StaticFile("/register.html", "./public/register.html")
	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/fonts", "./public/fonts")
	r.Static("/images", "./public/images")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/public/test.html")
	})

	// Routes
	pravah := r.Group("/pravah")
	{
		pravah.GET("/init", api.HandshakeHandler)
		pravah.GET("/health", api.HealthHandler)
		pravah.GET("/subscribe", api.SubscribeHandler)
		pravah.POST("/publish", middleware.ValidatePublishAuth(), api.PublishHandler)
	}

	apiV1 := r.Group("/api/v1")
	{
		apiV1.PUT("/register", middleware.ValidateAdminAuth(), api.RegisterAppHandler)
		apiV1.POST("/subscriptions/token", middleware.ValidatePublishAuth(), api.CreateSubscriptionTokenHandler)
		apiV1.GET("/config", middleware.ValidatePublishAuth(), api.GetAppConfigHandler)
		apiV1.PUT("/config", middleware.ValidatePublishAuth(), api.UpdateAppConfigHandler)
	}

	log.Printf("Server started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
