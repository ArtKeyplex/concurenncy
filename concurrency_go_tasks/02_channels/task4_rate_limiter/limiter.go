//Не вывез эту задачу...

package limiter

import _ "time"

// Limiter ограничивает количество событий до 5 в секунду.
type Limiter struct {
	tokens chan struct{}
}

// NewLimiter создаёт новый лимитер с ёмкостью 5 токенов.
func NewLimiter() *Limiter {
	// TODO: инициализировать канал токенов и запуск пополнения
	return &Limiter{}
}

// Allow возвращает true, если событие разрешено в текущий момент.
func (l *Limiter) Allow() bool {

	// TODO: реализовать получение токена из канала
	return false
}

// Stop останавливает лимитер.
func (l *Limiter) Stop() {
	// TODO: остановить пополнение токенов
}
