package pool

import (
	"sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	jobCh := make(chan int)
	resultCh := make(chan int)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobCh {
				resultCh <- job
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, job := range jobs {
			jobCh <- job
		}
		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	sum := 0
	for result := range resultCh {
		sum += result
	}

	return sum
}
