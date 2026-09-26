package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx

	ch := make(chan int)
	go func() {
		for n := range nums {
			n = n * 2
			select {
			case ch <- n:
			case <-ctx.Done():
				return
			}
		}

		close(ch)

	}()

	var n int
	var ok bool

	for {
		select {
		case n, ok = <-ch:

		}
	}

	return 0, nil
}
