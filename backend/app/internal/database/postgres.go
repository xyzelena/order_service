package database

import (
	"database/sql"
	"fmt"
	"order-service/internal/models"
	"order-service/pkg/config"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	db     *sql.DB
	logger *logrus.Logger
}

type OrderRepository interface {
	CreateOrder(orderFull *models.OrderFull) error
	GetOrderByUID(orderUID string) (*models.OrderFull, error)
	GetAllOrders(limit int) ([]models.OrderFull, error)
	OrderExists(orderUID string) (bool, error)
}

func NewPostgresDB(cfg *config.DatabaseConfig, logger *logrus.Logger) (*PostgresDB, error) {
	dsn := cfg.GetDSN()
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL")

	return &PostgresDB{
		db:     db,
		logger: logger,
	}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// CreateOrder сохраняет полный заказ в базу данных с использованием транзакции
func (p *PostgresDB) CreateOrder(orderFull *models.OrderFull) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Вставляем основной заказ
	orderQuery := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
						   customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(orderQuery,
		orderFull.OrderUID, orderFull.TrackNumber, orderFull.Entry, orderFull.Locale,
		orderFull.InternalSignature, orderFull.CustomerID, orderFull.DeliveryService,
		orderFull.Shardkey, orderFull.SmID, orderFull.DateCreated, orderFull.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// 2. Вставляем данные доставки
	if orderFull.Delivery != nil {
		deliveryQuery := `
			INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = tx.Exec(deliveryQuery,
			orderFull.OrderUID, orderFull.Delivery.Name, orderFull.Delivery.Phone,
			orderFull.Delivery.Zip, orderFull.Delivery.City, orderFull.Delivery.Address,
			orderFull.Delivery.Region, orderFull.Delivery.Email,
		)
		if err != nil {
			return fmt.Errorf("failed to insert delivery: %w", err)
		}
	}

	// 3. Вставляем платежные данные
	if orderFull.Payment != nil {
		paymentQuery := `
			INSERT INTO payments (order_uid, transaction, request_id, currency, provider, 
								 amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		_, err = tx.Exec(paymentQuery,
			orderFull.OrderUID, orderFull.Payment.Transaction, orderFull.Payment.RequestID,
			orderFull.Payment.Currency, orderFull.Payment.Provider, orderFull.Payment.Amount,
			orderFull.Payment.PaymentDt, orderFull.Payment.Bank, orderFull.Payment.DeliveryCost,
			orderFull.Payment.GoodsTotal, orderFull.Payment.CustomFee,
		)
		if err != nil {
			return fmt.Errorf("failed to insert payment: %w", err)
		}
	}

	// 4. Вставляем товары
	for _, item := range orderFull.Items {
		itemQuery := `
			INSERT INTO order_items (order_uid, chrt_id, track_number, price, rid, name, 
								   sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`
		_, err = tx.Exec(itemQuery,
			orderFull.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	p.logger.WithField("order_uid", orderFull.OrderUID).Info("Order saved successfully")
	return nil
}

// GetOrderByUID получает полную информацию о заказе по UID
func (p *PostgresDB) GetOrderByUID(orderUID string) (*models.OrderFull, error) {
	orderFull := &models.OrderFull{}

	// 1. Получаем основной заказ
	orderQuery := `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
			   delivery_service, shardkey, sm_id, date_created, oof_shard, created_at, updated_at
		FROM orders WHERE order_uid = $1
	`
	err := p.db.QueryRow(orderQuery, orderUID).Scan(
		&orderFull.OrderUID, &orderFull.TrackNumber, &orderFull.Entry, &orderFull.Locale,
		&orderFull.InternalSignature, &orderFull.CustomerID, &orderFull.DeliveryService,
		&orderFull.Shardkey, &orderFull.SmID, &orderFull.DateCreated, &orderFull.OofShard,
		&orderFull.CreatedAt, &orderFull.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// 2. Получаем данные доставки
	deliveryQuery := `
		SELECT id, order_uid, name, phone, zip, city, address, region, email, created_at
		FROM deliveries WHERE order_uid = $1
	`
	delivery := &models.Delivery{}
	err = p.db.QueryRow(deliveryQuery, orderUID).Scan(
		&delivery.ID, &delivery.OrderUID, &delivery.Name, &delivery.Phone, &delivery.Zip,
		&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email, &delivery.CreatedAt,
	)
	if err == nil {
		orderFull.Delivery = delivery
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	// 3. Получаем платежные данные
	paymentQuery := `
		SELECT id, order_uid, transaction, request_id, currency, provider, amount, payment_dt,
			   bank, delivery_cost, goods_total, custom_fee, created_at
		FROM payments WHERE order_uid = $1
	`
	payment := &models.Payment{}
	err = p.db.QueryRow(paymentQuery, orderUID).Scan(
		&payment.ID, &payment.OrderUID, &payment.Transaction, &payment.RequestID,
		&payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		&payment.CreatedAt,
	)
	if err == nil {
		orderFull.Payment = payment
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// 4. Получаем товары
	itemsQuery := `
		SELECT id, order_uid, chrt_id, track_number, price, rid, name, sale, size,
			   total_price, nm_id, brand, status, created_at
		FROM order_items WHERE order_uid = $1 ORDER BY id
	`
	rows, err := p.db.Query(itemsQuery, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderUID, &item.ChrtID, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status, &item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}
	orderFull.Items = items

	return orderFull, nil
}

// GetAllOrders получает все заказы с ограничением
func (p *PostgresDB) GetAllOrders(limit int) ([]models.OrderFull, error) {
	query := `
		SELECT order_uid FROM orders 
		ORDER BY created_at DESC LIMIT $1
	`
	rows, err := p.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get order list: %w", err)
	}
	defer rows.Close()

	var orders []models.OrderFull
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, fmt.Errorf("failed to scan order UID: %w", err)
		}

		orderFull, err := p.GetOrderByUID(orderUID)
		if err != nil {
			p.logger.WithError(err).WithField("order_uid", orderUID).Error("Failed to get full order")
			continue
		}
		if orderFull != nil {
			orders = append(orders, *orderFull)
		}
	}

	return orders, nil
}

// OrderExists проверяет существование заказа
func (p *PostgresDB) OrderExists(orderUID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)`
	err := p.db.QueryRow(query, orderUID).Scan(&exists)
	return exists, err
}
