# Инструкция по запуску и тестированию всех компонентов системы Order Service

## Быстрый запуск полной системы

### 1. Запуск инфраструктуры (PostgreSQL + Kafka)

```bash
cd order_service/backend/database
docker-compose up -d
```

Это запустит:
- **PostgreSQL** на порту 5432
- **Kafka** на порту 9092
- **Zookeeper** на порту 2181
- **Kafka UI** на порту 8080

### 2. Проверка готовности инфраструктуры

```bash
# Проверка PostgreSQL
docker exec -it order_service_postgres psql -U user -d order_service_db -c "\dt"

# Проверка Kafka
docker logs order_service_kafka | grep "started"

# Kafka UI (опционально)
# Откройте http://localhost:8080 в браузере
```

### 3. Запуск Go микросервиса

```bash
cd order_service/backend/app
make deps
make build
make run

# Или через Docker
make docker-build
make docker-run
```

### 4. Запуск Frontend (опционально)

**Варианты запуска:**

```bash
# Вариант 1: Go сервер (рекомендуется)
cd order_service/backend/app
make run-frontend

# Вариант 2: Python сервер
cd order_service/frontend
python -m http.server 3000

# Вариант 3: Node.js сервер
npx serve frontend -p 3000

# Откройте http://localhost:3000
```

**Почему нужен сервер?**
Frontend написан на чистом HTML/CSS/JavaScript, но браузеры блокируют CORS запросы к API при открытии файлов напрямую (file:// протокол). HTTP сервер решает эту проблему.

## Тестирование компонентов

### Тест 1: HTTP API и JSON формат

```bash
cd order_service/backend/app/tests
./api_test.sh
```

**Что проверяется:**
- Доступность API на порту 8081
- Корректность JSON ответов
- Структура API responses `{"success": true/false, "data": ..., "error": ...}`
- Обработка ошибок (404 для несуществующих заказов)
- Производительность API

**Ожидаемый результат:**
```json
{
  "success": true,
  "data": {
    "order_uid": "b563feb7b2b84b6test",
    "track_number": "WBILMTESTTRACK",
    "customer_id": "test",
    "delivery": { ... },
    "payment": { ... },
    "items": [ ... ]
  }
}
```

### Тест 2: Производительность кеша

```bash
cd order_service/backend/app/tests
go run cache_benchmark.go
```

**Что проверяется:**
- Первый запрос медленнее (загрузка из БД)
- Повторные запросы быстрее (из кеша)
- Ускорение кеша (измерение в разах)
- Параллельные запросы
- Статистика использования кеша

**Ожидаемые результаты:**
```
Первый запрос (БД):           50ms
Среднее время кеша:          5ms
Ускорение кеша:              10.00x
```

### Тест 3: Kafka интеграция

```bash
cd order_service/backend/database
go mod tidy
go run kafka-producer.go
```

**Что проверяется:**
- Подключение к Kafka
- Отправка сообщений в топик "orders"
- Обработка сообщений Go микросервисом
- Сохранение в БД через Kafka
- Автоматическое добавление в кеш

**Мониторинг:**
- Логи Go приложения: новые заказы из Kafka
- Kafka UI: http://localhost:8080
- БД: новые записи в таблице orders

### Тест 4: Frontend интерфейс

1. Откройте http://localhost:3000
2. Введите `b563feb7b2b84b6test` в поле поиска
3. Нажмите "Найти заказ"
4. Проверьте отображение данных

**Что проверяется:**
- Подключение к API на порту 8081
- Отображение полной информации о заказе
- Работа дополнительных функций (статистика, список заказов)

## Комплексный тест системы

### Сценарий полного цикла:

1. **Запустите Kafka producer** (эмулирует поступление заказов)
2. **Отправьте заказ через Kafka** 
3. **Проверьте появление в БД**
4. **Получите через API** 
5. **Убедитесь в работе кеша**

```bash
# Терминал 1: Запуск producer
cd order_service/backend/database
go run kafka-producer.go

# Терминал 2: Мониторинг логов Go приложения
cd order_service/backend/app
make run

# Терминал 3: Тестирование API
curl http://localhost:8081/api/v1/orders/b563feb7b2b84b6test | jq

# Терминал 4: Тест производительности
cd order_service/backend/app/tests
go run cache_benchmark.go
```

## Мониторинг и отладка

### Полезные команды:

```bash
# Логи всех контейнеров
docker-compose logs -f

# Статистика кеша
curl http://localhost:8081/api/v1/cache/stats | jq

# Проверка здоровья
curl http://localhost:8081/api/v1/health | jq

# Просмотр таблиц БД
docker exec -it order_service_postgres psql -U user -d order_service_db -c "
SELECT o.order_uid, o.customer_id, o.date_created 
FROM orders o 
ORDER BY o.created_at DESC 
LIMIT 5;"
```

### Kafka мониторинг:

```bash
# Просмотр топиков
docker exec -it order_service_kafka kafka-topics --bootstrap-server localhost:9092 --list

# Просмотр сообщений
docker exec -it order_service_kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic orders --from-beginning --max-messages 5
```

## Troubleshooting

### Проблема: API недоступен

```bash
# Проверить, что сервис запущен
ps aux | grep server

# Проверить порт
lsof -i :8081

# Проверить логи приложения
```

### Проблема: Kafka не работает

```bash
# Проверить статус контейнеров
docker-compose ps

# Перезапустить Kafka
docker-compose restart kafka

# Проверить логи
docker-compose logs kafka
```

### Проблема: БД недоступна

```bash
# Проверить подключение
docker exec -it order_service_postgres pg_isready

# Проверить таблицы
docker exec -it order_service_postgres psql -U user -d order_service_db -c "\dt"
```

## Ожидаемые результаты тестов

### HTTP API 
- Время ответа: < 100ms для кешированных данных
- Время ответа: < 500ms для загрузки из БД
- Формат JSON: всегда валидный с полями success/data/error

### Кеш
- Ускорение: 5-20x для повторных запросов
- Память: эффективное использование LRU
- Thread-safety: работа при параллельных запросах

### Kafka
- Latency: < 100ms от отправки до сохранения в БД
- Throughput: > 1000 сообщений/сек
- Надежность: гарантированная доставка сообщений

### База данных
- ACID: транзакционность операций
- Integrity: целостность связанных данных
- Performance: < 50ms для простых запросов

## Критерии успешного тестирования

1. **API отвечает** на `GET http://localhost:8081/api/v1/orders/<order_uid>`
2. **JSON корректный** с полями success, data, error
3. **Кеш ускоряет** повторные запросы минимум в 5 раз
4. **Kafka обрабатывает** сообщения онлайн
5. **БД сохраняет** данные транзакционно
6. **Frontend отображает** данные корректно

Все тесты должны пройти успешно для подтверждения готовности системы.
