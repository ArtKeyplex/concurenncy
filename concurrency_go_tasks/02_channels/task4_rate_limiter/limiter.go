package limiter

import (
	"sync/atomic"
	"time"
)

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
	done   chan struct{}
	atomic.Bool
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
func NewLimiter() *Limiter {
	// TODO: инициализировать канал токенов и запуск пополнения
	l := &Limiter{
		tokens: make(chan struct{}, 5),
		done:   make(chan struct{}),
	}
	for i := 0; i < 5; i++ {
		l.tokens <- struct{}{}
	}
	go func() {
		ticker := time.NewTicker(time.Millisecond * 200)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				select {
				case l.tokens <- struct{}{}:
				default:
				}
			case <-l.done:
				return

			}
		}
	}()

	return l
}

// Allow возвращает true, если событие разрешено в текущий момент.
func (l *Limiter) Allow() bool {
	// TODO: реализовать получение токена из канала

	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

// Stop останавливает лимитер.
func (l *Limiter) Stop() {
	// TODO: остановить пополнение токенов
	if l.CompareAndSwap(false, true) {
		close(l.done)
	}
}
