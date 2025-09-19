package limiter

import "time"

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
func NewLimiter() *Limiter {

	lim := &Limiter{
		tokens: make(chan struct{}, 5),
	}
	go lim.refill()
	for i := 0; i < 5; i++ {
		lim.tokens <- struct{}{}
	}
	return lim
}

// TODO: инициализировать канал токенов и запуск пополнения

// Allow возвращает true, если событие разрешено в текущий момент.
func (l *Limiter) refill() {
	// TODO: реализовать получение токена из канала

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C

		for i := 0; i < 5; i++ {
			select {
			case l.tokens <- struct{}{}:
			default:
				return
			}
		}
	}
}

func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	case <-time.After(0):
	}
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}

}
