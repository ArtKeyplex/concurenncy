package scheduler

import (
	"sync"
	"time"
)

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	if d <= 0 {
		return func() {}
	}

	ticker := time.NewTicker(d)
	stopChan := make(chan struct{})
	var once sync.Once

	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(stopChan)
		})
	}
}
