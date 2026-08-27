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
	var one sync.Once

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return

			case <-ticker.C:
				f()
			}
		}
	}()

	return func() {
		one.Do(func() {
			close(stopCh)
		})
	}

}
