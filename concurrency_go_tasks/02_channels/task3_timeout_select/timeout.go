package timeout

import (
	"context"
	"errors"
)

var ErrTimeout = errors.New("work timeout") // возвращается если работа заняла больше 100 мс
var ErrCanceled = errors.New("work canceled") // возвращается при отмене контекста

// Work выполняет длительную задачу и возвращает ошибку,
// если она заняла больше 100 мс или контекст был отменён.
func Work(ctx context.Context) error {
	// TODO: реализовать через select и time.After
	return nil
}
