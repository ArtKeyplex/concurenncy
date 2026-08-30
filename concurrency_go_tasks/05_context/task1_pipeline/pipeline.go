package pipelinectx

import (
	"context"
)

type result struct {
	value int
	err   error
}

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx
	ch := make(chan int)

	go func() {
		defer close(ch)

		for _, i := range nums {
			select {
			case ch <- i * 2:
			case <-ctx.Done():
				return
			}
		}
	}()

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return counter, ctx.Err()
		case val, ok := <-ch:
			if !ok {
				return counter, nil
			}
			counter += val
		}
	}
}
