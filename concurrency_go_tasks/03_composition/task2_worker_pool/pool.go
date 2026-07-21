package pool

import (
	"sync"
	_ "sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	// TODO: реализовать пул воркеров и сбор результатов
	if len(jobs) == 0 {
		return 0
	}
	if workers <= 0 {
		workers = 1
	}
	ch := make(chan int, len(jobs))
	resultCh := make(chan int, len(jobs))
	for _, v := range jobs {
		ch <- v
	}
	close(ch)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for rec := range ch {
				resultCh <- rec
			}
		}()
	}
	wg.Wait()
	close(resultCh)

	var ans int
	for num := range resultCh {
		ans += num
	}

	return ans
}
