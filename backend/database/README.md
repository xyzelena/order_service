# Database - Order Service

Папка содержит все файлы, связанные с базой данных PostgreSQL.

## Содержимое

- **docker-compose.yml** - конфигурация PostgreSQL контейнера
- **init.sql** - скрипт инициализации БД и пользователя
- **migrations/** - SQL миграции для создания таблиц
- **DATABASE_SCHEMA.md** - подробная документация схемы БД

## Быстрый запуск

```bash
# Запуск PostgreSQL
docker-compose up -d postgres

# Проверка подключения
docker exec -it order_service_postgres psql -U user -d order_service_db -c "\dt"
```

## Структура данных

База данных содержит 4 основные таблицы:
- `orders` - основная информация о заказах
- `deliveries` - данные о доставке
- `payments` - платежная информация  
- `order_items` - товары в заказе

## Примеры запросов

### Получение заказа с полной информацией:
```sql
SELECT 
    o.order_uid,
    o.track_number,
    d.name AS customer_name,
    d.city,
    p.amount,
    COUNT(i.id) AS items_count
FROM orders o
LEFT JOIN deliveries d ON o.order_uid = d.order_uid
LEFT JOIN payments p ON o.order_uid = p.order_uid
LEFT JOIN order_items i ON o.order_uid = i.order_uid
WHERE o.order_uid = 'b563feb7b2b84b6test'
GROUP BY o.order_uid, d.id, p.id;
```

### Статистика по заказам:
```sql
SELECT 
    COUNT(*) as total_orders,
    AVG(p.amount) as avg_amount,
    COUNT(DISTINCT o.customer_id) as unique_customers
FROM orders o
LEFT JOIN payments p ON o.order_uid = p.order_uid;
```
