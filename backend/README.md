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

# Если Docker не запущен, запустите Docker Desktop
# или выполните: sudo systemctl start docker (на Linux)
```

### 2. Запуск PostgreSQL

```bash
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
# Остановка контейнера
docker-compose down

# Подключение к БД через psql
docker exec -it order_service_postgres psql -U user -d order_service_db

# Просмотр логов
docker-compose logs postgres
```

## Структура файлов

```
backend/
├── docker-compose.yml    # Конфигурация PostgreSQL
├── init.sql             # SQL скрипт инициализации БД
└── README.md            # Документация backend
```

## Решение проблем

### Ошибка "Cannot connect to the Docker daemon"
```bash
# Проверьте, запущен ли Docker
docker --version

# Запустите Docker Desktop (macOS/Windows) или Docker daemon (Linux)
# macOS/Windows: откройте Docker Desktop
# Linux: sudo systemctl start docker
```

### Warning о версии docker-compose
Предупреждение `version is obsolete` - не критично, но для чистоты убрана строка `version` из docker-compose.yml

## Примечания

- База данных автоматически инициализируется при первом запуске с помощью скрипта `init.sql`
- Пользователь `user` создается с полными правами на базу `order_service_db`
- Убедитесь, что порт 5432 свободен перед запуском контейнера
