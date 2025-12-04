package pipeline

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	input := make(chan int)

	go func() {
		defer close(input)
		for _, n := range nums {
			input <- n
		}
	}()

	squared := make(chan int)

	go func() {
		defer close(squared)
		for n := range input {
			squared <- n * n
		}
	}()

	multiplied := make(chan int)

	go func() {
		defer close(multiplied)
		for n := range squared {
			multiplied <- n * 2
		}
	}()

	sum := 0
	for n := range multiplied {
		sum += n
	}

	return sum
}
