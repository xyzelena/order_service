#!/bin/bash

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8081/api/v1"
TEST_ORDER_ID="b563feb7b2b84b6test"

echo -e "${BLUE}🚀 Тестирование Order Service API${NC}"
echo "=================================="
echo ""

# Функция для проверки HTTP ответа
check_response() {
    local url=$1
    local expected_status=${2:-200}
    local description=$3
    
    echo -e "${YELLOW}📡 Тестирование: $description${NC}"
    echo "URL: $url"
    
    # Выполняем запрос и сохраняем результат
    response=$(curl -s -w "\n%{http_code}" "$url")
    http_code=$(echo "$response" | tail -n1)
    json_body=$(echo "$response" | head -n -1)
    
    echo "HTTP Status: $http_code"
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✅ Статус корректный${NC}"
    else
        echo -e "${RED}❌ Ожидался статус $expected_status, получен $http_code${NC}"
    fi
    
    # Проверяем валидность JSON
    if echo "$json_body" | jq . > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Ответ является валидным JSON${NC}"
        
        # Проверяем структуру ответа
        success=$(echo "$json_body" | jq -r '.success // "null"')
        if [ "$success" = "true" ] || [ "$success" = "false" ]; then
            echo -e "${GREEN}✅ Структура ответа корректная (поле success присутствует)${NC}"
        else
            echo -e "${RED}❌ Некорректная структура ответа (отсутствует поле success)${NC}"
        fi
        
        # Красиво выводим JSON
        echo -e "${BLUE}📄 Ответ:${NC}"
        echo "$json_body" | jq .
    else
        echo -e "${RED}❌ Ответ не является валидным JSON${NC}"
        echo "Тело ответа: $json_body"
    fi
    
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Проверяем доступность API
echo -e "${YELLOW}🔍 Проверка доступности API...${NC}"
if curl -s --connect-timeout 5 "$API_BASE/health" > /dev/null; then
    echo -e "${GREEN}✅ API доступен${NC}"
else
    echo -e "${RED}❌ API недоступен. Убедитесь, что сервис запущен на порту 8081${NC}"
    exit 1
fi
echo ""

# Тест 1: Health Check
check_response "$API_BASE/health" 200 "Health Check"

# Тест 2: Получение заказа по ID
check_response "$API_BASE/orders/$TEST_ORDER_ID" 200 "Получение заказа по ID"

# Тест 3: Получение несуществующего заказа
check_response "$API_BASE/orders/nonexistent_order" 404 "Получение несуществующего заказа"

# Тест 4: Получение списка заказов
check_response "$API_BASE/orders?limit=5" 200 "Получение списка заказов"

# Тест 5: Статистика кеша
check_response "$API_BASE/cache/stats" 200 "Статистика кеша"

# Тест 6: Некорректный endpoint
check_response "$API_BASE/invalid_endpoint" 404 "Некорректный endpoint"

# Тест производительности
echo -e "${YELLOW}⚡ Тест производительности (10 запросов)${NC}"
echo "========================================"

total_time=0
success_count=0

for i in {1..10}; do
    start_time=$(date +%s%N)
    
    response=$(curl -s -w "\n%{http_code}" "$API_BASE/orders/$TEST_ORDER_ID")
    http_code=$(echo "$response" | tail -n1)
    
    end_time=$(date +%s%N)
    duration=$((($end_time - $start_time) / 1000000)) # Переводим в миллисекунды
    
    total_time=$(($total_time + $duration))
    
    if [ "$http_code" -eq "200" ]; then
        success_count=$(($success_count + 1))
        echo -e "Запрос $i: ${GREEN}$duration мс${NC}"
    else
        echo -e "Запрос $i: ${RED}ОШИБКА (HTTP $http_code)${NC}"
    fi
done

if [ $success_count -gt 0 ]; then
    avg_time=$(($total_time / $success_count))
    echo ""
    echo -e "${BLUE}📊 Статистика производительности:${NC}"
    echo "Успешных запросов: $success_count/10"
    echo "Среднее время ответа: $avg_time мс"
    echo "Общее время: $total_time мс"
fi

echo ""
echo -e "${GREEN}🎉 Тестирование API завершено!${NC}"

# Дополнительная информация
echo ""
echo -e "${BLUE}📋 Краткая справка по API:${NC}"
echo "• GET /api/v1/orders/{order_uid} - получить заказ по ID"
echo "• GET /api/v1/orders?limit=N - получить список заказов"
echo "• GET /api/v1/cache/stats - статистика кеша"
echo "• GET /api/v1/health - проверка здоровья сервиса"
echo ""
echo -e "${BLUE}🌐 Frontend интерфейс:${NC}"
echo "• http://localhost:3000 - веб-интерфейс (если запущен)"
echo ""
echo -e "${BLUE}🔧 Kafka UI:${NC}"
echo "• http://localhost:8080 - интерфейс управления Kafka (если запущен)"
