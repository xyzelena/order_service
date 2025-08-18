# Makefile для Order Service

.PHONY: build run test clean docker-build docker-run deps help

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
	go run ./cmd/frontend-server ../../frontend

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
	docker run --rm -p 8080:8080 --env-file env.example $(DOCKER_IMAGE)

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
	@echo "$(YELLOW)Скопируйте env.example в .env и настройте переменные$(NC)"

# По умолчанию показываем справку
.DEFAULT_GOAL := help
