# Инструкция по запуску и тестированию всех компонентов системы Order Service

## Быстрый запуск полной системы

### 1. Запуск инфраструктуры (PostgreSQL + Kafka)

```bash
cd order_service/backend/database
docker-compose up -d
```

Это запустит:
- **PostgreSQL** на порту 5433 (внешний)
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
Frontend использует ES6 модули и современную JavaScript архитектуру. Браузеры требуют HTTP сервер для:
- Загрузки ES6 модулей (import/export)
- Предотвращения CORS ошибок при API запросах
- Корректной работы модульной архитектуры

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
- Создание случайных заказов через POST /api/v1/orders/random
- Статистика кеша через GET /api/v1/cache/stats
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

### Тест 4: Модульная архитектура Frontend

```bash
# Автоматическая проверка модулей
cd order_service
make check-frontend
```

**Что проверяется:**
- Наличие всех ES6 модулей
- Корректность экспортов (export)
- Корректность импортов (import)
- Структура модульной архитектуры

### Тест 5: Frontend интерфейс

1. Откройте http://localhost:3000
2. Введите `b563feb7b2b84b6test` в поле поиска
3. Нажмите "Найти заказ"
4. Проверьте отображение данных
5. Протестируйте дополнительные функции:
   - "Статистика кеша" - модальное окно с данными кеша
   - "Список заказов" - интерактивный список для выбора
   - "Создать случайный заказ" - генерация и отображение нового заказа
   - "Проверка сервиса" - статус всех компонентов системы

**Что проверяется:**
- Загрузка ES6 модулей (api.js, modal.js, notifications.js, etc.)
- Модульная архитектура JavaScript
- Подключение к API на порту 8081
- Отображение полной информации о заказе
- Работа модальных окон и уведомлений
- Интерактивность и UX элементов

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

# Тестирование новых функций
curl -X POST http://localhost:8081/api/v1/orders/random | jq
curl http://localhost:8081/api/v1/cache/stats | jq

# Терминал 4: Тест производительности
cd order_service/backend/app/tests
go run cache_benchmark.go

# Терминал 5: Проверка модульной архитектуры
cd order_service
make check-frontend
```

## Мониторинг и отладка

### Полезные команды:

```bash
# Быстрая проверка всей системы
cd order_service
make quick-test

# Проверка модульной структуры frontend
make check-frontend

# Логи всех контейнеров
docker-compose logs -f

# Статистика кеша
curl http://localhost:8081/api/v1/cache/stats | jq

# Проверка здоровья
curl http://localhost:8081/api/v1/health | jq

# Создание тестового заказа
curl -X POST http://localhost:8081/api/v1/orders/random | jq

# Список заказов
curl "http://localhost:8081/api/v1/orders?limit=5" | jq

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
ps aux | grep order-service

# Проверить порт 8081
lsof -i :8081

# Проверить логи приложения
cd order_service/backend/app
make run

# Быстрая проверка API
curl http://localhost:8081/api/v1/health
```

### Проблема: Frontend модули не загружаются

```bash
# Проверить структуру модулей
cd order_service
make check-frontend

# Проверить HTTP сервер
lsof -i :3000

# Проверить консоль браузера на ошибки CORS/модулей
# Убедитесь, что используется HTTP сервер, а не file://
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
6. **Frontend ES6 модули** загружаются корректно
7. **Модульная архитектура** работает без ошибок
8. **CORS запросы** выполняются успешно
9. **Генерация заказов** работает через POST API
10. **Статистика кеша** доступна через API
11. **Модальные окна** и **уведомления** функционируют
12. **make check-frontend** проходит без ошибок

### Команды для финальной проверки:

```bash
# Полная проверка системы
cd order_service
make test-full

# Быстрая проверка основных функций  
make quick-test

# Проверка модульной архитектуры
make check-frontend
```

Все тесты должны пройти успешно для подтверждения готовности системы.
