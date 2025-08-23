import { formatDate } from './utils.js';

import { showModal, closeModal } from './modal.js';
import { showMessage } from './notifications.js';
import { displayOrder } from './orderRenderer.js';
import {
    searchOrder as apiSearchOrder,
    getCacheStats as apiGetCacheStats,
    getAllOrders as apiGetAllOrders,
    healthCheck as apiHealthCheck,
    createRandomOrder,
    getRandomOrderExample
} from './api.js';

const searchForm = document.getElementById('searchForm');
const orderUIDInput = document.getElementById('orderUID');
const searchBtn = document.querySelector('.search-btn');
const resultsSection = document.getElementById('resultsSection');
const errorSection = document.getElementById('errorSection');
const errorMessage = document.getElementById('errorMessage');
const orderDetails = document.getElementById('orderDetails');
const loader = document.getElementById('loader');


document.addEventListener('DOMContentLoaded', function () {
    searchForm.addEventListener('submit', handleSearch);

    orderUIDInput.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleSearch(e);
        }
    });

    orderUIDInput.focus();

    const randomOrderBtn = document.getElementById('randomOrderBtn');
    if (randomOrderBtn) {
        randomOrderBtn.addEventListener('click', generateRandomOrder);
    }
});

async function handleSearch(e) {
    e.preventDefault();

    const orderUID = orderUIDInput.value.trim();
    if (!orderUID) {
        showError('Пожалуйста, введите ID заказа');
        return;
    }

    await searchOrder(orderUID);
}

async function searchOrder(orderUID) {
    try {
        setLoading(true);
        hideError();
        hideResults();

        console.log(`Searching for order: ${orderUID}`);

        const data = await apiSearchOrder(orderUID);

        if (data.success && data.data) {
            displayOrder(data.data);
            console.log('Order found:', data.data);
        } else {
            showError(data.error || 'Заказ не найден');
        }
    } catch (error) {
        console.error('Search error:', error);
        showError(`Ошибка сети: ${error.message}. Убедитесь, что backend сервис запущен на порту 8080.`);
    } finally {
        setLoading(false);
    }
}

function setLoading(loading) {
    if (loading) {
        searchBtn.classList.add('loading');
        searchBtn.disabled = true;
    } else {
        searchBtn.classList.remove('loading');
        searchBtn.disabled = false;
    }
}

function showResults() {
    resultsSection.style.display = 'block';
    resultsSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

function hideResults() {
    resultsSection.style.display = 'none';
}

function showError(message) {
    errorMessage.textContent = message;
    errorSection.style.display = 'block';
    errorSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

function hideError() {
    errorSection.style.display = 'none';
}

function clearError() {
    hideError();
    orderUIDInput.focus();
}

async function getRandomOrder() {
    const randomOrder = getRandomOrderExample();
    orderUIDInput.value = randomOrder;
    await searchOrder(randomOrder);
}

async function getCacheStats() {
    try {
        const data = await apiGetCacheStats();

        if (data.success) {
            showModal('Статистика кеша', `
                <div class="order-grid">
                    <div class="order-field">
                        <span class="field-label">Размер кеша:</span>
                        <span class="field-value">${data.data.size}</span>
                    </div>
                    <div class="order-field">
                        <span class="field-label">Максимальный размер:</span>
                        <span class="field-value">${data.data.capacity}</span>
                    </div>
                    <div class="order-field">
                        <span class="field-label">Загруженность:</span>
                        <span class="field-value">${Math.round((data.data.size / data.data.capacity) * 100)}%</span>
                    </div>
                </div>
            `);
        } else {
            showError(data.error || 'Не удалось получить статистику кеша');
        }
    } catch (error) {
        showError(`Ошибка получения статистики: ${error.message}`);
    }
}

async function getAllOrders() {
    try {
        const data = await apiGetAllOrders(10);

        if (data.success && data.data.orders) {
            const ordersList = data.data.orders.map(order => `
                <div class="order-field" style="cursor: pointer;" onclick="loadOrder('${order.order_uid}')">
                    <span class="field-label">${order.order_uid}</span>
                    <span class="field-value">${formatDate(order.date_created)}</span>
                </div>
            `).join('');

            showModal(`Последние заказы (${data.data.count})`, `
                <div style="max-height: 400px; overflow-y: auto;">
                    ${ordersList || '<p>Заказы не найдены</p>'}
                </div>
                <p style="margin-top: 16px; color: #6b7280; font-size: 0.9rem;">
                    Нажмите на заказ, чтобы загрузить его
                </p>
            `);
        } else {
            showError(data.error || 'Не удалось получить список заказов');
        }
    } catch (error) {
        showError(`Ошибка получения заказов: ${error.message}`);
    }
}

async function healthCheck() {
    try {
        const data = await apiHealthCheck();

        if (data.success) {
            const status = data.data;
            showModal('Состояние сервиса', `
                <div class="order-grid">
                    <div class="order-field">
                        <span class="field-label">Статус:</span>
                        <span class="field-value status-badge status-success">${status.status}</span>
                    </div>
                    <div class="order-field">
                        <span class="field-label">Кеш:</span>
                        <span class="field-value">${status.cache.size}/${status.cache.capacity}</span>
                    </div>
                    <div class="order-field">
                        <span class="field-label">Время проверки:</span>
                        <span class="field-value">${new Date().toLocaleString('ru-RU')}</span>
                    </div>
                </div>
            `);
        } else {
            showError('Сервис недоступен');
        }
    } catch (error) {
        showError(`Сервис недоступен: ${error.message}`);
    }
}

// Загрузка конкретного заказа из списка
async function loadOrder(orderUID) {
    closeModal();
    orderUIDInput.value = orderUID;
    await searchOrder(orderUID);
}

// Генерация случайного заказа
async function generateRandomOrder(event) {
    const randomBtn = event ? event.target : document.getElementById('randomOrderBtn');

    if (!randomBtn) {
        return;
    }

    const originalText = randomBtn.textContent;

    try {
        // Блокируем кнопку и меняем текст
        randomBtn.disabled = true;
        randomBtn.style.backgroundColor = '#9ca3af';
        randomBtn.style.cursor = 'not-allowed';
        randomBtn.textContent = 'Создается...';

        hideError();

        const result = await createRandomOrder();

        if (result.success && result.data) {
            orderUIDInput.value = result.data.order_uid;
            displayOrder(result.data);
            showMessage(`Заказ создан: ${result.data.order_uid}`, 'success');
        } else {
            showError(result.error || 'Ошибка создания заказа');
        }
    } catch (error) {
        showMessage(`Ошибка: ${error.message}`, 'error');
    } finally {
        randomBtn.disabled = false;
        randomBtn.style.backgroundColor = '';
        randomBtn.style.cursor = '';
        randomBtn.textContent = originalText;
    }
}


