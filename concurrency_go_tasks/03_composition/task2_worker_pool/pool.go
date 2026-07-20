package pool

import (
	"sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	// TODO: реализовать пул воркеров и сбор результатов
	if workers <= 0 {
		workers = 1
	}

	jobsCh := make(chan int)
	resultCh := make(chan int)
	var wg sync.WaitGroup


	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sum := 0
			for job := range jobsCh {
				sum += job
			}
			resultCh <- sum
		}()
	}

	go func() {
		defer close(jobsCh)

		for _, task := range jobs {
			jobsCh <- task
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result := 0

	for num := range resultCh {
		result += num
	}

	return result
}
