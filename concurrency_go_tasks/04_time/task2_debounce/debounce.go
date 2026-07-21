package debounce

import (
	"time"
)

// Debounce принимает значения и отдаёт только последнее после паузы d.
func Debounce(d time.Duration, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		var (
			timer    *time.Timer
			timerCh  <-chan time.Time
			lastVal  int
			hasVal   bool
			inClosed bool
		)

		for {
			select {
			case val, ok := <-in:
				if !ok {
					inClosed = true
					in = nil

					if !hasVal {
						return
					}
					continue
				}

				lastVal = val
				hasVal = true

				if timer != nil {
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
				}

				timer = time.NewTimer(d)
				timerCh = timer.C

			case <-timerCh:
				if hasVal {
					out <- lastVal
					hasVal = false
				}
				timerCh = nil

				if inClosed {
					return
				}
			}
		}
	}()

	return out
}
