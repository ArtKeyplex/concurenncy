package counter

import (
	"sync"
)

// Counter хранит целое значение и мьютекс для безопасного доступа.
type Counter struct {
	mu sync.Mutex
	v  int
}

// Inc увеличивает счётчик на 1 с защитой от гонок.
func (c *Counter) Inc() {
	// TODO: реализовать инкремент с использованием мьютекса
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v++
}

// Value возвращает текущее значение счётчика безопасно для гонок.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: вернуть значение с учётом блокировки
	return c.v
}

// Метод для безопасной установки значения счетчика
func (c *Counter) SetValue(v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v = v
}

// Метод для безопасного уменьшения значения счетчика на 1.
func (c *Counter) Decrement() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v--
}
