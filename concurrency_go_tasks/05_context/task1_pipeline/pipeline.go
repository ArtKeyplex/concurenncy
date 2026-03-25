package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx
	in := make(chan int)
	double := make(chan int)

	go func() {
		defer close(in)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case in <- n:
			}
		}
	}()

	go func() {
		defer close(double)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
					case <-ctx.Done():
						return
					case double <- v * 2:
				}
			}
		}
	}()

	sum := 0

	for {
		select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case v, ok := <-double:
				if !ok {
					return sum, nil
				}
				sum += v
		}
	}

}
