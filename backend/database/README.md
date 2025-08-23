# Database & Infrastructure - Order Service

Папка содержит все файлы, связанные с инфраструктурой: PostgreSQL, Kafka и дополнительные сервисы.

## Содержимое

- **docker-compose.yml** - конфигурация всей инфраструктуры (PostgreSQL, Kafka, Zookeeper, Kafka UI)
- **init.sql** - скрипт инициализации БД и пользователя
- **migrations/** - SQL миграции для создания таблиц
- **kafka-producer.go** - тестовый Kafka producer для отправки сообщений
- **DATABASE_SCHEMA.md** - подробная документация схемы БД
- **go.mod** / **go.sum** - зависимости для Kafka producer

## Быстрый запуск

```bash
# Запуск всей инфраструктуры
docker-compose up -d

# Проверка статуса всех сервисов
docker-compose ps

# Проверка подключения к PostgreSQL
docker exec -it order_service_postgres psql -U user -d order_service_db -c "\dt"

# Тестирование Kafka
go run kafka-producer.go
```

## Сервисы и порты

| Сервис | Порт | Описание |
|--------|------|----------|
| PostgreSQL | 5433 | База данных (внутренний: 5432) |
| Kafka | 9092 | Брокер сообщений |
| Zookeeper | 2181 | Координация Kafka |
| Kafka UI | 8080 | Веб-интерфейс для мониторинга Kafka |

## Структура данных

База данных содержит 4 основные таблицы:
- `orders` - основная информация о заказах
- `deliveries` - данные о доставке
- `payments` - платежная информация  
- `order_items` - товары в заказе

## Примеры запросов

### Получение заказа с полной информацией:
```sql
SELECT 
    o.order_uid,
    o.track_number,
    d.name AS customer_name,
    d.city,
    p.amount,
    COUNT(i.id) AS items_count
FROM orders o
LEFT JOIN deliveries d ON o.order_uid = d.order_uid
LEFT JOIN payments p ON o.order_uid = p.order_uid
LEFT JOIN order_items i ON o.order_uid = i.order_uid
WHERE o.order_uid = 'b563feb7b2b84b6test'
GROUP BY o.order_uid, d.id, p.id;
```

### Статистика по заказам:
```sql
SELECT 
    COUNT(*) as total_orders,
    AVG(p.amount) as avg_amount,
    COUNT(DISTINCT o.customer_id) as unique_customers
FROM orders o
LEFT JOIN payments p ON o.order_uid = p.order_uid;
```

## Kafka интеграция

### Тестирование отправки сообщений:
```bash
# Отправка тестового заказа в Kafka
go run kafka-producer.go

# Мониторинг топиков через Kafka UI
open http://localhost:8080
```

### Структура сообщений Kafka:
```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": { ... },
  "payment": { ... },
  "items": [ ... ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```

## Управление инфраструктурой

```bash
# Остановка всех сервисов
docker-compose down

# Остановка с удалением данных
docker-compose down -v

# Просмотр логов
docker-compose logs postgres
docker-compose logs kafka

# Перезапуск отдельного сервиса
docker-compose restart postgres
```
