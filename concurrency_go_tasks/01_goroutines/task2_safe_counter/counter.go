// Counter хранит целое значение и мьютекс для безопасного доступа.
// Inc увеличивает счётчик на 1 с защитой от гонок.
// Value возвращает текущее значение счётчика безопасно для гонок.

// TODO: реализовать инкремент с использованием мьютекса
// TODO: вернуть значение с учётом блокировки

package counter

import (
	"sync"
)

type Counter struct {
	mu sync.Mutex
	v  int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.v++
	c.mu.Unlock()
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}
