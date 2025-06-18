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
	return &Cache{}
}

// Set сохраняет значение по ключу.
func (c *Cache) Set(key string, value interface{}) {
	// TODO: реализовать запись с использованием RWMutex
}

// Get возвращает значение по ключу и признак его наличия.
func (c *Cache) Get(key string) (interface{}, bool) {
	// TODO: реализовать чтение с использованием RWMutex
	return nil, false
}
