package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sampledge-sse/pravah-go/internal/core"
	"github.com/sampledge-sse/pravah-go/internal/db"
	"github.com/sampledge-sse/pravah-go/internal/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoIntegration(t *testing.T) {
	uri := os.Getenv("PRAVAH_MONGODB_TEST_URI")
	if uri == "" {
		t.Skip("set PRAVAH_MONGODB_TEST_URI to run MongoDB integration tests")
	}

	t.Setenv("SUBSCRIPTION_JWT_SECRET", "0123456789abcdef0123456789abcdef0123456789abcdef")
	t.Setenv("ADMIN_API_KEY", "abcdef0123456789abcdef0123456789abcdef0123456789")

	gin.SetMode(gin.TestMode)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("connect test MongoDB: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Disconnect(context.Background())
	})

	database := client.Database(fmt.Sprintf("pravah_test_%d", time.Now().UnixNano()))
	db.Client = client
	db.Database = database
	db.Ping(context.Background())
	t.Cleanup(func() {
		_ = database.Drop(context.Background())
	})
	db.EnsureIndexes()

	router := gin.New()
	router.GET("/pravah/init", HandshakeHandler)
	router.GET("/pravah/health", HealthHandler)
	router.GET("/pravah/subscribe", SubscribeHandler)
	router.POST("/pravah/publish", middleware.ValidatePublishAuth(), PublishHandler)
	router.PUT("/api/v1/register", middleware.ValidateAdminAuth(), RegisterAppHandler)
	router.POST("/api/v1/subscriptions/token", middleware.ValidatePublishAuth(), CreateSubscriptionTokenHandler)

	appOne := registerIntegrationApp(t, router, "integration-one")
	appTwo := registerIntegrationApp(t, router, "integration-two")
	tokenOne := issueIntegrationToken(t, router, appOne.apiKey, []string{"project:integration"})
	tokenTwo := issueIntegrationToken(t, router, appTwo.apiKey, []string{"project:integration"})
	_ = tokenTwo

	invalidPublish := integrationRequest(t, router, http.MethodPost, "/pravah/publish", "invalid-api-key", nil,
		strings.NewReader(`{"streamId":"project:integration","event":"invalid","data":{"id":"1"}}`),
	)
	if invalidPublish.Code != http.StatusForbidden {
		t.Fatalf("expected invalid API key rejection, got %d: %s", invalidPublish.Code, invalidPublish.Body.String())
	}

	handshake := integrationRequest(t, router, http.MethodGet, "/pravah/init", "", nil)
	cookie := handshake.Result().Cookies()[0]

	unauthorized := integrationRequest(t, router, http.MethodGet,
		"/pravah/subscribe?streams=private:stream&token="+tokenOne,
		"", []*http.Cookie{cookie},
	)
	if unauthorized.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden subscription, got %d: %s", unauthorized.Code, unauthorized.Body.String())
	}

	oversized := strings.NewReader(`{"streamId":"project:integration","event":"invalid","data":{"id":"1","metadata":{"secret":true}}}`)
	rejected := integrationRequest(t, router, http.MethodPost, "/pravah/publish", appOne.apiKey, nil, oversized)
	if rejected.Code != http.StatusBadRequest {
		t.Fatalf("expected metadata rejection, got %d: %s", rejected.Code, rejected.Body.String())
	}

	first := publishIntegrationSignal(t, router, appOne.apiKey, "first-event", "resource-1")

	second := publishIntegrationSignal(t, router, appOne.apiKey, "second-event", "resource-2")
	replay := startIntegrationSubscription(t, router, cookie, tokenOne, "project:integration", first.eventID)
	waitForIntegrationText(t, replay, "second-event", 3*time.Second)
	if !strings.Contains(replay.String(), second.eventID) {
		t.Fatalf("replay stream is missing event ID %s: %s", second.eventID, replay.String())
	}
	replay.cancel()
	waitForIntegrationCondition(t, 3*time.Second, func() bool {
		return core.GetClientManager().GetConnectionCount() == 0
	}, "client registry cleanup")

	tenantTwo := startIntegrationSubscription(t, router, cookie, tokenTwo, "project:integration", first.eventID)
	tenantOne := publishIntegrationSignal(t, router, appOne.apiKey, "tenant-one-only", "resource-3")
	_ = tenantOne
	time.Sleep(100 * time.Millisecond)
	if strings.Contains(tenantTwo.String(), "tenant-one-only") {
		t.Fatal("tenant two received tenant one's replay/live signal")
	}
	tenantTwo.cancel()
	waitForIntegrationCondition(t, 3*time.Second, func() bool {
		return core.GetClientManager().GetConnectionCount() == 0
	}, "tenant client registry cleanup")

	health := integrationRequest(t, router, http.MethodGet, "/pravah/health", "", nil)
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), `"status":"ok"`) {
		t.Fatalf("unexpected health response %d: %s", health.Code, health.Body.String())
	}

	eventCount, err := database.Collection("events").CountDocuments(context.Background(), bson.M{"appId": appOne.appID})
	if err != nil {
		t.Fatalf("count tenant events: %v", err)
	}
	if eventCount != 3 {
		t.Fatalf("expected 3 tenant events, got %d", eventCount)
	}
}

