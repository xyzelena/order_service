# Order Service

Микросервис для обработки заказов с использованием PostgreSQL, Kafka и кеширования в памяти.

## Структура проекта

```
order_service/
├── backend/          # Backend приложение (Go, PostgreSQL, Kafka)
│   ├── app/         # Go микросервис
│   └── database/    # PostgreSQL конфигурация и схема
├── frontend/        # Frontend веб-интерфейс (HTML/CSS/JS)
│   ├── index.html   # Главная страница
│   ├── styles.css   # CSS стили
│   └── script.js    # JavaScript логика
└── README.md        # Основная документация
```

## Быстрый старт

### Автоматический запуск (рекомендуется)
```bash
# Запуск всей системы одной командой
make setup  # Только первый раз
make start  # Запуск инфраструктуры + сервисов

# Проверка статуса
make status

# Тестирование
make test
```

### Ручной запуск
```bash
# 1. Инфраструктура (PostgreSQL + Kafka)
cd backend/database
docker-compose up -d

# 2. Go приложение
cd ../app
make deps && make build && make run

# 3. Frontend (опционально)
# Go сервер (рекомендуется):
cd ../app && make run-frontend
# Или Python: cd ../../frontend && python -m http.server 3000
```

## Архитектура

Сервис состоит из:
- **Backend**: Go микросервис с REST API, Kafka consumer, PostgreSQL
- **Frontend**: Современный веб-интерфейс для поиска заказов (HTML/CSS/JS)
- **База данных**: PostgreSQL с нормализованной схемой
- **Очередь сообщений**: Kafka для получения данных заказов
- **Кеширование**: LRU кеш в памяти для быстрого доступа

## Доступные интерфейсы

### Веб-интерфейс (рекомендуется)
- **URL**: http://localhost:3000 (после запуска frontend)
- **Функции**: Поиск заказов, детальная информация, статистика
- **Технологии**: HTML5, CSS3, JavaScript

### REST API  
- **URL**: http://localhost:8081/api/v1
- **Документация**: http://localhost:8081 (информация об эндпоинтах)
- **Формат**: JSON API с CORS поддержкой

### Kafka UI
- **URL**: http://localhost:8080 (мониторинг сообщений)
- **Функции**: Просмотр топиков, сообщений, consumer groups

## Тестирование

### Быстрая проверка
```bash
make quick-test  # Проверка основных функций
```

### Полное тестирование
```bash
make test-api     # HTTP API и JSON формат
make test-cache   # Производительность кеша  
make test-kafka   # Kafka интеграция
```

### Ожидаемые результаты:
1. **API отвечает** на `GET http://localhost:8081/order/<order_uid>`
2. **JSON корректный** с полями success, data, error
3. **Кеш ускоряет** повторные запросы минимум в 5 раз
4. **Kafka обрабатывает** сообщения онлайн

Подробное руководство: [TESTING_GUIDE.md](TESTING_GUIDE.md)

