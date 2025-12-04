package pool

import (
	"sync"
	_ "sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	if workers <= 0 {
		workers = 1
	}

	jobsChan := make(chan int, len(jobs))
	resultChan := make(chan int, len(jobs))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobsChan, resultChan, &wg)
	}

	for _, job := range jobs {
		jobsChan <- job
	}
	close(jobsChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	total := 0
	for result := range resultChan {
		total += result
	}

	return total
}

func worker(jobs <-chan int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		//time.Sleep(time.Second * 1) // имитация работы
		result <- job
	}
}
