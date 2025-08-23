export const API_BASE_URL = 'http://localhost:8081/api/v1';

export const searchOrder = async (orderUID) => {
    const response = await fetch(`${API_BASE_URL}/orders/${encodeURIComponent(orderUID)}`);
    return await response.json();
};

export const getCacheStats = async () => {
    const response = await fetch(`${API_BASE_URL}/cache/stats`);
    return await response.json();
};

export const getAllOrders = async (limit = 10) => {
    const response = await fetch(`${API_BASE_URL}/orders?limit=${limit}`);
    return await response.json();
};

export const healthCheck = async () => {
    const response = await fetch(`${API_BASE_URL}/health`);
    return await response.json();
};

export const createRandomOrder = async () => {
    const response = await fetch(`${API_BASE_URL}/orders/random`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
    });

    if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return await response.json();
};

export const getRandomOrderExample = () => {
    const exampleOrders = [
        'b563feb7b2b84b6test',
        'sample_order_001',
        'demo_order_123',
        'test_order_456'
    ];

    return exampleOrders[Math.floor(Math.random() * exampleOrders.length)];
};
