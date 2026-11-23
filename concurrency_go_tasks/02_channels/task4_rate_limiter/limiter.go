package limiter

import (
	"time"
)

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
// TODO: инициализировать канал токенов и запуск пополнения
func NewLimiter() *Limiter {
	l := &Limiter{}
	l.tokens = make(chan struct{}, 5)
	l.stop = make(chan struct{})
	for i := 0; i < 5; i++ {
		l.tokens <- struct{}{}
	}
	go func() {
		for {
			select {
			case <-l.stop:
				return
			case <-time.After(200 * time.Millisecond):
				select {
				case <-l.stop:
					return
				default:
				}

				select {
				case l.tokens <- struct{}{}:
				default:
				}
			}
		}
	}()
	return l
}

// Allow возвращает true, если событие разрешено в текущий момент.
// TODO: реализовать получение токена из канала
func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

// Stop останавливает лимитер.
// TODO: остановить пополнение токенов
func (l *Limiter) Stop() {
	select {
	case <-l.stop:
	default:
		close(l.stop)
	}
}
