package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaOrderMessage структура сообщения заказа для Kafka
type KafkaOrderMessage struct {
	OrderUID          string                `json:"order_uid"`
	TrackNumber       string                `json:"track_number"`
	Entry             string                `json:"entry"`
	Delivery          KafkaDelivery         `json:"delivery"`
	Payment           KafkaPayment          `json:"payment"`
	Items             []KafkaOrderItem      `json:"items"`
	Locale            string                `json:"locale"`
	InternalSignature string                `json:"internal_signature"`
	CustomerID        string                `json:"customer_id"`
	DeliveryService   string                `json:"delivery_service"`
	Shardkey          string                `json:"shardkey"`
	SmID              int                   `json:"sm_id"`
	DateCreated       string                `json:"date_created"`
	OofShard          string                `json:"oof_shard"`
}

type KafkaDelivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type KafkaPayment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type KafkaOrderItem struct {
	ChrtID      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	// Настройка Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "orders",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: 1, // kafka.RequireOne equivalent
		Async:        false,
	})
	defer writer.Close()

	// Данные для генерации случайных заказов
	names := []string{"Иван Иванов", "Мария Петрова", "Алексей Сидоров", "Елена Кузнецова", "Дмитрий Волков"}
	cities := []string{"Москва", "Санкт-Петербург", "Новосибирск", "Екатеринбург", "Казань"}
	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "New Balance"}
	products := []string{"Кроссовки", "Футболка", "Шорты", "Куртка", "Джинсы"}
	
	log.Println("Kafka Producer запущен")
	log.Println("Отправка тестовых сообщений каждые 5 секунд...")
	log.Println("Нажмите Ctrl+C для остановки")

	// Отправляем стандартный тестовый заказ
	log.Println("Отправка тестового заказа...")
	testOrder := createTestOrder()
	if err := sendOrder(writer, testOrder); err != nil {
		log.Printf("Ошибка отправки тестового заказа: %v", err)
	} else {
		log.Printf("Тестовый заказ отправлен: %s", testOrder.OrderUID)
	}

	// Бесконечный цикл отправки случайных заказов
	for {
		time.Sleep(5 * time.Second)
		
		order := generateRandomOrder(names, cities, brands, products)
		if err := sendOrder(writer, order); err != nil {
			log.Printf("Ошибка отправки заказа %s: %v", order.OrderUID, err)
		} else {
			log.Printf("Заказ отправлен: %s (клиент: %s, товаров: %d)", 
				order.OrderUID, order.CustomerID, len(order.Items))
		}
	}
}

func createTestOrder() KafkaOrderMessage {
	return KafkaOrderMessage{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: KafkaDelivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: KafkaPayment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []KafkaOrderItem{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          "1",
	}
}

func generateRandomOrder(names, cities, brands, products []string) KafkaOrderMessage {
	orderUID := fmt.Sprintf("order_%d_%d", time.Now().Unix(), rand.Intn(10000))
	trackNumber := fmt.Sprintf("WB%d", rand.Intn(1000000))
	customerID := fmt.Sprintf("customer_%d", rand.Intn(10000))
	
	// Случайная доставка
	delivery := KafkaDelivery{
		Name:    names[rand.Intn(len(names))],
		Phone:   fmt.Sprintf("+7%09d", rand.Intn(1000000000)),
		Zip:     fmt.Sprintf("%06d", rand.Intn(1000000)),
		City:    cities[rand.Intn(len(cities))],
		Address: fmt.Sprintf("ул. %s, д. %d", "Тестовая", rand.Intn(100)+1),
		Region:  "Тестовый регион",
		Email:   fmt.Sprintf("user%d@test.com", rand.Intn(1000)),
	}

	// Случайный платеж
	goodsTotal := rand.Intn(10000) + 1000
	deliveryCost := rand.Intn(500) + 200
	payment := KafkaPayment{
		Transaction:  orderUID,
		RequestID:    "",
		Currency:     "RUB",
		Provider:     "wbpay",
		Amount:       goodsTotal + deliveryCost,
		PaymentDt:    time.Now().Unix(),
		Bank:         "sberbank",
		DeliveryCost: deliveryCost,
		GoodsTotal:   goodsTotal,
		CustomFee:    0,
	}

	// Случайные товары (1-3 товара)
	itemCount := rand.Intn(3) + 1
	items := make([]KafkaOrderItem, itemCount)
	
	for i := 0; i < itemCount; i++ {
		price := rand.Intn(5000) + 500
		sale := rand.Intn(50)
		totalPrice := price * (100 - sale) / 100
		
		items[i] = KafkaOrderItem{
			ChrtID:      int64(rand.Intn(10000000) + 1000000),
			TrackNumber: trackNumber,
			Price:       price,
			Rid:         fmt.Sprintf("rid_%d_%d", time.Now().Unix(), i),
			Name:        products[rand.Intn(len(products))],
			Sale:        sale,
			Size:        fmt.Sprintf("%d", rand.Intn(10)+35),
			TotalPrice:  totalPrice,
			NmID:        int64(rand.Intn(10000000) + 1000000),
			Brand:       brands[rand.Intn(len(brands))],
			Status:      202,
		}
	}

	return KafkaOrderMessage{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            "ru",
		InternalSignature: "",
		CustomerID:        customerID,
		DeliveryService:   "cdek",
		Shardkey:          fmt.Sprintf("%d", rand.Intn(10)),
		SmID:              rand.Intn(100) + 1,
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          "1",
	}
}

func sendOrder(writer *kafka.Writer, order KafkaOrderMessage) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: data,
	}

	return writer.WriteMessages(context.Background(), message)
}
