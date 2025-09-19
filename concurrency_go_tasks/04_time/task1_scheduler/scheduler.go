package scheduler

import "time"

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	// TODO: периодический вызов функции с возможностью остановки

	ticker := time.NewTicker(d)

	stopChan := make(chan struct{})

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f()
			case <-stopChan:
				return
			}
		}
	}()
	stop = func() {
		close(stopChan)
	}

	return stop
}
