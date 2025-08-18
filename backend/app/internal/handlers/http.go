package handlers

import (
	"encoding/json"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/database"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	db     database.OrderRepository
	cache  cache.OrderCache
	logger *logrus.Logger
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewHTTPHandler создает новый HTTP handler
func NewHTTPHandler(db database.OrderRepository, cache cache.OrderCache, logger *logrus.Logger) *HTTPHandler {
	return &HTTPHandler{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

// SetupRoutes настраивает маршруты
func (h *HTTPHandler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// API маршруты
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/orders/{order_uid}", h.GetOrder).Methods("GET")
	api.HandleFunc("/orders", h.GetAllOrders).Methods("GET")
	api.HandleFunc("/cache/stats", h.GetCacheStats).Methods("GET")
	api.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Простая главная страница с информацией об API
	r.HandleFunc("/", h.APIInfoPage).Methods("GET")

	// Middleware для логирования
	r.Use(h.loggingMiddleware)
	r.Use(h.corsMiddleware)

	return r
}

// GetOrder возвращает заказ по UID (сначала из кеша, затем из БД)
func (h *HTTPHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if strings.TrimSpace(orderUID) == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "order_uid is required")
		return
	}

	h.logger.WithField("order_uid", orderUID).Debug("Getting order")

	// Сначала проверяем кеш
	if order, found := h.cache.Get(orderUID); found {
		h.logger.WithField("order_uid", orderUID).Debug("Order found in cache")
		h.writeSuccessResponse(w, order)
		return
	}

	// Если не найден в кеше, ищем в БД
	order, err := h.db.GetOrderByUID(orderUID)
	if err != nil {
		h.logger.WithError(err).WithField("order_uid", orderUID).Error("Failed to get order from database")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if order == nil {
		h.logger.WithField("order_uid", orderUID).Debug("Order not found")
		h.writeErrorResponse(w, http.StatusNotFound, "Order not found")
		return
	}

	// Добавляем найденный заказ в кеш
	h.cache.Set(orderUID, order)
	h.logger.WithField("order_uid", orderUID).Debug("Order found in database and added to cache")

	h.writeSuccessResponse(w, order)
}

// GetAllOrders возвращает список всех заказов
func (h *HTTPHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // по умолчанию

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	orders, err := h.db.GetAllOrders(limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get orders from database")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.writeSuccessResponse(w, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
		"limit":  limit,
	})
}

// GetCacheStats возвращает статистику кеша
func (h *HTTPHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	stats := h.cache.GetStats()
	h.writeSuccessResponse(w, stats)
}

// HealthCheck проверяет здоровье сервиса
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Простая проверка подключения к БД
	_, err := h.db.GetAllOrders(1)
	if err != nil {
		h.writeErrorResponse(w, http.StatusServiceUnavailable, "Database is unavailable")
		return
	}

	h.writeSuccessResponse(w, map[string]interface{}{
		"status":    "ok",
		"timestamp": "2024-01-01T00:00:00Z", // можно использовать time.Now()
		"cache":     h.cache.GetStats(),
	})
}

// APIInfoPage отображает информацию об API
func (h *HTTPHandler) APIInfoPage(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Order Service API</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1000px; margin: 0 auto; }
        .card { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        h1 { color: #333; margin-bottom: 10px; }
        .subtitle { color: #666; font-size: 1.1rem; margin-bottom: 30px; }
        .endpoint { background: #f8f9fa; padding: 15px; border-radius: 5px; margin: 10px 0; border-left: 4px solid #007bff; }
        .method { background: #007bff; color: white; padding: 3px 8px; border-radius: 3px; font-size: 0.8rem; margin-right: 10px; }
        .path { font-family: monospace; font-weight: bold; }
        .description { color: #666; margin-top: 5px; }
        .example { background: #e9ecef; padding: 10px; border-radius: 3px; font-family: monospace; margin-top: 5px; }
        .frontend-link { background: #28a745; color: white; padding: 15px 25px; text-decoration: none; border-radius: 5px; display: inline-block; margin: 20px 0; }
        .frontend-link:hover { background: #218838; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <h1>🚀 Order Service API</h1>
            <p class="subtitle">RESTful API для работы с заказами</p>
            
            <a href="http://localhost:3000" class="frontend-link">📱 Открыть веб-интерфейс</a>
            
            <h2>📡 Доступные эндпоинты:</h2>
            
            <div class="endpoint">
                <div><span class="method">GET</span><span class="path">/api/v1/orders/{order_uid}</span></div>
                <div class="description">Получить информацию о заказе по ID</div>
                <div class="example">curl http://localhost:8080/api/v1/orders/b563feb7b2b84b6test</div>
            </div>
            
            <div class="endpoint">
                <div><span class="method">GET</span><span class="path">/api/v1/orders?limit=10</span></div>
                <div class="description">Получить список заказов с ограничением</div>
                <div class="example">curl http://localhost:8080/api/v1/orders?limit=5</div>
            </div>
            
            <div class="endpoint">
                <div><span class="method">GET</span><span class="path">/api/v1/cache/stats</span></div>
                <div class="description">Получить статистику кеша</div>
                <div class="example">curl http://localhost:8080/api/v1/cache/stats</div>
            </div>
            
            <div class="endpoint">
                <div><span class="method">GET</span><span class="path">/api/v1/health</span></div>
                <div class="description">Проверить состояние сервиса</div>
                <div class="example">curl http://localhost:8080/api/v1/health</div>
            </div>
            
            <h2>📄 Документация:</h2>
            <ul>
                <li>Все ответы возвращаются в JSON формате</li>
                <li>Успешные ответы: <code>{"success": true, "data": {...}}</code></li>
                <li>Ошибки: <code>{"success": false, "error": "описание"}</code></li>
                <li>CORS включен для всех origins</li>
                <li>Кеш сначала проверяется в памяти, затем в БД</li>
            </ul>
            
            <h2>🛠️ Технологии:</h2>
            <ul>
                <li><strong>Backend:</strong> Go, Gorilla Mux, PostgreSQL, Kafka</li>
                <li><strong>Frontend:</strong> HTML, CSS, JavaScript (отдельный сервер)</li>
                <li><strong>Кеш:</strong> In-memory LRU кеш</li>
                <li><strong>База данных:</strong> PostgreSQL с нормализованной схемой</li>
            </ul>
        </div>
    </div>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(tmpl))
}

// Вспомогательные методы для ответов
func (h *HTTPHandler) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	response := APIResponse{
		Success: true,
		Data:    data,
	}
	h.writeJSONResponse(w, http.StatusOK, response)
}

func (h *HTTPHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, statusCode, response)
}

func (h *HTTPHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Middleware для логирования запросов
func (h *HTTPHandler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Обертка для захвата статус кода
		wrapper := &responseWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapper, r)
		
		h.logger.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": wrapper.statusCode,
			"duration":    time.Since(start),
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
		}).Info("HTTP request")
	})
}

// Middleware для CORS
func (h *HTTPHandler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
