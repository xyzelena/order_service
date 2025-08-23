# Backend - Order Service

Backend микросервиса для обработки заказов на Go с поддержкой PostgreSQL, Kafka и кеширования.

## Технологический стек

- **Язык**: Go 1.21+
- **База данных**: PostgreSQL 15+ (с миграциями)
- **Очередь сообщений**: Kafka 7.4.0 (с Zookeeper)
- **Кеширование**: LRU кеш в памяти (размер конфигурируется)
- **API**: REST API с CORS поддержкой
- **Архитектура**: Clean Architecture с внутренними модулями

## Настройка и запуск

### 1. Предварительные требования

**Убедитесь, что Docker запущен:**
```bash
# Проверка статуса Docker
docker --version
```

### 2. Запуск инфраструктуры (PostgreSQL + Kafka)

```bash
# Переход в папку database
cd database

# Запуск всей инфраструктуры
docker-compose up -d

# Проверка статуса
docker-compose ps
```

### 3. Настройки подключения

**PostgreSQL:**
- **Host**: localhost
- **Port**: 5433 (внешний), 5432 (внутренний)
- **Database**: order_service_db
- **Username**: user
- **Password**: 0000

**Kafka:**
- **Bootstrap Servers**: localhost:9092
- **Topic**: orders
- **Consumer Group**: order-service-group

**Kafka UI:**
- **URL**: http://localhost:8080

### 4. Дополнительные команды

```bash
# Остановка контейнера (из папки database)
docker-compose down

# Подключение к БД через psql
docker exec -it order_service_postgres psql -U user -d order_service_db

# Просмотр логов
docker-compose logs postgres
docker-compose logs kafka

# Тестирование Kafka
cd database && go run kafka-producer.go
```

## Структура файлов

```
backend/
├── app/                       # Go приложение Order Service
│   ├── cmd/server/           # Точка входа
│   ├── internal/             # Внутренние пакеты
│   │   ├── models/          # Модели данных
│   │   ├── database/        # PostgreSQL
│   │   ├── cache/           # Кеш в памяти
│   │   ├── kafka/           # Kafka consumer
│   │   └── handlers/        # HTTP API
│   ├── pkg/config/          # Конфигурация
│   ├── Dockerfile           # Docker образ
│   ├── Makefile            # Команды сборки
│   └── README.md           # Документация приложения
├── database/               # Все файлы базы данных
│   ├── docker-compose.yml  # Конфигурация PostgreSQL
│   ├── init.sql            # SQL скрипт инициализации БД
│   ├── migrations/         # SQL миграции
│   └── DATABASE_SCHEMA.md  # Документация схемы БД
└── README.md              # Документация backend
```

## Схема базы данных

Создается 4 основные таблицы:
- **orders** - основная информация о заказах
- **deliveries** - данные о доставке  
- **payments** - платежная информация
- **order_items** - товары в заказе

Подробное описание схемы см. в [database/DATABASE_SCHEMA.md](database/DATABASE_SCHEMA.md)

## Go приложение

### Запуск Order Service:

```bash
# Переход в папку приложения
cd app

# Установка зависимостей
make deps

# Сборка и запуск
make build
make run

# Или через Docker
make docker-build
make docker-run
```

### Доступные эндпоинты:
- **HTTP API**: http://localhost:8081/api/v1/
- **Health check**: http://localhost:8081/api/v1/health
- **Frontend server**: http://localhost:3000 (через make run-frontend)

### Основные API эндпоинты:
- `GET /api/v1/orders/{order_uid}` - получение заказа по ID
- `GET /api/v1/orders?limit=N` - список заказов
- `POST /api/v1/orders/random` - создание случайного заказа
- `GET /api/v1/cache/stats` - статистика кеша
- `GET /api/v1/health` - проверка здоровья сервиса

Подробная документация в [app/README.md](app/README.md)

## Примечания

- База данных автоматически инициализируется при первом запуске с помощью скрипта `init.sql`
- Пользователь `user` создается с полными правами на базу `order_service_db`
- Таблицы создаются автоматически через миграции
- Вставляется тестовый заказ для демонстрации работы
- Go приложение автоматически восстанавливает кеш из БД при запуске
- Убедитесь, что порты 5433, 9092, 8080, 8081 свободны перед запуском
- Frontend требует HTTP сервер для поддержки ES6 модулей
- Kafka может потребовать время для инициализации после запуска