type integrationApp struct {
	appID  string
	apiKey string
}

type integrationEvent struct {
	eventID string
}

func registerIntegrationApp(t *testing.T, router http.Handler, name string) integrationApp {
	t.Helper()
	response := integrationRequest(t, router, http.MethodPut, "/api/v1/register", "", nil,
		strings.NewReader(fmt.Sprintf(`{"name":%q}`, name)),
	)
	if response.Code != http.StatusCreated {
		t.Fatalf("register app: %d: %s", response.Code, response.Body.String())
	}
	var payload struct {
		Data struct {
			App struct {
				AppID  string `json:"appId"`
				APIKey string `json:"apiKey"`
			} `json:"app"`
		} `json:"data"`
	}
	decodeIntegrationJSON(t, response, &payload)
	return integrationApp{appID: payload.Data.App.AppID, apiKey: payload.Data.App.APIKey}
}

func issueIntegrationToken(t *testing.T, router http.Handler, apiKey string, streams []string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{"streams": streams, "expiresIn": 300})
	response := integrationRequest(t, router, http.MethodPost, "/api/v1/subscriptions/token", apiKey, nil, bytes.NewReader(body))
	if response.Code != http.StatusCreated {
		t.Fatalf("issue token: %d: %s", response.Code, response.Body.String())
	}
	var payload struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	decodeIntegrationJSON(t, response, &payload)
	return payload.Data.Token
}

func publishIntegrationSignal(t *testing.T, router http.Handler, apiKey, event, resource string) integrationEvent {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{
		"streamId": "project:integration",
		"event":    event,
		"data":     map[string]string{"id": resource},
	})
	response := integrationRequest(t, router, http.MethodPost, "/pravah/publish", apiKey, nil, bytes.NewReader(body))
	if response.Code != http.StatusAccepted {
		t.Fatalf("publish signal: %d: %s", response.Code, response.Body.String())
	}
	var payload struct {
		EventID string `json:"eventId"`
	}
	decodeIntegrationJSON(t, response, &payload)
	return integrationEvent{eventID: payload.EventID}
}

func integrationRequest(
	t *testing.T,
	router http.Handler,
	method string,
	target string,
	apiKey string,
	cookies []*http.Cookie,
	bodies ...io.Reader,
) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if len(bodies) > 0 {
		body = bodies[0]
	}
	request := httptest.NewRequest(method, target, body)
	request.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		request.Header.Set("X-API-KEY", apiKey)
	}
	if strings.Contains(target, "/api/v1/register") {
		request.Header.Set("X-ADMIN-KEY", os.Getenv("ADMIN_API_KEY"))
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

type integrationSubscription struct {
	mu     sync.Mutex
	body   strings.Builder
	cancel context.CancelFunc
}

func startIntegrationSubscription(
	t *testing.T,
	router http.Handler,
	cookie *http.Cookie,
	token string,
	stream string,
	lastEventID string,
) *integrationSubscription {
	t.Helper()
	target := "/pravah/subscribe?streams=" + stream + "&token=" + token
	ctx, cancel := context.WithCancel(context.Background())
	request := httptest.NewRequest(http.MethodGet, target, nil).WithContext(ctx)
	request.Header.Set("Last-Event-ID", lastEventID)
	request.AddCookie(cookie)
	subscription := &integrationSubscription{cancel: cancel}
	writer := &integrationResponseWriter{
		header:       make(http.Header),
		subscription: subscription,
	}

	go func() {
		router.ServeHTTP(writer, request)
	}()
	t.Cleanup(cancel)
	waitForIntegrationConnection(t, subscription)
	return subscription
}

type integrationResponseWriter struct {
	header       http.Header
	subscription *integrationSubscription
	status       int
}

func (w *integrationResponseWriter) Header() http.Header {
	return w.header
}

func (w *integrationResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (w *integrationResponseWriter) Write(data []byte) (int, error) {
	w.subscription.mu.Lock()
	defer w.subscription.mu.Unlock()
	return w.subscription.body.Write(data)
}

func (w *integrationResponseWriter) Flush() {}

func (w *integrationResponseWriter) CloseNotify() <-chan bool {
	channel := make(chan bool)
	return channel
}

func (s *integrationSubscription) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.body.String()
}

func waitForIntegrationConnection(t *testing.T, subscription *integrationSubscription) {
	t.Helper()
	waitForIntegrationText(t, subscription, "event: connected", 3*time.Second)
}

func waitForIntegrationText(t *testing.T, subscription *integrationSubscription, text string, timeout time.Duration) {
	t.Helper()
	waitForIntegrationCondition(t, timeout, func() bool {
		return strings.Contains(subscription.String(), text)
	}, fmt.Sprintf("%q in subscription output: %s", text, subscription.String()))
}

func waitForIntegrationCondition(t *testing.T, timeout time.Duration, condition func() bool, description string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", description)
}

func decodeIntegrationJSON(t *testing.T, response *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response %q: %v", response.Body.String(), err)
	}
}
