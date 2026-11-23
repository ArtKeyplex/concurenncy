package cache

import "sync"

// Cache представляет потокобезопасный кэш.
type Cache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// New создаёт новый кэш.
// TODO: инициализировать структуру кэша
func New() *Cache { //new создаёт новый кэш
	return &Cache{
		data: make(map[string]interface{}), // cоздаём кэш и сразу создаём пустую карту внутри него
	}
}

// Set сохраняет значение по ключу.
// TODO: реализовать запись с использованием RWMutex
func (c *Cache) Set(key string, value interface{}) { // set сохраняет значение по ключу
	c.mu.Lock()         // закрываем замок для записи, чтобы другие горутины не могли одновременно менять карту
	c.data[key] = value // записываем значение в карту по ключу
	c.mu.Unlock()       // открываем замок после записи
}

// Get возвращает значение по ключу и признак его наличия.
// TODO: реализовать чтение с использованием RWMutex
func (c *Cache) Get(key string) (interface{}, bool) { // get возвращает значение по ключу и признак его наличия
	c.mu.RLock()             // закрываем замок для чтения, чтобы безопасно прочитать данные из карты
	value, ok := c.data[key] // читаем значение из карты, ok будет true, если ключ есть, и false, если ключа нет
	c.mu.RUnlock()           // открываем замок после чтения
	return value, ok         // возвращаем значение и признак наличия ключа
}
