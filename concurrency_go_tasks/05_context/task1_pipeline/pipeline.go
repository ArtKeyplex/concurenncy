package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	doubledChan := make(chan int)

	go func() {
		defer close(doubledChan)

		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case doubledChan <- n * 2:
			}
		}
	}()

	sum := 0

	for {
		select {
		case <-ctx.Done():
			return sum, ctx.Err()
		case n, ok := <-doubledChan:
			if !ok {
				return sum, nil
			}
			sum += n
		}
	}
	// TODO: реализовать конвейер с остановкой по ctx
}
