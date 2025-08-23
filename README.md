# Order Service

Микросервис для обработки заказов с модульной архитектурой, использующий PostgreSQL, Kafka и кеширование в памяти.

## Структура проекта

```
order_service/
├── backend/                # Backend приложение (Go, PostgreSQL, Kafka)
│   ├── app/               # Go микросервис с REST API
│   │   ├── cmd/           # Точки входа приложения
│   │   ├── internal/      # Внутренняя логика
│   │   ├── pkg/           # Общие пакеты
│   │   └── Makefile       # Команды сборки и запуска
│   └── database/          # PostgreSQL конфигурация и схема
│       ├── migrations/    # SQL миграции
│       └── docker-compose.yml
├── frontend/              # Модульный frontend (ES6, HTML/CSS/JS)
│   ├── index.html        # Главная страница
│   ├── styles/           # CSS стили
│   │   └── main.css
│   ├── scripts/          # JavaScript ES6 модули
│   │   ├── main.js       # Основная логика
│   │   ├── api.js        # API взаимодействие
│   │   ├── orderRenderer.js  # Отображение заказов
│   │   ├── modal.js      # Модальные окна
│   │   ├── notifications.js  # Уведомления
│   │   └── utils.js      # Утилиты
│   └── README.md         # Документация frontend
├── Makefile              # Команды управления всей системой
├── README.md             # Основная документация
└── TESTING_GUIDE.md      # Руководство по тестированию
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
# Или из корня: make start-frontend
```

## Архитектура

Сервис состоит из:
- **Backend**: Go микросервис с REST API, Kafka consumer, PostgreSQL
- **Frontend**: Модульный веб-интерфейс на ES6 модулях (HTML/CSS/JS)
- **База данных**: PostgreSQL с нормализованной схемой и миграциями
- **Очередь сообщений**: Kafka для получения данных заказов в реальном времени
- **Кеширование**: LRU кеш в памяти для быстрого доступа к заказам

## Доступные интерфейсы

### Веб-интерфейс (рекомендуется)
- **URL**: http://localhost:3000 (после запуска frontend)
- **Функции**: Поиск заказов, детальная информация, статистика, генерация тестовых заказов
- **Технологии**: HTML5, CSS3, ES6 Modules, модульная архитектура
- **Особенности**: Модальные окна, уведомления, адаптивный дизайн

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
make test           # Запустить все тесты (API + кеш)
make test-api       # HTTP API и JSON формат
make test-cache     # Производительность кеша  
make test-kafka     # Kafka интеграция
make test-full      # Полный интеграционный тест
make check-frontend # Проверка модульной структуры frontend
```

### Ожидаемые результаты:
1. **API отвечает** на `GET http://localhost:8081/api/v1/orders/<order_uid>`
2. **JSON корректный** с полями success, data, error
3. **Кеш ускоряет** повторные запросы минимум в 5 раз
4. **Kafka обрабатывает** сообщения онлайн
5. **Frontend модули** загружаются и работают корректно
6. **Генерация заказов** работает через `POST /api/v1/orders/random`

Подробное руководство: [TESTING_GUIDE.md](TESTING_GUIDE.md)

