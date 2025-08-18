CREATE DATABASE order_service_db;

CREATE USER "user" WITH PASSWORD '0000';

GRANT ALL PRIVILEGES ON DATABASE order_service_db TO "user";

\c order_service_db;

GRANT ALL ON SCHEMA public TO "user";

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "user";
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO "user";

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO "user";

-- Создание таблиц для заказов
\i /docker-entrypoint-initdb.d/migrations/001_create_orders_tables.sql

-- Вставка тестовых данных
\i /docker-entrypoint-initdb.d/migrations/002_insert_test_data.sql
