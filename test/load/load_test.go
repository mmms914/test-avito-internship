//go:build load
// +build load

package load

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

type TestResult struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	TotalDuration   time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
	Durations       []time.Duration
	Errors          map[string]int64
	mu              sync.Mutex
}

type LoadTestConfig struct {
	URL         string
	Duration    time.Duration
	Concurrency int
	TargetRPS   int
	RoomID      string
	AdminToken  string
	UserToken   string
}

type Slot struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}

type BookingResponse struct {
	Booking struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"booking"`
}

func TestLoad(t *testing.T) {
	// 1. Получаем токены через dummyLogin
	adminToken, userToken := getTokens(t)

	// 2. Создаем тестовую переговорку
	roomID := createTestRoom(t, adminToken)

	// 3. Создаем расписание
	createTestSchedule(t, adminToken, roomID)

	// 4. Запускаем нагрузочный тест
	cfg := &LoadTestConfig{
		URL:         "http://localhost:8081",
		Duration:    2 * time.Minute,
		Concurrency: 50,
		TargetRPS:   100,
		RoomID:      roomID,
		AdminToken:  adminToken,
		UserToken:   userToken,
	}

	result := RunLoadTest(t, cfg)

	// 5. Выводим результаты
	printResults(result)

	// 6. Проверяем требования
	if float64(result.SuccessRequests)/float64(result.TotalRequests) < 0.999 {
		t.Errorf("Success rate below 99.9%%: %.2f%%",
			float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
	}

	avgDuration := time.Duration(result.TotalDuration.Nanoseconds() / result.SuccessRequests)
	if avgDuration > 200*time.Millisecond {
		t.Errorf("Average response time > 200ms: %v", avgDuration)
	}
}

func getTokens(t *testing.T) (adminToken, userToken string) {
	client := &http.Client{Timeout: 10 * time.Second}

	adminToken = getDummyLoginToken(t, client, "admin")
	userToken = getDummyLoginToken(t, client, "user")

	return adminToken, userToken
}

func getDummyLoginToken(t *testing.T, client *http.Client, role string) string {
	reqBody := map[string]string{"role": role}
	body, _ := json.Marshal(reqBody)

	resp, err := client.Post("http://localhost:8081/dummyLogin", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode token response: %v", err)
	}

	return result.Token
}

func createTestRoom(t *testing.T, adminToken string) string {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := map[string]interface{}{
		"name":        fmt.Sprintf("Load Test Room %s", uuid.New().String()[:8]),
		"description": "Room for load testing",
		"capacity":    10,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "http://localhost:8081/rooms/create", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode room response: %v", err)
	}

	return result.Room.ID
}

func createTestSchedule(t *testing.T, adminToken, roomID string) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := map[string]interface{}{
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "00:00",
		"endTime":    "23:59",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:8081/rooms/%s/schedule/create", roomID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create schedule, status: %d", resp.StatusCode)
	}
}

func RunLoadTest(t *testing.T, cfg *LoadTestConfig) *TestResult {
	result := &TestResult{
		MinDuration: time.Hour,
		Errors:      make(map[string]int64),
		Durations:   make([]time.Duration, 0),
	}

	// Контроль RPS
	ticker := time.NewTicker(time.Second / time.Duration(cfg.TargetRPS))
	defer ticker.Stop()

	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(cfg, result, &wg, stopChan, ticker)
	}

	time.Sleep(cfg.Duration)
	close(stopChan)
	wg.Wait()

	return result
}

func worker(cfg *LoadTestConfig, result *TestResult, wg *sync.WaitGroup, stopChan <-chan struct{}, ticker *time.Ticker) {
	defer wg.Done()

	client := &http.Client{Timeout: 10 * time.Second}

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			// Генерируем случайную дату в будущем (до 5 лет вперед)
			randomDays := 1 + time.Duration(rand.Intn(5*365))*24*time.Hour
			futureDate := time.Now().UTC().Add(randomDays)
			dateStr := futureDate.Format("2006-01-02")

			// Запрос на получение слотов
			start := time.Now()
			err := makeSlotsRequest(client, cfg, dateStr)
			duration := time.Since(start)

			result.mu.Lock()
			result.TotalRequests++
			result.Durations = append(result.Durations, duration)

			if duration < result.MinDuration {
				result.MinDuration = duration
			}
			if duration > result.MaxDuration {
				result.MaxDuration = duration
			}
			result.TotalDuration += duration

			if err != nil {
				result.FailedRequests++
				result.Errors[err.Error()]++
			} else {
				result.SuccessRequests++
			}
			result.mu.Unlock()
		}
	}
}

func makeSlotsRequest(client *http.Client, cfg *LoadTestConfig, dateStr string) error {
	url := fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", cfg.URL, cfg.RoomID, dateStr)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.UserToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var result struct {
		Slots []Slot `json:"slots"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}

func printResults(result *TestResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("LOAD TEST RESULTS")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("Total Requests:     %d\n", result.TotalRequests)
	fmt.Printf("Successful:         %d (%.2f%%)\n",
		result.SuccessRequests,
		float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("Failed:             %d (%.2f%%)\n",
		result.FailedRequests,
		float64(result.FailedRequests)/float64(result.TotalRequests)*100)

	fmt.Println("\nResponse Times:")
	fmt.Printf("  Min:              %v\n", result.MinDuration)
	fmt.Printf("  Max:              %v\n", result.MaxDuration)
	fmt.Printf("  Avg:              %v\n",
		time.Duration(result.TotalDuration.Nanoseconds()/result.SuccessRequests))

	// Перцентили
	durations := result.Durations
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	p50 := durations[len(durations)*50/100]
	p95 := durations[len(durations)*95/100]
	p99 := durations[len(durations)*99/100]

	fmt.Printf("  P50:              %v\n", p50)
	fmt.Printf("  P95:              %v\n", p95)
	fmt.Printf("  P99:              %v\n", p99)

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for err, count := range result.Errors {
			fmt.Printf("  %s: %d times\n", err, count)
		}
	}

	// RPS
	rps := float64(result.TotalRequests) / result.TotalDuration.Seconds()
	fmt.Printf("\nRPS:                %.2f\n", rps)

	fmt.Println(strings.Repeat("=", 60))
}
