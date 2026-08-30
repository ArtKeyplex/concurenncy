package debounce

import "time"

// Debounce принимает значения и отдаёт только последнее после паузы d.
func Debounce(d time.Duration, in <-chan int) <-chan int {
	// TODO: реализовать дебаунс значений из канала
	out := make(chan int)

	go func() {
		defer close(out)

		var last int
		var timerCh <-chan time.Time

		timer := time.NewTimer(d)
		if !timer.Stop() {
			<-timer.C
		}
		defer timer.Stop()

		for {
			select {
			case value, ok := <-in:

				if !ok {
					if timerCh != nil {
						<-timer.C
						out <- last
					}
					return
				}

				last = value
				timer.Reset(d)
				timerCh = timer.C

			case <-timerCh:
				out <- last
				timerCh = nil
			}
		}
	}()

	return out
}
