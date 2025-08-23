const API_BASE_URL = 'http://localhost:8081/api/v1';

const searchForm = document.getElementById('searchForm');
const orderUIDInput = document.getElementById('orderUID');
const searchBtn = document.querySelector('.search-btn');
const resultsSection = document.getElementById('resultsSection');
const errorSection = document.getElementById('errorSection');
const errorMessage = document.getElementById('errorMessage');
const orderDetails = document.getElementById('orderDetails');
const loader = document.getElementById('loader');
const modal = document.getElementById('modal');
const modalTitle = document.getElementById('modalTitle');
const modalBody = document.getElementById('modalBody');

document.addEventListener('DOMContentLoaded', function () {
    // Обработчик формы
    searchForm.addEventListener('submit', handleSearch);

    // Обработчик Enter в поле ввода
    orderUIDInput.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleSearch(e);
        }
    });

    // Автофокус на поле ввода
    orderUIDInput.focus();

    // Обработчик для кнопки случайного заказа
    const randomOrderBtn = document.getElementById('randomOrderBtn');
    if (randomOrderBtn) {
        randomOrderBtn.addEventListener('click', generateRandomOrder);
    }
});

// Обработчик поиска заказа
async function handleSearch(e) {
    e.preventDefault();

    const orderUID = orderUIDInput.value.trim();
    if (!orderUID) {
        showError('Пожалуйста, введите ID заказа');
        return;
    }

    await searchOrder(orderUID);
}

