package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	API_BASE_URL = "http://localhost:8081/api/v1"
	TEST_ORDER_ID = "b563feb7b2b84b6test"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type BenchmarkResult struct {
	RequestNum int
	Duration   time.Duration
	Success    bool
	Source     string // "cache" или "database"
}

func main() {
	fmt.Println("🚀 Тестирование производительности кеша Order Service")
	fmt.Println("=================================================")
	
	// Проверяем доступность API
	if !checkAPIHealth() {
		fmt.Println("❌ API недоступен. Убедитесь, что сервис запущен на порту 8081")
		return
	}
	
	fmt.Println("✅ API доступен, начинаем тестирование...")
	fmt.Println()

	// Тест 1: Первый запрос (из БД)
	fmt.Println("📊 Тест 1: Первый запрос (загрузка из БД)")
	firstRequest := benchmarkSingleRequest(TEST_ORDER_ID, 1)
	fmt.Printf("Время первого запроса: %v\n", firstRequest.Duration)
	
	if !firstRequest.Success {
		fmt.Println("❌ Первый запрос неуспешен, завершаем тестирование")
		return
	}
	
	time.Sleep(100 * time.Millisecond) // Небольшая пауза
	
	// Тест 2: Повторные запросы (из кеша)
	fmt.Println("\n📊 Тест 2: Повторные запросы (из кеша)")
	cacheRequests := make([]BenchmarkResult, 0, 10)
	
	for i := 2; i <= 11; i++ {
		result := benchmarkSingleRequest(TEST_ORDER_ID, i)
		cacheRequests = append(cacheRequests, result)
		fmt.Printf("Запрос %d: %v\n", i, result.Duration)
		time.Sleep(10 * time.Millisecond) // Небольшая пауза между запросами
	}
	
	// Тест 3: Параллельные запросы
	fmt.Println("\n📊 Тест 3: Параллельные запросы (50 одновременных)")
	parallelResults := benchmarkParallelRequests(TEST_ORDER_ID, 50)
	
	// Анализ результатов
	fmt.Println("\n📈 Анализ результатов:")
	fmt.Println("=====================")
	
	// Средние времена
	avgCacheTime := calculateAverageTime(cacheRequests)
	avgParallelTime := calculateAverageTime(parallelResults)
	
	fmt.Printf("Первый запрос (БД):           %v\n", firstRequest.Duration)
	fmt.Printf("Среднее время кеша:          %v\n", avgCacheTime)
	fmt.Printf("Среднее время параллельно:   %v\n", avgParallelTime)
	
	// Ускорение кеша
	if avgCacheTime > 0 {
		speedup := float64(firstRequest.Duration) / float64(avgCacheTime)
		fmt.Printf("Ускорение кеша:              %.2fx\n", speedup)
	}
	
	// Тест 4: Статистика кеша
	fmt.Println("\n📊 Тест 4: Статистика кеша")
	cacheStats, err := getCacheStats()
	if err != nil {
		fmt.Printf("❌ Ошибка получения статистики кеша: %v\n", err)
	} else {
		fmt.Printf("Размер кеша: %v\n", cacheStats["size"])
		fmt.Printf("Максимальный размер: %v\n", cacheStats["capacity"])
		if size, ok := cacheStats["size"].(float64); ok {
			if capacity, ok := cacheStats["capacity"].(float64); ok {
				usage := (size / capacity) * 100
				fmt.Printf("Использование: %.1f%%\n", usage)
			}
		}
	}
	
	// Тест 5: Нагрузочное тестирование
	fmt.Println("\n📊 Тест 5: Нагрузочное тестирование (100 запросов)")
	loadTestResults := performLoadTest(TEST_ORDER_ID, 100)
	
	successCount := 0
	totalTime := time.Duration(0)
	for _, result := range loadTestResults {
		if result.Success {
			successCount++
			totalTime += result.Duration
		}
	}
	
	if successCount > 0 {
		avgLoadTime := totalTime / time.Duration(successCount)
		successRate := float64(successCount) / float64(len(loadTestResults)) * 100
		
		fmt.Printf("Успешных запросов: %d/%d (%.1f%%)\n", successCount, len(loadTestResults), successRate)
		fmt.Printf("Среднее время: %v\n", avgLoadTime)
		fmt.Printf("RPS: %.1f\n", 1000.0/float64(avgLoadTime.Milliseconds()))
	}
	
	fmt.Println("\n✅ Тестирование завершено!")
}

func checkAPIHealth() bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(API_BASE_URL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

func benchmarkSingleRequest(orderID string, requestNum int) BenchmarkResult {
	start := time.Now()
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/orders/%s", API_BASE_URL, orderID))
	
	duration := time.Since(start)
	
	if err != nil {
		return BenchmarkResult{
			RequestNum: requestNum,
			Duration:   duration,
			Success:    false,
		}
	}
	defer resp.Body.Close()
	
	success := resp.StatusCode == http.StatusOK
	
	return BenchmarkResult{
		RequestNum: requestNum,
		Duration:   duration,
		Success:    success,
	}
}

func benchmarkParallelRequests(orderID string, count int) []BenchmarkResult {
	results := make([]BenchmarkResult, count)
	var wg sync.WaitGroup
	
	start := time.Now()
	
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			requestStart := time.Now()
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get(fmt.Sprintf("%s/orders/%s", API_BASE_URL, orderID))
			requestDuration := time.Since(requestStart)
			
			success := err == nil && resp != nil && resp.StatusCode == http.StatusOK
			if resp != nil {
				resp.Body.Close()
			}
			
			results[index] = BenchmarkResult{
				RequestNum: index + 1,
				Duration:   requestDuration,
				Success:    success,
			}
		}(i)
	}
	
	wg.Wait()
	totalDuration := time.Since(start)
	
	fmt.Printf("Все %d запросов выполнены за %v\n", count, totalDuration)
	
	return results
}

func performLoadTest(orderID string, count int) []BenchmarkResult {
	results := make([]BenchmarkResult, count)
	
	start := time.Now()
	for i := 0; i < count; i++ {
		results[i] = benchmarkSingleRequest(orderID, i+1)
		
		// Небольшая пауза для имитации реальной нагрузки
		if i%10 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	totalDuration := time.Since(start)
	
	fmt.Printf("Нагрузочный тест завершен за %v\n", totalDuration)
	
	return results
}

func calculateAverageTime(results []BenchmarkResult) time.Duration {
	if len(results) == 0 {
		return 0
	}
	
	total := time.Duration(0)
	count := 0
	
	for _, result := range results {
		if result.Success {
			total += result.Duration
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return total / time.Duration(count)
}

func getCacheStats() (map[string]interface{}, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(API_BASE_URL + "/cache/stats")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	
	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}
	
	if stats, ok := apiResp.Data.(map[string]interface{}); ok {
		return stats, nil
	}
	
	return nil, fmt.Errorf("unexpected response format")
}
