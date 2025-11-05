package pipeline

import "sync"

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	// TODO: реализовать конвейер обработки чисел
	res := 0
	ch1 := make(chan int)
	ch2 := make(chan int)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for _, v := range nums {
			ch1 <- v * v
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < len(nums); i++ {
			num := <-ch1
			ch2 <- num * 2
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < len(nums); i++ {
			res += <-ch2
		}
	}()
	wg.Wait()
	return res

}
