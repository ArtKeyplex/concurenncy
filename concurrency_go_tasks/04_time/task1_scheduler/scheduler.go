package scheduler

import (
	"sync/atomic"
	"time"
)

// Every запускает f каждые d и возвращает функцию для остановки.
func Every(d time.Duration, f func()) (stop func()) {
	// TODO: периодический вызов функции с возможностью остановки
	ticker := time.NewTicker(d)
	done := make(chan struct{})
	atom := &atomic.Bool{}
	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			case <-done:
				ticker.Stop()
				return
			}

		}
	}()
	return func() {
		if atom.CompareAndSwap(false, true) {
			close(done) // 100 раз закрывать канал)))
		}

	}
}
