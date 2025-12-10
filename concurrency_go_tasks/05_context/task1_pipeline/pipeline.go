package pipelinectx

import "context"

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.
func Run(ctx context.Context, nums []int) (int, error) {
	// TODO: реализовать конвейер с остановкой по ctx
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	two, err := double(ctx, nums)
	if err != nil {
		return 0, err
	}
	plus, err := sum(ctx, two)
	if err != nil {
		return 0, err
	}
	return plus, ctx.Err()
}

func double(ctx context.Context, nums []int) ([]int, error) {
	for i := range nums {
		nums[i] *= 2
	}
	return nums, ctx.Err()
}

func sum(ctx context.Context, nums []int) (int, error) {
	var total int
	for i := range nums {
		total += nums[i]
	}
	return total, ctx.Err()

}
