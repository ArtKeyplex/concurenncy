package scheduler

import "time"

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	// TODO: периодический вызов функции с возможностью остановки
	st := make(chan int)
	ticker := time.NewTicker(d)
	stop = func() {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()
		st <- 1
		close(st)
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			case <-st:
				return
			}
		}

	}()
	return stop
}
