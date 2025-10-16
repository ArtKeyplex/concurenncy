package pipelinectx

import (
	"context"
)

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	doubledCh := make(chan int)

	go func() {
		defer close(doubledCh)

		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case doubledCh <- num * 2:
			}
		}
	}()

	sum := 0
	for num := range doubledCh {
		select {
		case <-ctx.Done():
			return sum, ctx.Err()
		default:
			sum += num
		}
	}

	if err := ctx.Err(); err != nil {
		return sum, err
	}

	return sum, nil
}
