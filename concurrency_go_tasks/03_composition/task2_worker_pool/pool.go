package pool

import (
	"sync"
	_ "sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.
func RunPool(jobs []int, workers int) int {
	// TODO: реализовать пул воркеров и сбор результатов
	var mu sync.Mutex
	var wg sync.WaitGroup
	res := 0
	tasks := make(chan int)
	if workers <= 0 {
		workers = 1
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for job := range tasks {
				mu.Lock()
				res += job
				mu.Unlock()
			}
		}()
	}
	go func() {
		for _, job := range jobs {
			tasks <- job
		}
		close(tasks)
	}()
	wg.Wait()
	return res
}
