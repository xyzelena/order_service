# Схема базы данных Order Service

## Обзор структуры

База данных спроектирована для хранения данных заказов с нормализованной структурой, состоящей из 4 основных таблиц:

```
orders (1) ──┬── deliveries (1)
             ├── payments (1)  
             └── order_items (N)
```

## Таблицы

### 1. `orders` - Основная таблица заказов

**Назначение**: Хранение ключевой информации о заказах

| Поле | Тип | Описание |
|------|-----|----------|
| `order_uid` | VARCHAR(255) PRIMARY KEY | Уникальный идентификатор заказа |
| `track_number` | VARCHAR(255) NOT NULL | Трек-номер для отслеживания |
| `entry` | VARCHAR(50) NOT NULL | Точка входа (например, WBIL) |
| `locale` | VARCHAR(10) | Локаль пользователя |
| `customer_id` | VARCHAR(255) NOT NULL | Идентификатор клиента |
| `delivery_service` | VARCHAR(100) NOT NULL | Служба доставки |
| `date_created` | TIMESTAMP WITH TIME ZONE | Дата создания заказа |

### 2. `deliveries` - Информация о доставке

**Назначение**: Хранение адресных данных и информации о получателе

| Поле | Тип | Описание |
|------|-----|----------|
| `order_uid` | VARCHAR(255) UNIQUE FK | Связь с заказом |
| `name` | VARCHAR(255) NOT NULL | Имя получателя |
| `phone` | VARCHAR(50) NOT NULL | Телефон получателя |
| `address` | TEXT NOT NULL | Полный адрес доставки |
| `city` | VARCHAR(100) NOT NULL | Город |
| `region` | VARCHAR(100) NOT NULL | Регион |
| `email` | VARCHAR(255) NOT NULL | Email получателя |

### 3. `payments` - Платежная информация

**Назначение**: Хранение финансовых данных по заказу

| Поле | Тип | Описание |
|------|-----|----------|
| `order_uid` | VARCHAR(255) UNIQUE FK | Связь с заказом |
| `transaction` | VARCHAR(255) NOT NULL | ID транзакции |
| `currency` | VARCHAR(10) NOT NULL | Валюта платежа |
| `provider` | VARCHAR(100) NOT NULL | Платежный провайдер |
| `amount` | INTEGER NOT NULL | Общая сумма (в копейках) |
| `delivery_cost` | INTEGER NOT NULL | Стоимость доставки |
| `goods_total` | INTEGER NOT NULL | Стоимость товаров |

### 4. `order_items` - Товары в заказе

**Назначение**: Хранение информации о товарах в заказе

| Поле | Тип | Описание |
|------|-----|----------|
| `order_uid` | VARCHAR(255) FK | Связь с заказом |
| `chrt_id` | BIGINT NOT NULL | ID характеристики товара |
| `name` | VARCHAR(255) NOT NULL | Название товара |
| `brand` | VARCHAR(255) NOT NULL | Бренд товара |
| `price` | INTEGER NOT NULL | Цена товара |
| `sale` | INTEGER | Скидка в процентах |
| `total_price` | INTEGER NOT NULL | Итоговая цена |
| `nm_id` | BIGINT NOT NULL | Номенклатурный номер |

## Индексы

### Производительность запросов оптимизирована индексами:

- **orders**: `track_number`, `customer_id`, `date_created`, `delivery_service`
- **deliveries**: `order_uid`, `city`, `region`  
- **payments**: `order_uid`, `transaction`, `provider`
- **order_items**: `order_uid`, `chrt_id`, `nm_id`, `brand`

## Типовые запросы

### Получение полной информации о заказе:
```sql
SELECT 
    o.*,
    d.name, d.phone, d.address, d.city,
    p.amount, p.currency, p.provider,
    array_agg(json_build_object(
        'name', i.name,
        'brand', i.brand, 
        'price', i.price,
        'total_price', i.total_price
    )) as items
FROM orders o
LEFT JOIN deliveries d ON o.order_uid = d.order_uid
LEFT JOIN payments p ON o.order_uid = p.order_uid  
LEFT JOIN order_items i ON o.order_uid = i.order_uid
WHERE o.order_uid = 'b563feb7b2b84b6test'
GROUP BY o.order_uid, d.id, p.id;
```

### Поиск заказов по клиенту:
```sql
SELECT order_uid, track_number, date_created, delivery_service
FROM orders 
WHERE customer_id = 'test'
ORDER BY date_created DESC;
```

