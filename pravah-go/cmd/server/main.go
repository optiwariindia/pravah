package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/optiwariindia/pravah/pravah-go/internal/api"
	"github.com/optiwariindia/pravah/pravah-go/internal/core"
	"github.com/optiwariindia/pravah/pravah-go/internal/db"
	"github.com/optiwariindia/pravah/pravah-go/internal/middleware"
	"github.com/optiwariindia/pravah/pravah-go/internal/utils"
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
	r := gin.Default()

	// Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-KEY"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		apiV1.PUT("/register", api.RegisterAppHandler)
	}

	// Start background broadcast for timestamp
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		cm := core.GetClientManager()
		for {
			select {
			case t := <-ticker.C:
				cm.Broadcast("timestamp", map[string]interface{}{
					"server_time": t.Format(time.RFC3339),
				})
			}
		}
	}()

	log.Printf("Server started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
