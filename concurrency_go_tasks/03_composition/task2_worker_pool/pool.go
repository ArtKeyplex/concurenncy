package pool

import (
	"sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	if len(jobs) == 0 {
		return 0
	}

	if workers <= 0 {
		workers = 1
	}

	chJobs := make(chan int, len(jobs))
	chResult := make(chan int, len(jobs))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range chJobs {
				chResult <- n
			}
		}()
	}

	for _, v := range jobs {
		chJobs <- v
	}
	close(chJobs)

	go func() {
		wg.Wait()
		close(chResult)
	}()

	total_sum := 0
	for sum := range chResult {
		total_sum += sum
	}

	return total_sum
}
