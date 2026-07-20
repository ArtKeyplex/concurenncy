package scheduler

import (
	"sync"
	"time"
)

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	// TODO: периодический вызов функции с возможностью остановки

	ticker := time.NewTicker(d)
	stopCh := make(chan struct{})
	once := sync.Once{}

	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			case <-stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(stopCh)
		})
	}

}
