package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx
	if len(nums) == 0 {
		return 0, ctx.Err()
	}

	doubled := make(chan int)

	// Стадия 1: Удвоение
	go func() {
		defer close(doubled)
		for _, v := range nums {
			select {
			case doubled <- v * 2:
			case <-ctx.Done():
				return
			}
		}
	}()
	sum := 0
	for {
		select {
		case v, ok := <-doubled:
			if !ok {
				return sum, nil
			}
			sum += v

		case <-ctx.Done():
			return sum, ctx.Err()
		}
	}
}
