package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx

	ch1 := make(chan int)
	ch2 := make(chan int)
	resultCh := make(chan int)

	go func() {
		defer close(ch2)

		for {
			select {
			case num, ok := <-ch1:
				if !ok {
					return
				}
				ch2 <- num * 2

			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(resultCh)

		var result int
		for {
			select {
			case num, ok := <-ch2:

				if !ok {
					resultCh <- result
					return
				}

				result += num

			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(ch1)

		for _, num := range nums {
			select {
			case ch1 <- num:
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
