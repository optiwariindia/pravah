package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type registeredApp struct {
	Data struct {
		App struct {
			APIKey string `json:"apiKey"`
		} `json:"app"`
	} `json:"data"`
}

type tokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type publishResponse struct {
	EventID string `json:"eventId"`
}

func main() {
	baseURL := requiredEnv("PRAVAH_URL")
	adminKey := requiredEnv("ADMIN_API_KEY")
	connections := positiveEnv("LOAD_CONNECTIONS", 100)
	events := positiveEnv("LOAD_EVENTS", 50)
	rate := positiveEnv("LOAD_RATE_PER_SECOND", 50)
	maxLatency := time.Duration(positiveEnv("LOAD_MAX_LATENCY_MS", 1000)) * time.Millisecond

	client := &http.Client{}
	app := requestJSON[registeredApp](client, http.MethodPut, baseURL+"/api/v1/register", adminKey, "", map[string]any{
		"name":           fmt.Sprintf("load-test-%d", time.Now().UnixNano()),
		"retentionHours": 1,
	})
	apiKey := app.Data.App.APIKey
	token := requestJSON[tokenResponse](client, http.MethodPost, baseURL+"/api/v1/subscriptions/token", "", apiKey, map[string]any{
		"streams":   []string{"load:stream"},
		"expiresIn": 3600,
	}).Data.Token
	cookie := initializeCookie(client, baseURL)

	var delivered atomic.Int64
	connected := make(chan struct{}, connections)
	observer := make(chan string, events)
	testContext, cancelTest := context.WithCancel(context.Background())
	defer cancelTest()
	var clients sync.WaitGroup

	for index := 0; index < connections; index++ {
		clients.Add(1)
		go func(observerClient bool) {
			defer clients.Done()
			subscribe(testContext, client, baseURL, token, cookie, observerClient, connected, observer, &delivered)
		}(index == 0)
	}
	awaitConnections(connected, connections)

	started := make(map[string]time.Time, events)
	var startedMu sync.Mutex
	interval := time.Second / time.Duration(rate)
	testStarted := time.Now()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for index := 0; index < events; index++ {
		resourceID := fmt.Sprintf("resource-%03d", index)
		startedMu.Lock()
		started[resourceID] = time.Now()
		startedMu.Unlock()

		requestJSON[publishResponse](client, http.MethodPost, baseURL+"/pravah/publish", "", apiKey, map[string]any{
			"streamId": "load:stream",
			"event":    "load-event",
			"data":     map[string]string{"id": resourceID},
		})
		if index+1 < events {
			<-ticker.C
		}
	}
	publishDuration := time.Since(testStarted)

	observed := make(map[string]time.Duration, events)
	deadline := time.After(10 * time.Second)
	for len(observed) < events {
		select {
		case resourceID := <-observer:
			startedMu.Lock()
			startedAt := started[resourceID]
			startedMu.Unlock()
			observed[resourceID] = time.Since(startedAt)
		case <-deadline:
			log.Fatalf("observer received only %d/%d events", len(observed), events)
		}
	}

	expected := int64(connections * events)
	for delivered.Load() < expected && time.Now().Before(testStarted.Add(15*time.Second)) {
		time.Sleep(20 * time.Millisecond)
	}
	cancelTest()

	latencies := make([]time.Duration, 0, len(observed))
	for _, latency := range observed {
		latencies = append(latencies, latency)
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	p95 := latencies[(len(latencies)*95-1)/100]
	throughput := float64(events) / publishDuration.Seconds()

	fmt.Printf("connections=%d events=%d delivered=%d/%d throughput=%.2f/s p95=%s\n",
		connections, events, delivered.Load(), expected, throughput, p95.Round(time.Millisecond))
	if delivered.Load() != expected {
		log.Fatalf("event drops detected")
	}
	if p95 > maxLatency {
		log.Fatalf("p95 latency %s exceeds %s", p95, maxLatency)
	}
	if throughput < float64(rate)*0.9 {
		log.Fatalf("throughput %.2f/s is below 90%% of target %d/s", throughput, rate)
	}
	clients.Wait()
}

func subscribe(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	token string,
	cookie *http.Cookie,
	observer bool,
	connected chan<- struct{},
	observerEvents chan<- string,
	delivered *atomic.Int64,
) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/pravah/subscribe?streams=load:stream&token="+token, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.AddCookie(cookie)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	eventName := ""
	for scanner.Scan() {
		line := scanner.Text()
		if stringsHasPrefix(line, "event: ") {
			eventName = line[len("event: "):]
		} else if stringsHasPrefix(line, "data: ") && eventName == "connected" {
			connected <- struct{}{}
			eventName = ""
		} else if stringsHasPrefix(line, "data: ") && eventName == "signal" {
			var signal struct {
				Data struct {
					ID string `json:"id"`
				} `json:"data"`
			}
			if json.Unmarshal([]byte(line[len("data: "):]), &signal) == nil {
				delivered.Add(1)
				if observer {
					select {
					case observerEvents <- signal.Data.ID:
					case <-ctx.Done():
						return
					}
				}
			}
			eventName = ""
		}
	}
	if err := scanner.Err(); err != nil {
		select {
		case <-ctx.Done():
		default:
			log.Printf("SSE scanner failed: %v", err)
		}
	}
}

func stringsHasPrefix(value, prefix string) bool {
	return len(value) >= len(prefix) && value[:len(prefix)] == prefix
}

func initializeCookie(client *http.Client, baseURL string) *http.Cookie {
	response, err := client.Get(baseURL + "/pravah/init")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	cookies := response.Cookies()
	if len(cookies) == 0 {
		log.Fatal("handshake did not set a cookie")
	}
	return cookies[0]
}

func awaitConnections(connected <-chan struct{}, expected int) {
	deadline := time.After(10 * time.Second)
	for count := 0; count < expected; count++ {
		select {
		case <-connected:
		case <-deadline:
			log.Fatalf("only %d/%d SSE clients connected", count, expected)
		}
	}
}

func requestJSON[T any](client *http.Client, method, url, adminKey, apiKey string, payload any) T {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	if adminKey != "" {
		request.Header.Set("X-ADMIN-KEY", adminKey)
	}
	if apiKey != "" {
		request.Header.Set("X-API-KEY", apiKey)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Fatalf("%s %s: %d: %s", method, url, response.StatusCode, data)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalf("decode %s response: %v", url, err)
	}
	return result
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s is required", name)
	}
	return value
}

func positiveEnv(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		log.Fatalf("%s must be a positive integer", name)
	}
	return number
}
