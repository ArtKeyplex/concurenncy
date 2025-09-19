package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx

	doubled := make(chan int, len(nums))

	go func() {

		defer close(doubled)
		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case doubled <- num * 2:
			}
		}
	}()

	sum := 0
	for num := range doubled {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			sum += num
		}
	}

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return sum, nil
	}
}
