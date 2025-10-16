package scheduler

import "time"

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(d)
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
		close(stopCh)
	}
}
