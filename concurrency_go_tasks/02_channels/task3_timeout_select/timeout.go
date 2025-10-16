package timeout

import (
	"context"
	"errors"
	"time"
)

// Work выполняет длительную задачу и возвращает ошибку,
// если она заняла больше 100 мс или контекст был отменён.
func Work(ctx context.Context) error {
	ch := make(chan struct{})

	go func() {
		time.Sleep(200 * time.Millisecond)
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return errors.New("контекст отменён")
	case <-time.After(100 * time.Millisecond):
		return errors.New("таймаут 100 мс")
	case <-ch:
		return nil
	}
}
