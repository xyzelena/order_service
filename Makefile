# Makefile для Order Service
# Автоматизация запуска и тестирования всей системы

.PHONY: help setup start stop test test-api test-cache test-kafka clean status

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
BLUE=\033[0;34m
NC=\033[0m

help: ## Показать справку по командам
	@echo "$(GREEN)Order Service - Команды управления$(NC)"
	@echo "====================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

setup: ## Настройка окружения (только первый раз)
	@echo "$(BLUE) Настройка окружения...$(NC)"
	@cd backend/app && cp env.example .env 2>/dev/null || echo "Файл .env уже существует"
	@cd backend/database && docker-compose pull
	@cd backend/app && make deps
	@cd backend/app && make build
	@echo "$(GREEN) Окружение настроено$(NC)"

start: ## Запустить всю систему (инфраструктура + сервисы)
	@echo "$(BLUE) Запуск системы...$(NC)"
	@echo "$(YELLOW) Запуск инфраструктуры (PostgreSQL + Kafka)...$(NC)"
	@cd backend/database && docker-compose up -d
	@echo "$(YELLOW) Ожидание готовности инфраструктуры...$(NC)"
	@sleep 15
	@echo "$(YELLOW) Сборка и запуск Go приложения...$(NC)"
	@cd backend/app && make build && make run &
	@echo "$(GREEN) Система запущена!$(NC)"
	@echo ""
	@echo "$(BLUE) Доступные сервисы:$(NC)"
	@echo "  • API: http://localhost:8081/api/v1"
	@echo "  • Health: http://localhost:8081/api/v1/health"
	@echo "  • Kafka UI: http://localhost:8080"
	@echo "  • Frontend: make start-frontend"

start-app: ## Запустить только Go приложение (инфраструктура должна быть запущена)
	@echo "$(BLUE) Запуск Go приложения...$(NC)"
	@cd backend/app && make run

start-producer: ## Запустить Kafka producer для эмуляции сообщений
	@echo "$(BLUE) Запуск Kafka producer...$(NC)"
	@cd backend/database && go mod tidy && go run kafka-producer.go

start-frontend: ## Запустить frontend на порту 3000
	@echo "$(BLUE) Запуск frontend...$(NC)"
	@echo "$(YELLOW)Доступные варианты:$(NC)"
	@echo "  1. Go сервер:     cd backend/app && make run-frontend"
	@echo "  2. Python сервер: cd frontend && python -m http.server 3000"
	@echo "  3. Node.js:       npx serve frontend -p 3000"
	@echo ""
	@echo "$(BLUE)Запускаем Go сервер...$(NC)"
	@cd backend/app && make run-frontend

stop: ## Остановить всю систему
	@echo "$(YELLOW) Остановка системы...$(NC)"
	@cd backend/database && docker-compose down
	@pkill -f "order-service" 2>/dev/null || true
	@echo "$(GREEN) Система остановлена$(NC)"

status: ## Показать статус всех компонентов
	@echo "$(BLUE) Статус системы$(NC)"
	@echo "=================="
	@echo ""
	@echo "$(YELLOW) Docker контейнеры:$(NC)"
	@cd backend/database && docker-compose ps || echo "Docker compose не запущен"
	@echo ""
	@echo "$(YELLOW) Сетевые порты:$(NC)"
	@echo "PostgreSQL (5432):" && (lsof -i :5432 >/dev/null 2>&1 && echo "$(GREEN) Активен$(NC)" || echo "$(RED) Не активен$(NC)")
	@echo "Kafka (9092):     " && (lsof -i :9092 >/dev/null 2>&1 && echo "$(GREEN) Активен$(NC)" || echo "$(RED) Не активен$(NC)")
	@echo "API (8081):       " && (lsof -i :8081 >/dev/null 2>&1 && echo "$(GREEN) Активен$(NC)" || echo "$(RED) Не активен$(NC)")
	@echo "Frontend (3000):  " && (lsof -i :3000 >/dev/null 2>&1 && echo "$(GREEN) Активен$(NC)" || echo "$(RED) Не активен$(NC)")

test: test-api test-cache ## Запустить все тесты
	@echo "$(GREEN) Все тесты завершены!$(NC)"

test-api: ## Тестировать HTTP API и JSON формат
	@echo "$(BLUE) Тестирование HTTP API...$(NC)"
	@if [ -f backend/app/tests/api_test.sh ]; then \
		cd backend/app/tests && chmod +x api_test.sh && ./api_test.sh; \
	else \
		echo "$(YELLOW) Тестовый скрипт не найден, запускаем простую проверку...$(NC)"; \
		curl -s http://localhost:8081/api/v1/health || echo "$(RED) API недоступен$(NC)"; \
	fi

test-cache: ## Тестировать производительность кеша
	@echo "$(BLUE) Тестирование производительности кеша...$(NC)"
	@if [ -f backend/app/tests/cache_benchmark.go ]; then \
		cd backend/app/tests && go run cache_benchmark.go; \
	else \
		echo "$(YELLOW) Бенчмарк кеша не найден, пропускаем...$(NC)"; \
	fi

test-kafka: ## Проверить работу Kafka (отправить тестовое сообщение)
	@echo "$(BLUE) Тестирование Kafka интеграции...$(NC)"
	@cd backend/database && timeout 10 go run kafka-producer.go || echo "$(YELLOW) Тест Kafka завершен по таймауту$(NC)"

test-full: ## Полный интеграционный тест всей системы
	@echo "$(BLUE) Полное интеграционное тестирование...$(NC)"
	@$(MAKE) status
	@$(MAKE) test-api
	@$(MAKE) test-cache
	@echo "$(GREEN) Интеграционное тестирование завершено$(NC)"

clean: ## Очистить систему (остановить контейнеры, удалить volumes)
	@echo "$(YELLOW) Очистка системы...$(NC)"
	@cd backend/database && docker-compose down -v
	@docker system prune -f >/dev/null 2>&1 || true
	@echo "$(GREEN) Система очищена$(NC)"

logs: ## Показать логи всех компонентов
	@echo "$(BLUE) Логи системы$(NC)"
	@echo "==============="
	@cd backend/database && docker-compose logs --tail=50

restart: stop start ## Перезапустить всю систему
	@echo "$(GREEN) Система перезапущена$(NC)"

quick-test: ## Быстрая проверка основных функций
	@echo "$(BLUE) Быстрая проверка системы...$(NC)"
	@echo "$(YELLOW)1. Проверка API...$(NC)"
	@curl -s http://localhost:8081/api/v1/health >/dev/null && echo "$(GREEN) API работает$(NC)" || echo "$(RED) API недоступен$(NC)"
	@echo "$(YELLOW)2. Проверка заказа...$(NC)"
	@curl -s http://localhost:8081/api/v1/orders/b563feb7b2b84b6test >/dev/null && echo "$(GREEN) Заказы загружаются$(NC)" || echo "$(RED) Ошибка загрузки заказов$(NC)"
	@echo "$(YELLOW)3. Проверка кеша...$(NC)"
	@curl -s http://localhost:8081/api/v1/cache/stats >/dev/null && echo "$(GREEN) Кеш работает$(NC)" || echo "$(RED) Кеш недоступен$(NC)"
	@echo "$(GREEN) Быстрая проверка завершена$(NC)"

# Разработка
dev: ## Запустить в режиме разработки с автоперезагрузкой
	@echo "$(BLUE) Режим разработки$(NC)"
	@cd backend/app && DEBUG=true make run

build: ## Собрать все компоненты
	@echo "$(BLUE) Сборка проекта...$(NC)"
	@cd backend/app && make build
	@cd backend/database && go mod tidy
	@echo "$(GREEN) Сборка завершена$(NC)"

# По умолчанию показываем справку
.DEFAULT_GOAL := help
