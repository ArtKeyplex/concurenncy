package timeout

import (
	"context"
	"fmt"
	"time"
)

// Work выполняет длительную задачу и возвращает ошибку,
// если она заняла больше 100 мс или контекст был отменён.
func Work(ctx context.Context) error {
	ch := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
	}()

	select {
	case result := <-ch:
		fmt.Printf("Work успешно завершен с результатом %d\n", result)
		return nil

	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("операция превысила таймаут 100ms")

	case <-ctx.Done():
		fmt.Printf("Work: операция отменена. Причина: %v\n", ctx.Err())
		return fmt.Errorf("операция отменена: %w", ctx.Err())

	}

}

// TODO: реализовать через select и time.After