// Поиск заказа через API
async function searchOrder(orderUID) {
    try {
        setLoading(true);
        hideError();
        hideResults();

        console.log(`Searching for order: ${orderUID}`);

        const response = await fetch(`${API_BASE_URL}/orders/${encodeURIComponent(orderUID)}`);
        const data = await response.json();

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

// Отображение информации о заказе
function displayOrder(order) {
    const orderInfo = document.createElement('div');
    orderInfo.className = 'order-info fade-in';

    orderInfo.innerHTML = `
        <!-- Основная информация о заказе -->
        <div class="order-section">
            <h3>📦 Основная информация</h3>
            <div class="order-grid">
                <div class="order-field">
                    <span class="field-label">ID заказа:</span>
                    <span class="field-value">${order.order_uid}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Трек-номер:</span>
                    <span class="field-value">${order.track_number}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Клиент:</span>
                    <span class="field-value">${order.customer_id}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Дата создания:</span>
                    <span class="field-value">${formatDate(order.date_created)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Служба доставки:</span>
                    <span class="field-value">${order.delivery_service}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Локаль:</span>
                    <span class="field-value">${order.locale || 'не указано'}</span>
                </div>
            </div>
        </div>

        <!-- Информация о доставке -->
        ${order.delivery ? `
        <div class="order-section">
            <h3>🚚 Доставка</h3>
            <div class="order-grid">
                <div class="order-field">
                    <span class="field-label">Получатель:</span>
                    <span class="field-value">${order.delivery.name}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Телефон:</span>
                    <span class="field-value">${order.delivery.phone}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Email:</span>
                    <span class="field-value">${order.delivery.email}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Адрес:</span>
                    <span class="field-value">${order.delivery.address}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Город:</span>
                    <span class="field-value">${order.delivery.city}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Регион:</span>
                    <span class="field-value">${order.delivery.region}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Индекс:</span>
                    <span class="field-value">${order.delivery.zip}</span>
                </div>
            </div>
        </div>
        ` : ''}

        <!-- Платежная информация -->
        ${order.payment ? `
        <div class="order-section">
            <h3>💳 Платеж</h3>
            <div class="order-grid">
                <div class="order-field">
                    <span class="field-label">Сумма заказа:</span>
                    <span class="field-value">${formatCurrency(order.payment.amount, order.payment.currency)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Стоимость товаров:</span>
                    <span class="field-value">${formatCurrency(order.payment.goods_total, order.payment.currency)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Доставка:</span>
                    <span class="field-value">${formatCurrency(order.payment.delivery_cost, order.payment.currency)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Провайдер:</span>
                    <span class="field-value">${order.payment.provider}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Банк:</span>
                    <span class="field-value">${order.payment.bank}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Транзакция:</span>
                    <span class="field-value">${order.payment.transaction}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Дата платежа:</span>
                    <span class="field-value">${formatTimestamp(order.payment.payment_dt)}</span>
                </div>
            </div>
        </div>
        ` : ''}

        <!-- Товары -->
        ${order.items && order.items.length > 0 ? `
        <div class="order-section">
            <h3>🛍️ Товары (${order.items.length})</h3>
            <div class="items-list">
                ${order.items.map(item => `
                    <div class="item-card">
                        <div class="item-header">
                            <div class="item-name">${item.name}</div>
                            <div class="item-price">${formatCurrency(item.total_price, order.payment?.currency || 'USD')}</div>
                        </div>
                        <div class="item-details">
                            <div><strong>Бренд:</strong> ${item.brand}</div>
                            <div><strong>Размер:</strong> ${item.size || 'не указан'}</div>
                            <div><strong>Артикул:</strong> ${item.nm_id}</div>
                            <div><strong>Цена:</strong> ${formatCurrency(item.price, order.payment?.currency || 'USD')}</div>
                            ${item.sale > 0 ? `<div><strong>Скидка:</strong> ${item.sale}%</div>` : ''}
                            <div><strong>Статус:</strong> <span class="status-badge ${getStatusClass(item.status)}">${getStatusText(item.status)}</span></div>
                        </div>
                    </div>
                `).join('')}
            </div>
        </div>
        ` : ''}
    `;

    orderDetails.innerHTML = '';
    orderDetails.appendChild(orderInfo);
    showResults();
}

// Вспомогательные функции для форматирования
function formatDate(dateString) {
    try {
        const date = new Date(dateString);
        return date.toLocaleString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    } catch {
        return dateString;
    }
}

function formatTimestamp(timestamp) {
    try {
        const date = new Date(timestamp * 1000);
        return date.toLocaleString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    } catch {
        return 'Неизвестно';
    }
}

function formatCurrency(amount, currency = 'USD') {
    if (typeof amount !== 'number') return amount;

    // Преобразуем центы в основную валюту
    const value = amount / 100;

    const formatter = new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 2
    });

    return formatter.format(value);
}

function getStatusClass(status) {
    if (status >= 200 && status < 300) return 'status-success';
    if (status >= 100 && status < 200) return 'status-warning';
    return 'status-error';
}

function getStatusText(status) {
    const statusMap = {
        100: 'Обработка',
        200: 'Подтвержден',
        202: 'Принят',
        300: 'Доставлен',
        400: 'Ошибка',
        500: 'Отменен'
    };
    return statusMap[status] || `Статус ${status}`;
}

// Управление состоянием UI
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

// Дополнительные функции
async function getRandomOrder() {
    // Список примеров заказов для демонстрации
    const exampleOrders = [
        'b563feb7b2b84b6test',
        'sample_order_001',
        'demo_order_123',
        'test_order_456'
    ];

    const randomOrder = exampleOrders[Math.floor(Math.random() * exampleOrders.length)];
    orderUIDInput.value = randomOrder;
    await searchOrder(randomOrder);
}

async function getCacheStats() {
    try {
        const response = await fetch(`${API_BASE_URL}/cache/stats`);
        const data = await response.json();

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
        const response = await fetch(`${API_BASE_URL}/orders?limit=10`);
        const data = await response.json();

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
        const response = await fetch(`${API_BASE_URL}/health`);
        const data = await response.json();

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

// Управление модальным окном
function showModal(title, content) {
    modalTitle.textContent = title;
    modalBody.innerHTML = content;
    modal.style.display = 'flex';
}

function closeModal() {
    modal.style.display = 'none';
}

// Закрытие модального окна по клику вне его
window.onclick = function (event) {
    if (event.target === modal) {
        closeModal();
    }
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

        const response = await fetch(`${API_BASE_URL}/orders/random`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const result = await response.json();

        if (result.success && result.data) {
            // Автоматически отображаем созданный заказ
            orderUIDInput.value = result.data.order_uid;
            displayOrder(result.data);
            showMessage(`Заказ создан: ${result.data.order_uid}`, 'success');
        } else {
            showError(result.error || 'Ошибка создания заказа');
        }
    } catch (error) {
        showMessage(`Ошибка: ${error.message}`, 'error');
    } finally {
        // Разблокируем кнопку
        randomBtn.disabled = false;
        randomBtn.style.backgroundColor = '';
        randomBtn.style.cursor = '';
        randomBtn.textContent = originalText;
    }
}

// Показ уведомлений
function showMessage(message, type = 'info') {
    // Создаем элемент уведомления
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;

    // Добавляем стили
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        z-index: 1000;
        opacity: 0;
        transform: translateX(100%);
        transition: all 0.3s ease;
    `;

    // Устанавливаем цвет в зависимости от типа
    switch (type) {
        case 'success':
            notification.style.backgroundColor = '#10b981';
            break;
        case 'error':
            notification.style.backgroundColor = '#ef4444';
            break;
        default:
            notification.style.backgroundColor = '#3b82f6';
    }

    // Добавляем в DOM
    document.body.appendChild(notification);

    // Анимация появления
    setTimeout(() => {
        notification.style.opacity = '1';
        notification.style.transform = 'translateX(0)';
    }, 100);

    // Автоматическое удаление через 3 секунды
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }, 3000);
}
