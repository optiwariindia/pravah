package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/api"
	"github.com/sampledge-sse/pravah-go/internal/core"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/middleware"
	"github.com/sampledge-sse/pravah-go/internal/utils"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "check the local Pravah health endpoint")
	flag.Parse()
	if *healthcheck {
		os.Exit(runHealthcheck())
	}

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
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownContext, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to run server: %v", err)
		}
	case <-shutdownContext.Done():
		log.Println("Shutdown requested; closing SSE clients")
		core.GetClientManager().CloseAll("server_shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP shutdown did not complete gracefully: %v", err)
		}
		if db.Client != nil {
			disconnectContext, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer disconnectCancel()
			if err := db.Client.Disconnect(disconnectContext); err != nil {
				log.Printf("MongoDB disconnect failed: %v", err)
			}
		}
	}
}

func runHealthcheck() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Get("http://127.0.0.1:" + port + "/pravah/health")
	if err != nil {
		log.Printf("Healthcheck request failed: %v", err)
		return 1
	}
	defer response.Body.Close()
	return map[bool]int{true: 0, false: 1}[response.StatusCode == http.StatusOK]
}
