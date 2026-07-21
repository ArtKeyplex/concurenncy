package limiter

import (
	"sync"
	"time"
	_ "time"
)

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
	stop   chan struct{}
	once   sync.Once
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
func NewLimiter() *Limiter {
	l := &Limiter{
		tokens: make(chan struct{}, 5),
		stop:   make(chan struct{}),
	}

	for i := 0; i < 5; i++ {
		l.tokens <- struct{}{}
	}

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				select {
				case l.tokens <- struct{}{}:
				default:
				}
			case <-l.stop:
				for len(l.tokens) > 0 {
					<-l.tokens
				}
				return
			}
		}
	}()

	return l
}

// Allow возвращает true, если событие разрешено в текущий момент.
func (l *Limiter) Allow() bool {
	select {
	case _, ok := <-l.tokens:
		return ok
	default:
		return false
	}
}

// Stop останавливает лимитер.
func (l *Limiter) Stop() {
	l.once.Do(func() {
		close(l.stop)
	})
}
