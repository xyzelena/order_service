# Wildberries Go Order Service

Микросервис для обработки заказов с использованием PostgreSQL, Kafka и кеширования в памяти.

## Структура проекта

```
order_service/
├── backend/          # Backend приложение (Go, PostgreSQL, Kafka)
├── frontend/         # Frontend приложение (веб-интерфейс)
└── README.md         # Основная документация
```

## Быстрый старт

### Backend
```bash
cd backend
# См. backend/README.md для детальных инструкций
```

### Frontend  
```bash
cd frontend
# См. frontend/README.md для детальных инструкций
```

## Архитектура

Сервис состоит из:
- **Backend**: Go-приложение с REST API, обработка Kafka сообщений, PostgreSQL
- **Frontend**: Веб-интерфейс для отображения данных заказов
- **База данных**: PostgreSQL для хранения заказов
- **Очередь сообщений**: Kafka для получения данных заказов
- **Кеширование**: In-memory кеш для быстрого доступа к данным

