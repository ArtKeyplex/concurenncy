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

	wg := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	jobsCh := make(chan int, len(jobs))
	results := make(chan int)

	for _, j := range jobs {
		jobsCh <- j
	}
	close(jobsCh)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				results <- job
			}
		}()
	}

	var sum int
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for v := range results {
			sum += v
		}
	}()

	wg.Wait()
	close(results)

	wg2.Wait()
	return sum
}
