// Модуль для отображения заказов

import {
    formatDate,
    formatTimestamp,
    formatCurrency,
    getStatusClass,
    getStatusText
} from './utils.js';

const renderOrderBasicInfo = (order) => `
    <div class="order-section">
        <h3>Основная информация</h3>
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
`;

const renderDeliveryInfo = (delivery) => {
    if (!delivery) return '';

    return `
        <div class="order-section">
            <h3>Доставка</h3>
            <div class="order-grid">
                <div class="order-field">
                    <span class="field-label">Получатель:</span>
                    <span class="field-value">${delivery.name}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Телефон:</span>
                    <span class="field-value">${delivery.phone}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Email:</span>
                    <span class="field-value">${delivery.email}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Адрес:</span>
                    <span class="field-value">${delivery.address}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Город:</span>
                    <span class="field-value">${delivery.city}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Регион:</span>
                    <span class="field-value">${delivery.region}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Индекс:</span>
                    <span class="field-value">${delivery.zip}</span>
                </div>
            </div>
        </div>
    `;
};


const renderPaymentInfo = (payment) => {
    if (!payment) return '';

    return `
        <div class="order-section">
            <h3>Платеж</h3>
            <div class="order-grid">
                <div class="order-field">
                    <span class="field-label">Сумма заказа:</span>
                    <span class="field-value">${formatCurrency(payment.amount, payment.currency)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Стоимость товаров:</span>
                    <span class="field-value">${formatCurrency(payment.goods_total, payment.currency)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Доставка:</span>
                    <span class="field-value">${formatCurrency(payment.delivery_cost, payment.currency)}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Провайдер:</span>
                    <span class="field-value">${payment.provider}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Банк:</span>
                    <span class="field-value">${payment.bank}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Транзакция:</span>
                    <span class="field-value">${payment.transaction}</span>
                </div>
                <div class="order-field">
                    <span class="field-label">Дата платежа:</span>
                    <span class="field-value">${formatTimestamp(payment.payment_dt)}</span>
                </div>
            </div>
        </div>
    `;
};

const renderItemsList = (items, currency = 'USD') => {
    if (!items || items.length === 0) return '';

    return `
        <div class="order-section">
            <h3>Товары (${items.length})</h3>
            <div class="items-list">
                ${items.map(item => `
                    <div class="item-card">
                        <div class="item-header">
                            <div class="item-name">${item.name}</div>
                            <div class="item-price">${formatCurrency(item.total_price, currency)}</div>
                        </div>
                        <div class="item-details">
                            <div><strong>Бренд:</strong> ${item.brand}</div>
                            <div><strong>Размер:</strong> ${item.size || 'не указан'}</div>
                            <div><strong>Артикул:</strong> ${item.nm_id}</div>
                            <div><strong>Цена:</strong> ${formatCurrency(item.price, currency)}</div>
                            ${item.sale > 0 ? `<div><strong>Скидка:</strong> ${item.sale}%</div>` : ''}
                            <div><strong>Статус:</strong> <span class="status-badge ${getStatusClass(item.status)}">${getStatusText(item.status)}</span></div>
                        </div>
                    </div>
                `).join('')}
            </div>
        </div>
    `;
};


// Отображает полную информацию о заказе
export const displayOrder = (order) => {
    const orderDetails = document.getElementById('orderDetails');
    const resultsSection = document.getElementById('resultsSection');

    if (!orderDetails || !resultsSection) {
        console.error('Required DOM elements not found');
        return;
    }

    const orderInfo = document.createElement('div');
    orderInfo.className = 'order-info fade-in';

    const currency = order.payment?.currency || 'USD';

    orderInfo.innerHTML =
        renderOrderBasicInfo(order) +
        renderDeliveryInfo(order.delivery) +
        renderPaymentInfo(order.payment) +
        renderItemsList(order.items, currency);

    orderDetails.innerHTML = '';
    orderDetails.appendChild(orderInfo);

    resultsSection.style.display = 'block';
    resultsSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
};
