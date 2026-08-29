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
func NewLimiter() *Limiter {
	// TODO: инициализировать канал токенов и запуск пополнения
	ch := make(chan struct{}, 5)
	stop := make(chan struct{})
	ticker := time.NewTicker(200 * time.Millisecond)

	for range 5 {
		ch <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case ch <- struct{}{}:
				default:
					continue
				}
			case <-stop:
				return
			}
		}
	}()

	return &Limiter{tokens: ch, stop: stop}
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
	close(l.stop)
}
