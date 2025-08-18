-- Миграция для создания таблиц заказов
-- Версия: 001
-- Описание: Создание основных таблиц для хранения данных заказов

-- 1. Таблица заказов (основная информация)
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,           -- Уникальный идентификатор заказа
    track_number VARCHAR(255) NOT NULL,           -- Трек-номер заказа
    entry VARCHAR(50) NOT NULL,                   -- Точка входа (WBIL)
    locale VARCHAR(10) DEFAULT 'ru',              -- Локаль
    internal_signature TEXT,                      -- Внутренняя подпись
    customer_id VARCHAR(255) NOT NULL,            -- ID клиента
    delivery_service VARCHAR(100) NOT NULL,       -- Служба доставки
    shardkey VARCHAR(10) NOT NULL,                -- Ключ шарда
    sm_id INTEGER NOT NULL,                       -- SM ID
    date_created TIMESTAMP WITH TIME ZONE NOT NULL, -- Дата создания заказа
    oof_shard VARCHAR(10) NOT NULL,               -- OOF шард
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Время записи в БД
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- Время обновления
);

-- 2. Таблица информации о доставке
CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,                        -- Автоинкрементный ID
    order_uid VARCHAR(255) NOT NULL UNIQUE REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,                   -- Имя получателя
    phone VARCHAR(50) NOT NULL,                   -- Телефон получателя
    zip VARCHAR(20) NOT NULL,                     -- Почтовый индекс
    city VARCHAR(100) NOT NULL,                   -- Город
    address TEXT NOT NULL,                        -- Адрес доставки
    region VARCHAR(100) NOT NULL,                 -- Регион
    email VARCHAR(255) NOT NULL,                  -- Email получателя
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Таблица платежной информации
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,                        -- Автоинкрементный ID
    order_uid VARCHAR(255) NOT NULL UNIQUE REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR(255) NOT NULL,            -- ID транзакции
    request_id VARCHAR(255),                      -- ID запроса (может быть пустым)
    currency VARCHAR(10) NOT NULL,                -- Валюта платежа
    provider VARCHAR(100) NOT NULL,               -- Платежный провайдер
    amount INTEGER NOT NULL,                      -- Общая сумма (в копейках/центах)
    payment_dt BIGINT NOT NULL,                   -- Unix timestamp платежа
    bank VARCHAR(100) NOT NULL,                   -- Банк
    delivery_cost INTEGER NOT NULL,               -- Стоимость доставки
    goods_total INTEGER NOT NULL,                 -- Стоимость товаров
    custom_fee INTEGER DEFAULT 0,                 -- Таможенный сбор
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица товаров в заказе
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,                        -- Автоинкрементный ID
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,                      -- ID характеристики товара
    track_number VARCHAR(255) NOT NULL,           -- Трек-номер (дублируется из заказа)
    price INTEGER NOT NULL,                       -- Цена товара
    rid VARCHAR(255) NOT NULL,                    -- RID товара
    name VARCHAR(255) NOT NULL,                   -- Название товара
    sale INTEGER DEFAULT 0,                       -- Скидка в процентах
    size VARCHAR(50) NOT NULL,                    -- Размер товара
    total_price INTEGER NOT NULL,                 -- Итоговая цена товара
    nm_id BIGINT NOT NULL,                        -- Номенклатурный номер
    brand VARCHAR(255) NOT NULL,                  -- Бренд товара
    status INTEGER NOT NULL,                      -- Статус товара
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации запросов
CREATE INDEX idx_orders_track_number ON orders(track_number);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_date_created ON orders(date_created);
CREATE INDEX idx_orders_delivery_service ON orders(delivery_service);

CREATE INDEX idx_deliveries_order_uid ON deliveries(order_uid);
CREATE INDEX idx_deliveries_city ON deliveries(city);
CREATE INDEX idx_deliveries_region ON deliveries(region);

CREATE INDEX idx_payments_order_uid ON payments(order_uid);
CREATE INDEX idx_payments_transaction ON payments(transaction);
CREATE INDEX idx_payments_provider ON payments(provider);

CREATE INDEX idx_order_items_order_uid ON order_items(order_uid);
CREATE INDEX idx_order_items_chrt_id ON order_items(chrt_id);
CREATE INDEX idx_order_items_nm_id ON order_items(nm_id);
CREATE INDEX idx_order_items_brand ON order_items(brand);

-- Комментарии к таблицам
COMMENT ON TABLE orders IS 'Основная таблица заказов с ключевой информацией';
COMMENT ON TABLE deliveries IS 'Информация о доставке заказов';
COMMENT ON TABLE payments IS 'Платежная информация заказов';
COMMENT ON TABLE order_items IS 'Товары, входящие в заказы';

-- Комментарии к ключевым полям
COMMENT ON COLUMN orders.order_uid IS 'Уникальный идентификатор заказа, первичный ключ';
COMMENT ON COLUMN payments.amount IS 'Сумма в минимальных единицах валюты (копейки, центы)';
COMMENT ON COLUMN payments.payment_dt IS 'Unix timestamp времени платежа';
COMMENT ON COLUMN order_items.sale IS 'Размер скидки в процентах';
