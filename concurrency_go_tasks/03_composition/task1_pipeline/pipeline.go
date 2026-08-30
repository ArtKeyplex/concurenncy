package pipeline

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	// TODO: реализовать конвейер обработки чисел
	squareCh := make(chan int)
	sumCh := make(chan int)
	result := make(chan int)

	go func() {
		for _, i := range nums {
			squareCh <- i * i
		}
		close(squareCh)
	}()

	go func() {
		for i := range squareCh {
			sumCh <- i * 2
		}
		close(sumCh)
	}()

	go func() {
		counter := 0
		for i := range sumCh {
			counter += i
		}
		result <- counter
	}()
	return <-result
}
