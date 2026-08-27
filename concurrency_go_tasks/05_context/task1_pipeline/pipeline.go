package pipelinectx

import (
	"context"
)

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {

	generator := func(ctx context.Context, nums []int) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)

			for _, n := range nums {
				select {
				case <-ctx.Done():
					return

				case ch <- n:
				}
			}

		}()
		return ch

	}

	double := func(ctx context.Context, in <-chan int) <-chan int {
		ch := make(chan int)

		go func() {
			defer close(ch)
			for n := range in {
				select {
				case <-ctx.Done():
					return
				case ch <- n * 2:
				}
			}
		}()
		return ch
	}

	inCh := generator(ctx, nums)
	doubleCh := double(ctx, inCh)

	sum := 0
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case v, ok := <-doubleCh:
			if !ok {
				return sum, nil
			}
			sum += v
		}
	}
}
