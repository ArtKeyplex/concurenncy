package debounce

import "time"

func Debounce(d time.Duration, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		var (
			lastVal int
			timer   *time.Timer
		)

		for val := range in {
			lastVal = val
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(d, func() {
				out <- lastVal
			})
		}

		if timer != nil {
			timer.Stop()
			out <- lastVal
		}
		close(out)
	}()

	return out
}