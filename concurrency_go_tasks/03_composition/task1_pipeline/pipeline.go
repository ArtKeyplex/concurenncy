package pipeline

import "sync"

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {

	stage1 := make(chan int)
	stage2 := make(chan int)
	result := make(chan int)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(stage1)
		for _, num := range nums {
			stage1 <- num * num
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(stage2)
		for num := range stage1 {
			stage2 <- num * 2
		}

	}()
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(result)
		sum := 0
		for num := range stage2 {
			sum += num

		}
		result <- sum
	}()

	go func() {
		wg.Wait()
	}()
	return <-result
}
