package main

import (
	"context"
	"fmt"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/handlers"
	"order-service/internal/kafka"
	"order-service/pkg/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Настройка логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Если установлена переменная DEBUG, включаем отладочные логи
	if os.Getenv("DEBUG") == "true" {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.Info("Starting Order Service")

	// Загружаем конфигурацию
	cfg := config.LoadConfig()
	logger.WithFields(logrus.Fields{
		"server_port": cfg.Server.Port,
		"db_host":     cfg.Database.Host,
		"kafka_topic": cfg.Kafka.Topic,
		"cache_size":  cfg.Cache.MaxSize,
	}).Info("Configuration loaded")

	// Подключаемся к базе данных
	db, err := database.NewPostgresDB(&cfg.Database, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Создаем кеш
	orderCache := cache.NewMemoryCache(cfg.Cache.MaxSize, logger)

	// Восстанавливаем кеш из базы данных
	if err := restoreCache(db, orderCache, logger); err != nil {
		logger.WithError(err).Error("Failed to restore cache from database")
		// Не прерываем запуск, кеш будет заполняться по мере поступления запросов
	}

	// Создаем HTTP handler
	httpHandler := handlers.NewHTTPHandler(db, orderCache, logger)
	router := httpHandler.SetupRoutes()

	// Создаем и запускаем Kafka consumer
	consumer := kafka.NewConsumer(&cfg.Kafka, db, orderCache, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		logger.WithError(err).Fatal("Failed to start Kafka consumer")
	}
	defer consumer.Stop()

	// Настраиваем HTTP сервер
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		logger.WithField("address", server.Addr).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("HTTP server failed")
		}
	}()

	// Ожидание сигнала для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutdown signal received, starting graceful shutdown")

	// Останавливаем Kafka consumer
	if err := consumer.Stop(); err != nil {
		logger.WithError(err).Error("Failed to stop Kafka consumer")
	}

	// Останавливаем HTTP сервер с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Error("Failed to shutdown HTTP server gracefully")
	}

	logger.Info("Order Service stopped")
}

// restoreCache восстанавливает кеш из базы данных при запуске
func restoreCache(db database.OrderRepository, cache cache.OrderCache, logger *logrus.Logger) error {
	logger.Info("Restoring cache from database")

	// Получаем последние заказы из базы данных
	orders, err := db.GetAllOrders(1000) // загружаем до 1000 последних заказов
	if err != nil {
		return fmt.Errorf("failed to get orders from database: %w", err)
	}

	if len(orders) == 0 {
		logger.Info("No orders found in database")
		return nil
	}

	// Загружаем заказы в кеш
	cache.LoadFromDB(orders)

	logger.WithField("loaded_orders", len(orders)).Info("Cache restored successfully")
	return nil
}
