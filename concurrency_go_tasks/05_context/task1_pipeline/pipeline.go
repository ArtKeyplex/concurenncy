package pipelinectx

import (
	"context"
	"sync"
)

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx
	res := 0
	ch := make(chan int)
	ch2 := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(ch)
		for _, i := range nums {
			ch <- i * 2
		}
		return

	}()
	go func() {
		defer wg.Done()
		defer close(ch2)
		for num := range ch {
			res += num
		}
		ch2 <- 1

	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-ch2:
		return res, nil
	}

}
