package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func TestCORSAllowsConfiguredOriginOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(cors.New(corsConfig([]string{"https://app.example.com"})))
	router.OPTIONS("/pravah/publish", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	allowed := preflightRequest(router, "https://app.example.com")
	if allowed.Code != http.StatusNoContent {
		t.Fatalf("expected allowed preflight to return 204, got %d", allowed.Code)
	}
	if allowed.Header().Get("Access-Control-Allow-Origin") != "https://app.example.com" {
		t.Fatalf("unexpected allow-origin header: %q", allowed.Header().Get("Access-Control-Allow-Origin"))
	}

	blocked := preflightRequest(router, "https://attacker.example")
	if blocked.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("unexpected allow-origin header for blocked origin: %q", blocked.Header().Get("Access-Control-Allow-Origin"))
	}
}

func preflightRequest(handler http.Handler, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodOptions, "/pravah/publish", nil)
	request.Header.Set("Origin", origin)
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}
