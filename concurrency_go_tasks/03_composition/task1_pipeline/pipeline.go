package pipeline

import "sync"

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	squareCh := make(chan int)
	multiplierCh := make(chan int)
	sumCh := make(chan int)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer close(squareCh)
		defer wg.Done()

		for _, num := range nums {
			squareCh <- num * num
		}
	}()

	wg.Add(1)
	go func() {
		defer close(multiplierCh)
		defer wg.Done()

		for num := range squareCh {
			multiplierCh <- 2 * num
		}
	}()

	wg.Add(1)
	go func() {
		defer close(sumCh)
		defer wg.Done()

		sum := 0

		for num := range multiplierCh {
			sum += num
		}
		sumCh <- sum
	}()

	go func() {
		wg.Wait()
	}()

	return <-sumCh
}
