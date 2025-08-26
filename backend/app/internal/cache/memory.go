package cache

import (
	"container/list"
	"order-service/internal/models"
	"sync"

	"github.com/sirupsen/logrus"
)

// MemoryCache представляет LRU кеш в памяти для заказов
type MemoryCache struct {
	mu       sync.RWMutex
	capacity int
	cache    map[string]*list.Element
	lru      *list.List
	logger   *logrus.Logger
}

type cacheItem struct {
	key   string
	order *models.OrderFull
}

type OrderCache interface {
	Get(orderUID string) (*models.OrderFull, bool)
	Set(orderUID string, order *models.OrderFull)
	GetStats() CacheStats
	Clear()
	LoadFromDB(orders []models.OrderFull)
}

type CacheStats struct {
	Size     int `json:"size"`
	Capacity int `json:"capacity"`
	Hits     int `json:"hits"`
	Misses   int `json:"misses"`
}

// NewMemoryCache создает новый кеш в памяти с заданной емкостью
func NewMemoryCache(capacity int, logger *logrus.Logger) *MemoryCache {
	return &MemoryCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		lru:      list.New(),
		logger:   logger,
	}
}

// Get получает заказ из кеша
func (c *MemoryCache) Get(orderUID string) (*models.OrderFull, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, exists := c.cache[orderUID]; exists {
		// Перемещаем элемент в начало списка (recently used)
		c.lru.MoveToFront(elem)
		item := elem.Value.(*cacheItem)
		c.logger.WithField("order_uid", orderUID).Debug("Cache hit")
		return item.order, true
	}

	c.logger.WithField("order_uid", orderUID).Debug("Cache miss")
	return nil, false
}

// Set добавляет заказ в кеш
func (c *MemoryCache) Set(orderUID string, order *models.OrderFull) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если элемент уже существует, обновляем его
	if elem, exists := c.cache[orderUID]; exists {
		c.lru.MoveToFront(elem)
		item := elem.Value.(*cacheItem)
		item.order = order
		c.logger.WithField("order_uid", orderUID).Debug("Cache updated")
		return
	}

	// Добавляем новый элемент
	item := &cacheItem{
		key:   orderUID,
		order: order,
	}
	elem := c.lru.PushFront(item)
	c.cache[orderUID] = elem

	// Если превышена емкость, удаляем самый старый элемент
	if c.lru.Len() > c.capacity {
		c.evictOldest()
	}

	c.logger.WithField("order_uid", orderUID).Debug("Cache set")
}

// evictOldest удаляет самый старый элемент из кеша
func (c *MemoryCache) evictOldest() {
	elem := c.lru.Back()
	if elem != nil {
		c.lru.Remove(elem)
		item := elem.Value.(*cacheItem)
		delete(c.cache, item.key)
		c.logger.WithField("evicted_order_uid", item.key).Debug("Cache evicted oldest")
	}
}

// GetStats возвращает статистику кеша
func (c *MemoryCache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return CacheStats{
		Size:     len(c.cache),
		Capacity: c.capacity,
	}
}

// Clear очищает весь кеш
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*list.Element)
	c.lru = list.New()
	c.logger.Info("Cache cleared")
}

// LoadFromDB загружает заказы в кеш из среза
func (c *MemoryCache) LoadFromDB(orders []models.OrderFull) {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for _, order := range orders {
		if count >= c.capacity {
			break
		}

		item := &cacheItem{
			key:   order.OrderUID,
			order: &order,
		}
		elem := c.lru.PushFront(item)
		c.cache[order.OrderUID] = elem
		count++
	}

	c.logger.WithField("loaded_count", count).Info("Cache loaded from database")
}

// GetAllCachedOrders возвращает все заказы из кеша (для отладки)
func (c *MemoryCache) GetAllCachedOrders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]string, 0, len(c.cache))
	for elem := c.lru.Front(); elem != nil; elem = elem.Next() {
		item := elem.Value.(*cacheItem)
		orders = append(orders, item.key)
	}

	return orders
}
