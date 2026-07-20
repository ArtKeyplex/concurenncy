package pipeline

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	// TODO: реализовать конвейер обработки чисел

	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	resultCh := make(chan int)

	go func() {
		for num := range ch1 {
			ch2 <- num * num
		}
		close(ch2)
	}()

	go func() {
		for num := range ch2 {
			ch3 <- num * 2
		}
		close(ch3)
	}()

	go func() {
		var result int

		for num := range ch3 {
			result += num
		}

		resultCh <- result
	}()

	for _, num := range nums {
		ch1 <- num
	}
	close(ch1)

	return <-resultCh
}
