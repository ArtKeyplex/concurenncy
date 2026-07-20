package debounce

import (
	"time"
)

// Debounce принимает значения и отдаёт только последнее после паузы d.
func Debounce(d time.Duration, in <-chan int) <-chan int {
	// TODO: реализовать дебаунс значений из канала

	out := make(chan int)

	go func() {
		defer close(out)

		var lastValue int
		hasLast := false
		timer := time.NewTimer(d)

		for {
			select {
			case val, ok := <-in:
				if !ok {
					if hasLast {
						<-timer.C
						out <- lastValue
					}
					return
				}

				lastValue = val
				hasLast = true
				timer.Reset(d)

			case <-timer.C:
				if hasLast {
					out <- lastValue
					hasLast = false
					timer.Stop()
				}
			}
		}
	}()

	return out
}
