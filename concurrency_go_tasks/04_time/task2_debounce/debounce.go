package debounce

import "time"

// Debounce принимает значения и отдаёт только последнее после паузы d.
func Debounce(d time.Duration, in <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		var lastValue int
		var timer *time.Timer

		for value := range in {
			lastValue = value

			if timer != nil {
				timer.Stop()
			}

			timer = time.AfterFunc(d, func() {
				ch <- lastValue
			})
		}

		if timer != nil {
			timer.Stop()
			ch <- lastValue
		}
	}()

	return ch
}
