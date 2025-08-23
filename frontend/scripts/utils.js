export const formatDate = (dateString) => {
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

export const formatTimestamp = (timestamp) => {
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

export const formatCurrency = (amount, currency = 'USD') => {
    if (typeof amount !== 'number') return amount;

    const value = amount / 100;

    const formatter = new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 2
    });

    return formatter.format(value);
}

export const getStatusClass = (status) => {
    if (status >= 200 && status < 300) return 'status-success';
    if (status >= 100 && status < 200) return 'status-warning';
    return 'status-error';
}

export const getStatusText = (status) => {
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