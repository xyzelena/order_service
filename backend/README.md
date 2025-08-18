# Backend - Order Service

Backend микросервиса для обработки заказов на Go.

## Технологический стек

- **Язык**: Go
- **База данных**: PostgreSQL
- **Очередь сообщений**: Kafka
- **Кеширование**: In-memory cache

## Настройка и запуск

### 1. Предварительные требования

**Убедитесь, что Docker запущен:**
```bash
# Проверка статуса Docker
docker --version
```

### 2. Запуск PostgreSQL

```bash
# Переход в папку database
cd database

# Запуск PostgreSQL контейнера
docker-compose up -d postgres

# Проверка статуса
docker-compose ps
```

### 3. Настройки подключения к БД

- **Host**: localhost
- **Port**: 5432
- **Database**: order_service_db
- **Username**: user
- **Password**: 0000

### 4. Дополнительные команды

```bash
# Остановка контейнера (из папки database)
docker-compose down

# Подключение к БД через psql
docker exec -it order_service_postgres psql -U user -d order_service_db

# Просмотр логов
docker-compose logs postgres
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
- **HTTP API**: http://localhost:8080/api/v1/
- **Веб-интерфейс**: http://localhost:8080/
- **Health check**: http://localhost:8080/api/v1/health

Подробная документация в [app/README.md](app/README.md)

## Примечания

- База данных автоматически инициализируется при первом запуске с помощью скрипта `init.sql`
- Пользователь `user` создается с полными правами на базу `order_service_db`
- Таблицы создаются автоматически через миграции
- Вставляется тестовый заказ для демонстрации работы
- Go приложение автоматически восстанавливает кеш из БД при запуске
- Убедитесь, что порт 5432 свободен перед запуском контейнера
