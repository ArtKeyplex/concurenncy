package limiter

import (
	"time"
)

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
func NewLimiter() *Limiter {
	lim := &Limiter{
		tokens: make(chan struct{}, 5),
	}

	for i := 0; i < 5; i++ {
		lim.tokens <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			select {
			case lim.tokens <- struct{}{}:
			default:
			}
		}
	}()

	return lim
}

// Allow возвращает true, если событие разрешено в текущий момент.
func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}
