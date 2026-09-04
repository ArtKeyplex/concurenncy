package pool

import (
	"sync"
	_ "sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	// TODO: реализовать пул воркеров и сбор результатов
	if workers <= 0 {
		workers = 1
	}
	ch := make(chan int)

	sum := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	worker := func() {
		defer wg.Done()
		for job := range ch {
			mu.Lock()
			sum += job
			mu.Unlock()
		}
	}

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker()
	}

	for j := 0; j < len(jobs); j++ {
		ch <- jobs[j]
	}

	close(ch)

	wg.Wait()

	return sum
}
