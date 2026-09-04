package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx
	ch := make(chan int)

	x2 := func() {
		defer close(ch)
		for i := 0; i < len(nums); i++ {
			select {
			case ch <- nums[i] * 2:
			case <-ctx.Done():
				return
			}
		}
	}

	go x2()

	sum := 0
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				// Канал закрыт — первая стадия закончила работу.
				return sum, nil
			}

			sum += value

		case <-ctx.Done():
			return sum, ctx.Err()
		}
	}
}
