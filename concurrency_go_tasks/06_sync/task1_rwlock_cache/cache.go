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
	return &Cache{
		sync.RWMutex{},
		make(map[string]interface{}),
	}
}

// Set сохраняет значение по ключу.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
	// TODO: реализовать запись с использованием RWMutex
}

// Get возвращает значение по ключу и признак его наличия.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	val, ok := c.data[key]
	if !ok {
		return nil, ok
	}
	c.mu.RUnlock()
	// TODO: реализовать чтение с использованием RWMutex
	return val, ok
}
