package scheduler

import (
	"sync"
	"time"
)

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	// TODO: периодический вызов функции с возможностью остановки
	stopCh := make(chan struct{})
	ticker := time.NewTicker(d)
	var o sync.Once

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f()
			case <-stopCh:
				return
			}
		}
	}()

	return func() {
		o.Do(func() {
			close(stopCh)
		})
	}
}
