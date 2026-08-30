package pool

import (
	"sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {

	if workers <= 0 {
		workers = 1
	}

	jobsCh := make(chan int)
	resultCh := make(chan int)

	go func() {
		for _, job := range jobs {
			jobsCh <- job
		}
		close(jobsCh)
	}()

	var wg sync.WaitGroup
	wg.Add(workers)

	worker := func(jobs <-chan int, results chan<- int) {
		defer wg.Done()

		for job := range jobs {
			results <- job
		}
	}

	for i := 0; i < workers; i++ {
		go worker(jobsCh, resultCh)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	sum := 0
	for i := range resultCh {
		sum += i
	}
	return sum
}
