package scheduler

import (
	"time"
)

func Every(d time.Duration, f func()) (stop func()) {
	ticker := time.NewTicker(d)
	done := make(chan struct{})

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
		close(done)
	}
}
