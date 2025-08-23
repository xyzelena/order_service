# Makefile для Order Service

.PHONY: build run test clean docker-build docker-run deps help run-frontend lint fmt mod-verify dev-setup api-test random-order

# Переменные
APP_NAME=order-service
DOCKER_IMAGE=$(APP_NAME):latest
GO_VERSION=1.21

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

deps: ## Установить зависимости
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	go mod download
	go mod tidy

build: deps ## Собрать приложение
	@echo "$(GREEN)Сборка приложения...$(NC)"
	go build -o bin/$(APP_NAME) ./cmd/server

run: build ## Запустить приложение
	@echo "$(GREEN)Запуск приложения...$(NC)"
	./bin/$(APP_NAME)

run-frontend: ## Запустить frontend сервер на Go
	@echo "$(GREEN)Запуск frontend сервера...$(NC)"
	@if [ -d "../../frontend" ]; then \
		go run ./cmd/frontend-server ../../frontend; \
	else \
		echo "$(RED)Папка frontend не найдена в ../../frontend$(NC)"; \
		echo "$(YELLOW)Проверьте структуру проекта$(NC)"; \
	fi

test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	go test -v ./...

clean: ## Очистить собранные файлы
	@echo "$(GREEN)Очистка...$(NC)"
	rm -rf bin/
	go clean

docker-build: ## Собрать Docker образ
	@echo "$(GREEN)Сборка Docker образа...$(NC)"
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build ## Запустить в Docker контейнере
	@echo "$(GREEN)Запуск в Docker...$(NC)"
	docker run --rm -p 8081:8081 --env-file env.example $(DOCKER_IMAGE)

lint: ## Проверить код линтером
	@echo "$(GREEN)Проверка кода...$(NC)"
	golangci-lint run

fmt: ## Форматировать код
	@echo "$(GREEN)Форматирование кода...$(NC)"
	go fmt ./...

mod-verify: ## Проверить модули
	@echo "$(GREEN)Проверка модулей...$(NC)"
	go mod verify

dev-setup: ## Настройка среды разработки
	@echo "$(GREEN)Настройка среды разработки...$(NC)"
	@if [ ! -f env.example ]; then \
		echo "$(RED)Файл env.example не найден!$(NC)"; \
		exit 1; \
	fi
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "$(GREEN)Файл .env создан из env.example$(NC)"; \
		echo "$(YELLOW)Настройте переменные в .env при необходимости$(NC)"; \
	else \
		echo "$(YELLOW)Файл .env уже существует$(NC)"; \
	fi

api-test: ## Тестировать API эндпоинты
	@echo "$(GREEN)Тестирование API...$(NC)"
	@echo "$(YELLOW)Проверка health check...$(NC)"
	@curl -s http://localhost:8081/api/v1/health | jq '.success' || echo "$(RED)API недоступен$(NC)"
	@echo "$(YELLOW)Проверка поиска заказа...$(NC)"
	@curl -s http://localhost:8081/api/v1/orders/b563feb7b2b84b6test | jq '.success' || echo "$(RED)Заказ не найден$(NC)"
	@echo "$(YELLOW)Проверка статистики кеша...$(NC)"
	@curl -s http://localhost:8081/api/v1/cache/stats | jq '.success' || echo "$(RED)Кеш недоступен$(NC)"

random-order: ## Создать случайный заказ
	@echo "$(GREEN)Создание случайного заказа...$(NC)"
	@curl -s -X POST http://localhost:8081/api/v1/orders/random | jq '.data.order_uid' || echo "$(RED)Ошибка создания заказа$(NC)"

# По умолчанию показываем справку
.DEFAULT_GOAL := help
