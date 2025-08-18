# Frontend - Order Service

Frontend веб-интерфейса для отображения данных заказов.

## Описание

Веб-интерфейс для просмотра и управления данными заказов, получаемых от backend сервиса.

## Планируемый технологический стек

- **Framework**: React/Vue.js (на выбор)
- **Язык**: TypeScript/JavaScript
- **Стилизация**: CSS/SCSS или CSS-framework
- **HTTP клиент**: Axios/Fetch API
- **Build tool**: Vite/Webpack

## Функциональность

Планируемые возможности интерфейса:
- Отображение списка заказов
- Поиск и фильтрация заказов
- Просмотр детальной информации о заказе
- Real-time обновления данных
- Responsive дизайн

## Структура файлов

```
frontend/
├── src/                 # Исходный код приложения
├── public/             # Статические файлы
├── package.json        # Зависимости и скрипты
└── README.md           # Документация frontend
```

## Быстрый старт

```bash
# Установка зависимостей (после создания package.json)
npm install

# Запуск dev сервера
npm run dev

# Сборка для production
npm run build
```

## API Integration

Frontend будет взаимодействовать с backend API:
- **Base URL**: http://localhost:8080 (предполагаемый)
- **Endpoints**: REST API для работы с заказами
