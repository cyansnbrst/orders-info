package cache

import (
	"log"
	"sync"
	"time"

	"cyansnbrst.com/order-info/internal/orders"
)

// Orders cache item
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// Orders cache struct
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// Orders cache constructor
func NewInMemoryCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set new element to cache
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	expiration := time.Now().Add(duration).Unix()
	c.mu.Lock()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
	c.mu.Unlock()
}

// Get an element from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	if time.Now().Unix() > item.Expiration {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	c.mu.RUnlock()
	return item.Value, true
}

// Delete an element from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Clear orders cache
func (c *Cache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]CacheItem)
	c.mu.Unlock()
}

// Recover orders data from cache
func (c *Cache) Recover(ordersRepo orders.Repository) error {
	orders, err := ordersRepo.GetAll()
	if err != nil {
		return err
	}

	for _, order := range orders {
		c.Set(order.OrderUID, order, 5*time.Minute)
	}

	return nil
}

// Orders cache cleaner
func (c *Cache) StartCleaner(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now().Unix()
			c.mu.Lock()
			for key, item := range c.items {
				if now > item.Expiration {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

// Print orders cache
func (c *Cache) PrintCache() {
	log.Printf("current cache state: %v", c.items)
}
