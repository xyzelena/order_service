package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/models"
	"order-service/pkg/config"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	reader   *kafka.Reader
	db       database.OrderRepository
	cache    cache.OrderCache
	logger   *logrus.Logger
	stopChan chan struct{}
}

type MessageProcessor interface {
	Start(ctx context.Context) error
	Stop() error
}

// NewConsumer создает новый Kafka consumer
func NewConsumer(cfg *config.KafkaConfig, db database.OrderRepository, cache cache.OrderCache, logger *logrus.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		MinBytes:    1,        // 1 байт минимум
		MaxBytes:    10e6,     // 10MB максимум
		MaxWait:     100 * time.Millisecond,
		StartOffset: kafka.LastOffset, // Начинаем с последнего сообщения
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			logger.Errorf("Kafka error: "+msg, args...)
		}),
	})

	return &Consumer{
		reader:   reader,
		db:       db,
		cache:    cache,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

// Start запускает consumer в отдельной горутине
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("Starting Kafka consumer")

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Info("Context cancelled, stopping consumer")
				return
			case <-c.stopChan:
				c.logger.Info("Stop signal received, stopping consumer")
				return
			default:
				if err := c.processMessage(ctx); err != nil {
					c.logger.WithError(err).Error("Failed to process Kafka message")
					// Небольшая задержка перед следующей попыткой
					time.Sleep(time.Second)
				}
			}
		}
	}()

	return nil
}

// Stop останавливает consumer
func (c *Consumer) Stop() error {
	c.logger.Info("Stopping Kafka consumer")
	close(c.stopChan)
	return c.reader.Close()
}

// processMessage обрабатывает одно сообщение из Kafka
func (c *Consumer) processMessage(ctx context.Context) error {
	// Устанавливаем таймаут для чтения сообщения
	readCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	msg, err := c.reader.ReadMessage(readCtx)
	if err != nil {
		if err == context.DeadlineExceeded {
			// Таймаут - это нормально, просто нет сообщений
			return nil
		}
		return fmt.Errorf("failed to read message: %w", err)
	}

	c.logger.WithFields(logrus.Fields{
		"partition": msg.Partition,
		"offset":    msg.Offset,
		"topic":     msg.Topic,
	}).Debug("Received Kafka message")

	// Парсим и валидируем сообщение
	orderFull, err := c.parseAndValidateMessage(msg.Value)
	if err != nil {
		c.logger.WithError(err).WithField("raw_message", string(msg.Value)).Error("Failed to parse message")
		// Возвращаем nil, чтобы не останавливать consumer из-за одного невалидного сообщения
		return nil
	}

	// Проверяем, что заказ еще не существует
	exists, err := c.db.OrderExists(orderFull.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to check if order exists: %w", err)
	}

	if exists {
		c.logger.WithField("order_uid", orderFull.OrderUID).Info("Order already exists, skipping")
		return nil
	}

	// Сохраняем в базу данных
	if err := c.db.CreateOrder(orderFull); err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	// Добавляем в кеш
	c.cache.Set(orderFull.OrderUID, orderFull)

	c.logger.WithField("order_uid", orderFull.OrderUID).Info("Order processed successfully")
	return nil
}

// parseAndValidateMessage парсит JSON сообщение и конвертирует его в модель
func (c *Consumer) parseAndValidateMessage(data []byte) (*models.OrderFull, error) {
	var kafkaMsg models.KafkaOrderMessage
	if err := json.Unmarshal(data, &kafkaMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Валидация обязательных полей
	if err := c.validateKafkaMessage(&kafkaMsg); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// Конвертируем в наши модели
	orderFull, err := c.convertKafkaToModel(&kafkaMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to convert message: %w", err)
	}

	return orderFull, nil
}

// validateKafkaMessage проверяет валидность данных
func (c *Consumer) validateKafkaMessage(msg *models.KafkaOrderMessage) error {
	if strings.TrimSpace(msg.OrderUID) == "" {
		return fmt.Errorf("order_uid is required")
	}
	if strings.TrimSpace(msg.TrackNumber) == "" {
		return fmt.Errorf("track_number is required")
	}
	if strings.TrimSpace(msg.CustomerID) == "" {
		return fmt.Errorf("customer_id is required")
	}
	if len(msg.Items) == 0 {
		return fmt.Errorf("items are required")
	}

	// Валидация доставки
	if strings.TrimSpace(msg.Delivery.Name) == "" {
		return fmt.Errorf("delivery name is required")
	}
	if strings.TrimSpace(msg.Delivery.Phone) == "" {
		return fmt.Errorf("delivery phone is required")
	}

	// Валидация платежа
	if msg.Payment.Amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}
	if strings.TrimSpace(msg.Payment.Currency) == "" {
		return fmt.Errorf("payment currency is required")
	}

	return nil
}

// convertKafkaToModel конвертирует Kafka сообщение в наши модели
func (c *Consumer) convertKafkaToModel(msg *models.KafkaOrderMessage) (*models.OrderFull, error) {
	// Парсим дату
	dateCreated, err := time.Parse(time.RFC3339, msg.DateCreated)
	if err != nil {
		c.logger.WithError(err).WithField("date", msg.DateCreated).Warn("Failed to parse date, using current time")
		dateCreated = time.Now()
	}

	// Основной заказ
	order := models.Order{
		OrderUID:          msg.OrderUID,
		TrackNumber:       msg.TrackNumber,
		Entry:             msg.Entry,
		Locale:            msg.Locale,
		InternalSignature: msg.InternalSignature,
		CustomerID:        msg.CustomerID,
		DeliveryService:   msg.DeliveryService,
		Shardkey:          msg.Shardkey,
		SmID:              msg.SmID,
		DateCreated:       dateCreated,
		OofShard:          msg.OofShard,
	}

	// Доставка
	delivery := &models.Delivery{
		OrderUID: msg.OrderUID,
		Name:     msg.Delivery.Name,
		Phone:    msg.Delivery.Phone,
		Zip:      msg.Delivery.Zip,
		City:     msg.Delivery.City,
		Address:  msg.Delivery.Address,
		Region:   msg.Delivery.Region,
		Email:    msg.Delivery.Email,
	}

	// Платеж
	payment := &models.Payment{
		OrderUID:     msg.OrderUID,
		Transaction:  msg.Payment.Transaction,
		RequestID:    msg.Payment.RequestID,
		Currency:     msg.Payment.Currency,
		Provider:     msg.Payment.Provider,
		Amount:       msg.Payment.Amount,
		PaymentDt:    msg.Payment.PaymentDt,
		Bank:         msg.Payment.Bank,
		DeliveryCost: msg.Payment.DeliveryCost,
		GoodsTotal:   msg.Payment.GoodsTotal,
		CustomFee:    msg.Payment.CustomFee,
	}

	// Товары
	var items []models.OrderItem
	for _, kafkaItem := range msg.Items {
		item := models.OrderItem{
			OrderUID:    msg.OrderUID,
			ChrtID:      kafkaItem.ChrtID,
			TrackNumber: kafkaItem.TrackNumber,
			Price:       kafkaItem.Price,
			Rid:         kafkaItem.Rid,
			Name:        kafkaItem.Name,
			Sale:        kafkaItem.Sale,
			Size:        kafkaItem.Size,
			TotalPrice:  kafkaItem.TotalPrice,
			NmID:        kafkaItem.NmID,
			Brand:       kafkaItem.Brand,
			Status:      kafkaItem.Status,
		}
		items = append(items, item)
	}

	return &models.OrderFull{
		Order:    order,
		Delivery: delivery,
		Payment:  payment,
		Items:    items,
	}, nil
}
