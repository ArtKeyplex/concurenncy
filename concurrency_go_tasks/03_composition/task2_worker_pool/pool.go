package pool

import "sync"

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	// TODO: реализовать пул воркеров и сбор результатов

	jobChan := make(chan int, len(jobs))
	resultChan := make(chan int, len(jobs))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobChan, resultChan, &wg)
	}
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	sum := 0
	for result := range resultChan {
		sum += result
	}
	return sum
}

func worker(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		results <- job

	}
}
