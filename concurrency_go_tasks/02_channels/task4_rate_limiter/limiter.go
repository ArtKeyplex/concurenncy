package limiter

import (
	"sync"
	"time"
)

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
	end    chan struct{}
	wg     sync.WaitGroup
	once   sync.Once
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
func NewLimiter() *Limiter {
	// TODO: инициализировать канал токенов и запуск пополнения
	ch := &Limiter{tokens: make(chan struct{}, 5), end: make(chan struct{}), once: sync.Once{}}
	for i := 1; i < 6; i++ {
		ch.tokens <- struct{}{}
	}
	ch.wg.Add(1)
	go func() {
		defer ch.wg.Done()
		tick := time.NewTicker(200 * time.Millisecond)
		for {
			select {
			case <-tick.C:
				select {
				case ch.tokens <- struct{}{}:
				default:
				}
			case <-ch.end:
				return
			}
		}
	}()
	return ch
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
	l.once.Do(func() {
		close(l.end)
	})
	l.wg.Wait()
}
