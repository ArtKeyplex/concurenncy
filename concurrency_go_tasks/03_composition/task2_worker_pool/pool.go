package pool

import (
	"sync"
	_ "sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	ch1 := make(chan int)
	ch2 := make(chan int)
	res := 0
	var wg sync.WaitGroup
	if workers <= 0 {
		workers = 1
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for num := range ch1 {
				ch2 <- num
			}
		}()
	}
	go func() {
		wg.Wait()
		close(ch2)
	}()

	go func() {
		for _, v := range jobs {
			ch1 <- v
		}
		close(ch1)
	}()

	for v := range ch2 {
		res += v
	}
	return res
}
