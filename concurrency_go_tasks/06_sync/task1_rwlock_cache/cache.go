package cache

import "sync"

// Cache представляет потокобезопасный кэш.
type Cache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// New создаёт новый кэш.
func New() *Cache {
	// TODO: инициализировать структуру кэша
	return &Cache{mu: sync.RWMutex{}, data: make(map[string]interface{})}
}

// Set сохраняет значение по ключу.
func (c *Cache) Set(key string, value interface{}) {
	// TODO: реализовать запись с использованием RWMutex
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

// Get возвращает значение по ключу и признак его наличия.
func (c *Cache) Get(key string) (interface{}, bool) {
	// TODO: реализовать чтение с использованием RWMutex
	c.mu.RLock()
	res, ok := c.data[key]
	c.mu.RUnlock()
	return res, ok
}
