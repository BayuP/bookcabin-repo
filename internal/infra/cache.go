package infra

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     interface{}
	expiredAt time.Time
}

type Cache struct {
	mu sync.RWMutex
	m  map[string]cacheItem
}

func NewCache() *Cache {
	return &Cache{m: make(map[string]cacheItem)}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.m[k]
	c.mu.RUnlock()

	if !ok || time.Now().After(item.expiredAt) {
		return nil, false
	}
	return item.value, true
}

func (c *Cache) Set(k string, v interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.m[k] = cacheItem{value: v, expiredAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}
