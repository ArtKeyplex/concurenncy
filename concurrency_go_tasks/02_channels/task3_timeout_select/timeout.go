package timeout

// Work выполняет длительную задачу и возвращает ошибку,
// если она заняла больше 100 мс или контекст был отменён.

// TODO: реализовать через select и time.After

import (
	"context"
	"errors"
	"time"
)

var ErrTimeout = errors.New("work timeout")   // возвращается если работа заняла больше 100 мс
var ErrCanceled = errors.New("work canceled") // возвращается при отмене контекста

func Work(ctx context.Context) error {
	done := make(chan struct{}) // канал done, сигнализирующий, что длительная задача закончилась

	go func() { // запускаем горутину
		time.Sleep(150 * time.Millisecond) // ждем 150 миллисекунд
		close(done)                        // закрываем канал
	}()

	select { // запускаем конструкцию select, ждёт, какое событие произойдёт первым
	case <-done: // ждём сигнал из канала, канал закрылся
		return nil // ошибки нет
	case <-ctx.Done(): // контекст был отменён
		return ErrCanceled // возвращаем "work canceled"
	case <-time.After(100 * time.Millisecond): // прошло 100 миллисекунд
		return ErrTimeout // возвращаем ошибку "work timeout"
	}
}
