package scheduler

import (
	"time"
)

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	stopChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				f()
			}
		}
	}()

	return func() {
		select {
		case <-stopChan:
		default:
			close(stopChan)
		}
	}
}
