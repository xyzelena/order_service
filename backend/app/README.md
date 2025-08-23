# Order Service - Go Application

Микросервис для обработки заказов на Go с использованием PostgreSQL, Kafka и кеширования в памяти.

## 🚀 Функциональность

1. **HTTP API** - RESTful API для получения и создания заказов
2. **Kafka Consumer** - получение заказов из очереди сообщений в реальном времени
3. **PostgreSQL** - надежное хранение данных в реляционной БД с миграциями
4. **LRU Cache** - быстрый доступ к часто запрашиваемым заказам
5. **Cache Recovery** - автоматическое восстановление кеша из БД при запуске
6. **Frontend Server** - HTTP сервер для модульного ES6 frontend
7. **Random Order Generation** - создание тестовых заказов через API
8. **CORS Support** - поддержка кросс-доменных запросов
9. **Error Handling** - обработка ошибок и валидация данных
10. **Graceful Shutdown** - корректное завершение работы

## 📁 Структура проекта

```
app/
├── cmd/server/              # Точка входа в приложение
│   └── main.go             # Основной файл запуска
├── internal/               # Внутренние пакеты
│   ├── models/             # Модели данных
│   ├── database/           # Работа с PostgreSQL
│   ├── cache/              # Кеширование в памяти
│   ├── kafka/              # Kafka consumer
│   └── handlers/           # HTTP handlers и API
├── pkg/config/             # Конфигурация приложения
├── static/                 # Статические файлы для веб-интерфейса
├── Dockerfile              # Docker образ
├── Makefile               # Команды для сборки и запуска
├── go.mod                 # Go модуль и зависимости
└── env.example            # Пример переменных окружения
```

## 🛠️ Технологический стек

- **Язык**: Go 1.21+
- **База данных**: PostgreSQL 15
- **Очередь сообщений**: Apache Kafka
- **HTTP Router**: Gorilla Mux
- **Логирование**: Logrus
- **Containerization**: Docker

## ⚙️ Настройка и запуск

### Предварительные требования

1. Go 1.21+
2. PostgreSQL (через Docker или локально)
3. Apache Kafka (опционально, для полной функциональности)

### Быстрый старт

1. **Настройка переменных окружения**:
   ```bash
   cp env.example .env
   # Отредактируйте .env под ваши настройки
   ```

2. **Установка зависимостей**:
   ```bash
   make deps
   ```

3. **Сборка и запуск**:
   ```bash
   make build
   make run
   ```

4. **Проверка работы**:
   - API доступно по адресу http://localhost:8081/api/v1/
   - Frontend сервер: http://localhost:3000 (через make run-frontend)

### Запуск через Docker

```bash
make docker-build
make docker-run
```

## 📡 API Endpoints

### Основные эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/api/v1/orders/{order_uid}` | Получить заказ по UID |
| `GET` | `/api/v1/orders?limit=N` | Получить список заказов |
| `POST` | `/api/v1/orders/random` | Создать случайный заказ |
| `GET` | `/api/v1/cache/stats` | Статистика кеша |
| `GET` | `/api/v1/health` | Проверка здоровья сервиса |

### Frontend HTTP Server

| Команда | Описание |
|---------|----------|
| `make run-frontend` | Запуск HTTP сервера для ES6 модулей |
| `http://localhost:3000` | Модульный веб-интерфейс |

### Примеры запросов

**Получение заказа**:
```bash
curl http://localhost:8081/api/v1/orders/b563feb7b2b84b6test
```

**Создание случайного заказа**:
```bash
curl -X POST http://localhost:8081/api/v1/orders/random
```

**Статистика кеша**:
```bash
curl http://localhost:8081/api/v1/cache/stats
```

**Проверка здоровья**:
```bash
curl http://localhost:8081/api/v1/health
```

## 🗄️ Схема данных

Приложение работает с 4 основными таблицами:

- **orders** - основная информация о заказах
- **deliveries** - данные о доставке
- **payments** - платежная информация
- **order_items** - товары в заказе

Подробная схема БД описана в [../database/DATABASE_SCHEMA.md](../database/DATABASE_SCHEMA.md)

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `SERVER_HOST` | Хост HTTP сервера | `0.0.0.0` |
| `SERVER_PORT` | Порт HTTP сервера | `8081` |
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `user` |
| `DB_PASSWORD` | Пароль БД | `0000` |
| `DB_NAME` | Имя базы данных | `order_service_db` |
| `KAFKA_BROKERS` | Адреса Kafka брокеров | `localhost:9092` |
| `KAFKA_TOPIC` | Топик для заказов | `orders` |
| `KAFKA_GROUP_ID` | ID группы consumer'а | `order-service-group` |
| `CACHE_MAX_SIZE` | Максимальный размер кеша | `1000` |
| `DEBUG` | Режим отладки | `false` |

## 🎯 Архитектурные решения

### 1. Кеширование
- **LRU алгоритм** - вытеснение старых записей при переполнении
- **Thread-safe** - безопасная работа в многопоточной среде
- **Recovery** - восстановление при перезапуске из БД

### 2. Обработка ошибок
- **Валидация данных** - проверка обязательных полей
- **Graceful degradation** - продолжение работы при ошибках
- **Логирование** - детальное логирование всех операций

### 3. Kafka Integration
- **Идемпотентность** - повторная обработка сообщений безопасна
- **Error handling** - невалидные сообщения логируются и пропускаются
- **Graceful shutdown** - корректное завершение consumer'а

### 4. Database
- **Транзакции** - атомарное сохранение связанных данных
- **Connection pooling** - эффективное использование соединений
- **Normalized schema** - оптимизированная структура БД

## 🧪 Тестирование

```bash
# Запуск тестов
make test

# Проверка кода линтером
make lint

# Форматирование кода
make fmt
```

## 📊 Мониторинг и метрики

Приложение предоставляет следующие метрики:

- **Health check** - `/api/v1/health`
- **Cache statistics** - размер кеша, hit/miss ratio
- **HTTP request logs** - время ответа, статус коды
- **Database connection status**

## 🚀 Production Ready Features

1. **Graceful Shutdown** - корректное завершение всех соединений
2. **Configuration via Environment** - 12-factor app принципы
3. **Structured Logging** - JSON формат для централизованного логирования
4. **Health Checks** - эндпоинт для мониторинга
5. **Docker Support** - готовый образ для деплоя
6. **Error Recovery** - устойчивость к временным сбоям

## 🔍 Отладка

1. **Включить debug логи**:
   ```bash
   export DEBUG=true
   ./server
   ```

2. **Проверить подключение к БД**:
   ```bash
   curl http://localhost:8081/api/v1/health
   ```

3. **Просмотр статистики кеша**:
   ```bash
   curl http://localhost:8081/api/v1/cache/stats
   ```

4. **Тестирование frontend**:
   ```bash
   make run-frontend  # Запуск на порту 3000
   ```

## 📝 Примечания

- При отсутствии Kafka, приложение запустится только с HTTP API
- Кеш автоматически восстанавливается из БД при запуске
- Все API возвращают JSON с единообразной структурой ответов
- Frontend работает на модульной ES6 архитектуре
- CORS настроен для работы с frontend на порту 3000
- Поддерживается создание случайных заказов для тестирования
