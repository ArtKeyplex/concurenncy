package counter

import "sync"

// Counter хранит целое значение и мьютекс для безопасного доступа.
type Counter struct {
	mu sync.Mutex
	v  int
}

// Inc увеличивает счётчик на 1 с защитой от гонок.
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v++
}

// Value возвращает текущее значение счётчика безопасно для гонок.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}
