package pool

import (
	"sync"
	"time"
)

func RunPool(jobs []int, workers int) int {
	jobsCh := make(chan int, len(jobs))
	resultsCh := make(chan int, len(jobs))

	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				time.Sleep(10 * time.Millisecond)
				resultsCh <- job 
			}
		}()
	}

	for _, job := range jobs {
		jobsCh <- job
	}
	close(jobsCh)

	wg.Wait()
	close(resultsCh)

	sum := 0
	for res := range resultsCh {
		sum += res
	}

	return sum
}